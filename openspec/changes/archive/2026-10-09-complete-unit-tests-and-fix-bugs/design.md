# Design

## Context

A motivação está em proposal.md (seção Why) e o comportamento esperado nas specs `purchase-balance`, `stock-quote`, `reit-quote` e `cli`. Este documento cobre o estado atual que molda a abordagem. Todos os defeitos citados foram reproduzidos durante a exploração, com testes temporários via `go test -overlay` e com o `govulncheck`.

**Saída e encerramento do CLI**
- Os handlers de `cmd/` escrevem diretamente em `os.Stdout`: `common.NewTableWriter` fixa `os.Stdout` e os erros usam `fmt.Printf`.
- Erros de "not found" saem em stdout com código 0.
- `core.App.Run` chama `os.Exit(1)` por conta própria, o que impede testar o código de saída.

**Flags globais**
- `flagNoColor` e `flagCsv` são variáveis de pacote compartilhadas por todos os subcomandos.
- Um `--csv` em uma execução vaza para a próxima execução da mesma árvore.

**`getHtml`**
- Ignora o erro de `http.NewRequest`: um ticker como `%zz` gera uma requisição nula e causa panic.
- Com gzip inválido, chama `Close` em um `*gzip.Reader` nulo e causa panic.
- Não fecha o corpo da resposta quando o status é diferente de 200.
- Recebe o timeout como `time.Duration`, mas o multiplica por `time.Second`. `config` devolve `Duration(60)`, ou seja, 60 ns usados como contagem.

**Scraping de FII**
- Ignora o erro de `getHtml` e a `rs.url` da instância.
- O StatusInvest responde 200 tanto para FII inexistente quanto para `fundos-imobiliarios/` (ticker vazio). Isso foi verificado em 2026-10-09. Para ações inexistentes, responde 404.

**Testes de scraping**
- Os testes externos (`package scraping_test`) só constroem instâncias com a URL fixa do StatusInvest, então acessam a rede.

**Dependências**
- `golang.org/x/net` v0.33.0 e `github.com/antchfx/xpath` v1.3.3 têm 9 vulnerabilidades alcançáveis via `htmlquery.Parse` → `html.Parse`.

**CI**
- O workflow de release fixa Go 1.24 (o `go.mod` declara 1.25.0).
- Não roda testes.
- Usa `actions/upload-release-asset@v1`, que está arquivada e depende de `github.event.release.upload_url`. Esse valor é vazio no `workflow_dispatch`, então o disparo manual falha.

## Goals / Non-Goals

**Goals:**
- Nenhuma entrada do usuário ou resposta do provedor pode causar panic ou travamento.
- Cada cenário das specs tem ao menos um teste unitário determinístico, sem rede.
- As dependências ficam sem vulnerabilidades conhecidas alcançáveis, e o CI impede regressões.
- O repositório fica pronto para o próximo ciclo: CI em PR, release funcional, licença e documentação de desenvolvimento.

**Non-Goals:**
- Redesenhar a heurística de `PurchaseBalance` ("-1 cota" quando a parte é menor que o ativo mais caro). Ela é comportamento do produto, não defeito.
- Consultar tickers em paralelo, suportar outras moedas ou outros provedores. São features do próximo ciclo.
- Trocar cobra, go-pretty ou htmlquery.

## Decisions

### D1. `PurchaseBalance`: filtrar inválidos e redistribuir só entre os comprados
- Logo no início, montar a lista de ativos válidos (não nulos e com preço > 0) e usá-la em todo o cálculo.
- Fazer o loop de redistribuição percorrer `securitiesPurchaseSort` (os comprados) e encerrar quando o saldo for menor que o preço do mais barato entre eles, ou quando não houver comprado nenhum.
- **Por que termina:** enquanto o saldo for maior ou igual ao menor preço comprado, cada rodada compra ao menos esse ativo, reduzindo o saldo em um valor positivo. Os ativos caros não comprados custam mais que o saldo do momento em que foram avaliados, e o saldo só diminui. Por isso, ao fim, o saldo é menor que todos os preços válidos, como exige a spec.
- **Alternativa descartada:** `recover` ou um limite de iterações. Isso esconde o defeito e viola o invariante do saldo.

### D2. `getHtml` robusto e timeout como duração real
- Tratar o erro de `http.NewRequest`.
- Fechar o corpo com `defer` logo após um `client.Do` bem-sucedido, qualquer que seja o status.
- Só registrar o `Close` do leitor gzip depois de verificar o erro de `gzip.NewReader`.
- Manter a estrutura linear com early returns.
- `config.GetScrapingTimeout()` passa a devolver `60 * time.Second`, e `getHtml` usa o valor diretamente. Os testes passam `time.Second` ou menos.
- **Alternativa descartada:** manter a contagem em segundos e só renomear. O tipo `time.Duration` continuaria mentindo sobre a unidade.

### D3. Validação de ticker e consistência entre ação e FII
- Uma função `validateTicker` no pacote `scraping` aceita só `^[A-Za-z0-9]+$` e é chamada no início de `GetStockByTicker` e de `GetReitByTicker`. O ticker inválido gera erro sem requisição, o que cobre o vazio, `%zz`, espaços e `../`.
- `GetReitByTicker` propaga o erro de `getHtml` e monta a URL com `rs.url`.
- Os dois tratam o erro de `htmlquery.Parse`.
- O preço passa por `strings.TrimSpace` antes de `tools.ToFloat`.
- O nome da ação usa `strings.SplitN(title, "-", 2)` e só é preenchido quando há 2 partes.
- **Alternativa descartada:** `url.PathEscape` sem validar. Evitaria o panic, mas ainda geraria requisições inúteis e não resolveria o ticker vazio.

### D4. Testes de scraping com `httptest` e fixtures HTML, dentro do pacote
- Novos arquivos `internal/scraping/*_internal_test.go`, em `package scraping`, constroem `&stockScraping{url: srv.URL, config: fake}` e `&reitScraping{...}` diretamente. A API pública não muda.
- As fixtures ficam em `internal/scraping/testdata/`. São HTML mínimos escritos à mão, reproduzindo só a estrutura que os XPaths atuais alcançam.
- **Alternativa descartada:** um `RoundTripper` falso no `http.DefaultTransport`. Troca um estado global por outro e é menos legível.

### D5. Testes de rede atrás da build tag `integration`
- `internal/scraping/stock_test.go` e `reit_test.go` ganham `//go:build integration`.
- O alvo `make test/integration` executa `go test -tags integration ./internal/scraping/...`.
- **Alternativa descartada:** `testing.Short()`. O padrão seria acessar a rede.

### D6. Saída injetável, erros via `RunE` e código de saída no `main`
- `common.NewTableWriter(noColor bool, out io.Writer)`. Os handlers passam a ser `RunE` e escrevem resultados em `cmd.OutOrStdout()`.
- Para "not found", o handler define `cmd.SilenceUsage = true` e devolve `fmt.Errorf("ticker %q not found!", t)` ou `errors.New("tickers not found!")`. O cobra imprime `Error: <mensagem>` no stderr do comando. Os erros de uso (argumentos, flag obrigatória) continuam mostrando o uso.
- `rootCommand` também escreve a versão em `cmd.OutOrStdout()`.
- `core.App.Run(args []string) int` recebe os argumentos explicitamente (repassados ao cobra com `SetArgs`) e devolve 0 ou 1, sem chamar `os.Exit`. `main.go` passa a ter `run(args []string) int`, e `main()` fica reduzido a `os.Exit(run(os.Args[1:]))`.
- **Testes de `cmd/`:** cada teste monta uma árvore nova com serviços falsos e chama `SetOut`, `SetErr` (buffers separados), `SetArgs` e `Execute()`. Verifica stdout, stderr e o erro devolvido.
- **Testes de `core`:** chamam `Run` com argumentos que não acessam a rede (`[]string{"version"}`, `[]string{}` e `[]string{"stock", "get"}`) e verificam os códigos 0 e 1.
- **`main_test.go`:** chama `run([]string{})`. Precisa ser uma slice vazia, não `nil`, porque com `nil` o cobra usaria os argumentos do binário de teste.
- **Alternativa descartada:** capturar `os.Stdout` com `os.Pipe`. É frágil, não separa stderr e impede testes em paralelo.

### D7. Flags por comando
- `noColor` e `csv` viram campos de `stockCommand`, `reitCommand` e `securityCommand`, ao lado de `flagAmount`. `cmd/common.go` passa a conter só a função compartilhada `render(t table.Writer, csv bool)`.
- Como cada subcomando registra a flag apontando para o campo do próprio struct, o valor não vaza entre comandos.

### D8. Flags de `purchase-balance` e tickers vazios
- `MarkFlagRequired("amount")` nos três `purchase-balance`.
- Em `security purchase-balance`, `tools.SplitList(s, ",")` faz split, `TrimSpace` e descarta vazios.

### D9. Dependências e vulnerabilidades
- `go get golang.org/x/net@latest github.com/antchfx/xpath@latest github.com/antchfx/htmlquery@latest`, a atualização das demais dependências diretas (cobra, go-pretty) para o patch ou minor mais recente, e `go mod tidy`.
- O critério de aceite é o `govulncheck ./...` sem vulnerabilidades alcançáveis.
- Novo alvo `make test/vuln`, que executa `go run golang.org/x/vuln/cmd/govulncheck@latest ./...`.

### D10. CI e release
- **Novo `.github/workflows/ci.yaml`** (push e pull_request):
  - `actions/checkout@v4` e `actions/setup-go@v5` com `go-version-file: go.mod`;
  - `go vet ./...`;
  - `go test -race -coverpkg=./... ./...`;
  - `govulncheck`.
- **`build_and_release.yaml`:**
  - mesmo setup do Go;
  - um job `test` do qual os 3 builds dependem (`needs: test`);
  - upload trocado por `gh release upload "$TAG_NAME" <arquivo> --clobber`, com `GH_TOKEN`, o que funciona tanto em `release: created` quanto em `workflow_dispatch`.
- **Alternativa descartada:** matriz única para os 3 alvos. Seria uma melhoria estética que aumenta o diff sem resolver nenhum problema.

### D11. Limpeza do repositório
- **`.gitignore`:** `trader*` passa a ser `/trader`, `/trader.exe` e `/trader-*`, para os pacotes de release.
- **Código morto:** remover `tools.ToSnakeCase` e seus testes.
- **`LICENSE`:** texto MIT com "Copyright (c) 2025 Douglas Panhota". O ano é o do primeiro commit.
- **README:** seção "Desenvolvimento" (requisitos, `make build/*`, `make test/unit`, `make test/integration`, `make test/vuln`, cobertura); nota sobre erros em stderr com código 1 e `--amount` obrigatório; correção de erros de digitação ("portifolio", "disponivel", "o preços").

### D12. Cobertura entre pacotes
- `make test/unit` executa `go test -coverpkg=./... -coverprofile=coverage.out ./...`.
- Não há meta numérica obrigatória. O critério é cobrir todos os cenários das specs, e cada função nova ou alterada fica com cobertura total.

## Risks / Trade-offs

- **[Fixtures divergem do site real]** → `make test/integration` detecta mudanças de layout. As fixtures documentam o contrato do parser.
- **[Atualizar dependências muda o comportamento de parse ou de renderização]** → Os testes unitários de scraping e de `cmd` (comparando CSV) rodam depois da atualização, e `make test/integration` é executado antes do arquivamento.
- **[Erros em stderr com código 1 quebram scripts que liam stdout]** → Mudança registrada como BREAKING na proposal e documentada no README. O texto da mensagem não muda.
- **[Validação `[A-Za-z0-9]+` rejeita um formato de ticker legítimo]** → Os tickers da B3 (ações, units, FIIs, BDRs, fracionário com sufixo F) são alfanuméricos. Se surgir exceção, ampliar a expressão é uma mudança local.
- **[`go run govulncheck@latest` depende da rede e de uma versão móvel]** → É aceitável para o CI e para um alvo explícito. Não faz parte de `test/unit`.
- **[Mudar a assinatura de `NewTableWriter` e de `App.Run`]** → São APIs internas, e todos os chamadores são atualizados na mesma change.

## Migration Plan

Não há dados para migrar. As mudanças visíveis (`--amount` obrigatório, erros em stderr com código 1) entram na próxima release e ficam documentadas no README. O novo workflow de CI passa a valer no primeiro push. Para reverter, basta o revert do commit, e as actions antigas voltam junto.
