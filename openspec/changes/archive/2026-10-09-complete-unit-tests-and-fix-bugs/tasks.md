# Tasks

## 1. Dependências e segurança

- [x] 1.1 Atualizar `golang.org/x/net`, `github.com/antchfx/xpath` e `github.com/antchfx/htmlquery` para a versão mais recente, junto com cobra e go-pretty, e rodar `go mod tidy` (D9). Verificar que `go build ./...` e os testes atuais offline (`go test . ./internal/tools ./internal/common ./internal/config ./internal/service`) continuam passando
- [x] 1.2 Adicionar o alvo `test/vuln` ao `Makefile` (`go run golang.org/x/vuln/cmd/govulncheck@latest ./...`, D9). Verificar com `make test/vuln`, que deve relatar 0 vulnerabilidades alcançáveis pelo código

## 2. Modelos e utilitários

- [x] 2.1 Criar `internal/resource/currency_test.go` e `purchase_test.go` cobrindo `Currency.String` ("R$ BRL"), `TotalAmount`, `TotalCount`, `AmountSpent` e `RemainingBalance` (incluindo lista vazia). Verificar com `go test -cover ./internal/resource`, que deve mostrar 100%
- [x] 2.2 Adicionar `tools.SplitList(input, sep string) []string` (split + `TrimSpace` + descarte de vazios, D8) com testes de tabela para "", "A", "A,B", " A , ,B " e ",". Verificar com `go test ./internal/tools`
- [x] 2.3 Remover `tools.ToSnakeCase` e seus testes (D11). Verificar com `grep -rn ToSnakeCase .`, que não deve encontrar nada, e com `go build ./... && go test ./internal/tools`

## 3. Balanceamento de compra (spec `purchase-balance`)

- [x] 3.1 Em `PurchaseBalance`, filtrar no início os ativos nulos ou com preço <= 0 e usar só os válidos no cálculo (D1). Verificar com um teste em que "ZERO" a 0 e "A" a 10, com 100, termina, exclui "ZERO" e dá 10 cotas para "A"
- [x] 3.2 Fazer o loop de redistribuição percorrer só os ativos comprados e parar quando o saldo for menor que o menor preço entre eles, ou quando não houver comprado nenhum (D1). Verificar com um teste em que "High" a 900 e "Low" a 100, com 1000, termina sem panic, exclui "High", dá 10 cotas para "Low" e deixa saldo 0
- [x] 3.3 Reativar o teste comentado com as expectativas corretas ("Low"=4, "High"=1, nessa ordem, saldo 0) e adicionar testes para lista vazia (sem ativos, saldo = valor investido) e para "Ticker sem cotação" em `PurchaseBalancesBySecurities`. Verificar com `go test -cover ./internal/service`, com `PurchaseBalance` em 100%
- [x] 3.4 Adicionar um teste de propriedade com casos fixos variados que confirme que o saldo restante fica entre 0 e o preço de cada ativo válido. Verificar com `go test -run Property ./internal/service`

## 4. Scraping (specs `stock-quote` e `reit-quote`)

- [x] 4.1 Marcar `internal/scraping/stock_test.go` e `reit_test.go` com `//go:build integration` e adicionar o alvo `test/integration` ao `Makefile` (D5). Verificar que `go test ./internal/scraping` não acessa a rede e que `go vet -tags integration ./internal/scraping` compila os testes de rede
- [x] 4.2 Tornar `getHtml` robusto: tratar o erro de `http.NewRequest`, fechar o corpo qualquer que seja o status e registrar o `Close` do gzip só após um `gzip.NewReader` sem erro (D2). Ampliar `common_test.go` com os casos de URL inválida (`%zz`, devolve erro sem panic), gzip inválido (devolve erro sem panic) e corpo fechado em status 500 (servidor que registra o fechamento da conexão). Verificar com `go test -race ./internal/scraping`
- [x] 4.3 Fazer `config.GetScrapingTimeout()` devolver `60 * time.Second` e `getHtml` usar a duração diretamente (D2), atualizando `config_test.go` e os testes de `getHtml`, que passam `time.Second`. Verificar com `go test ./internal/config ./internal/scraping`, com o caso de timeout ainda falhando dentro de cerca de 1 s
- [x] 4.4 Adicionar `validateTicker` (`^[A-Za-z0-9]+$`) e chamá-la no início de `GetStockByTicker` e de `GetReitByTicker` (D3). Verificar com testes de tabela para "", "PE%ZZ", "../x", "MX RF11" (erro, com um servidor `httptest` que falha o teste se receber requisição) e "PETR4" (válido)
- [x] 4.5 Criar as fixtures em `internal/scraping/testdata/` reproduzindo a estrutura dos XPaths atuais (D4):
  - ação completa;
  - ação com nome com hífen;
  - ação com título sem `-`;
  - ação sem preço;
  - ação com preço cercado de espaços;
  - FII completo;
  - FII sem segmento;
  - FII "não encontrado".

  Verificar abrindo cada fixture nos testes das tarefas 4.6 e 4.7
- [x] 4.6 No scraping de ação, tratar o erro de `htmlquery.Parse`, aplicar `TrimSpace` ao preço e extrair o nome com `SplitN(title, "-", 2)` (D3). Criar `stock_internal_test.go` com `httptest` cobrindo cada cenário da spec `stock-quote`:
  - encontrada (todos os campos);
  - nome com hífen;
  - título sem separador;
  - preço "1.234,56" e preço com espaços;
  - preço ausente;
  - status 404;
  - gzip válido e gzip inválido;
  - lista com ticker inválido, na ordem.

  Verificar com `go test -cover ./internal/scraping`, com `GetStockByTicker` em 100%
- [x] 4.7 No scraping de FII, usar `rs.url`, propagar o erro de `getHtml` e tratar o erro de `htmlquery.Parse` (D3). Criar `reit_internal_test.go` com `httptest` cobrindo cada cenário da spec `reit-quote`:
  - encontrado (todos os campos);
  - segmento ausente;
  - preço "1.050,75";
  - "não encontrado" com status 200;
  - status 500 (erro e nenhum FII);
  - gzip válido e gzip inválido;
  - lista com ticker inexistente, na ordem.

  Verificar com `go test -cover ./internal/scraping`, com `GetReitByTicker` em 100%

## 5. Interface de linha de comando (spec `cli`)

- [x] 5.1 Mudar `common.NewTableWriter` para receber `out io.Writer` (D6) e atualizar `internal/common/table_test.go` para verificar que a renderização vai para o writer informado, com e sem cor. Verificar com `go test ./internal/common`
- [x] 5.2 Mover `noColor`/`csv` para campos de cada struct de comando e trocar `cmd/common.go` pela função `render(t, csv)` (D7). Verificar com `grep -n "flagNoColor\|flagCsv" cmd/`, que não deve encontrar nada, e com `go build ./...`
- [x] 5.3 Converter os handlers para `RunE`, escrevendo resultados em `cmd.OutOrStdout()` e devolvendo os erros de "not found" com `SilenceUsage` (D6). Escrever a versão em `cmd.OutOrStdout()`, aplicar `MarkFlagRequired("amount")` nos três `purchase-balance` e usar `tools.SplitList` em `security` (D8). Verificar com `go build ./...`
- [x] 5.4 Mudar `core.App.Run` para `Run(args []string) int` e `main.go` para `run(args) int` + `os.Exit(run(os.Args[1:]))` (D6). Criar `internal/core/app_test.go` cobrindo `version` (0), sem argumentos (0) e `stock get` sem ticker (1), e ajustar `main_test.go` para `run([]string{})`. Verificar com `go test . ./internal/core`
- [x] 5.5 Criar `cmd/mocks_test.go` (serviços e config falsos) e `cmd/root_test.go` cobrindo a ajuda sem argumentos e `version` ("development" em stdout). Verificar com `go test ./cmd`
- [x] 5.6 Criar `cmd/stock_test.go` e `cmd/reit_test.go` cobrindo, com stdout e stderr em buffers separados:
  - `get` encontrado (todas as linhas, via `--csv`);
  - `get` não encontrado (mensagem exata em stderr, stdout vazio, erro devolvido);
  - `get` sem argumento;
  - `list` com ativos (ordem, preço com 2 casas, data `AAAA-MM-DD HH:MM:SS`, moeda "R$ BRL");
  - `list` sem resultado;
  - `purchase-balance` com rodapés SPENT AMOUNT/REMAINING AMOUNT;
  - `purchase-balance` sem resultado e sem `--amount`;
  - `--no-color` sem códigos ANSI;
  - o cenário "Opção não vaza para outra execução".

  Verificar com `go test ./cmd`
- [x] 5.7 Criar `cmd/security_test.go` cobrindo:
  - só `--stocks` (o serviço de FII recebe lista vazia);
  - ações e FIIs com a coluna TYPE;
  - sem `--amount`;
  - nenhum ativo comprado (stderr + erro).

  Verificar com `go test -cover ./cmd`, com todos os handlers em 100%

## 6. Build, CI e repositório

- [x] 6.1 Atualizar o alvo `test/unit` do `Makefile` para `go test -coverpkg=./... -coverprofile=coverage.out ./...` (D12). Verificar com `make test/unit`, que deve passar sem rede e listar `cmd/`, `core/` e `resource/` com cobertura diferente de zero
- [x] 6.2 Criar `.github/workflows/ci.yaml` (push e pull_request) com:
  - `checkout@v4` e `setup-go@v5` com `go-version-file: go.mod`;
  - `go vet ./...`;
  - `go test -race -coverpkg=./... ./...`;
  - `govulncheck` (D10).

  Verificar com `actionlint` (via `go run github.com/rhysd/actionlint/cmd/actionlint@latest`) sem erros
- [x] 6.3 Atualizar `build_and_release.yaml` (D10):
  - `go-version-file: go.mod`;
  - actions v4/v5;
  - job `test` com `needs: test` nos 3 builds;
  - upload via `gh release upload "$TAG_NAME" <arquivo> --clobber` com `GH_TOKEN`, no lugar de `upload-release-asset@v1`.

  Verificar com `actionlint` sem erros e com `grep -n "upload-release-asset\|go-version: '1.24'"`, que não deve encontrar nada
- [x] 6.4 Trocar `trader*` no `.gitignore` por `/trader`, `/trader.exe` e `/trader-*` (D11). Verificar com `git check-ignore internal/trader_x.go`, que não deve ignorá-lo, e com `git check-ignore trader trader-v1-linux-x86_64.tar.bz2`, que deve ignorar os dois
- [x] 6.5 Preencher `LICENSE` com o texto MIT e "Copyright (c) 2025 Douglas Panhota" (D11). Verificar que o arquivo não está vazio e que o GitHub/`licensee` o reconheceria como MIT (cabeçalho "MIT License" e texto padrão íntegro)
- [x] 6.6 Atualizar o `README.md` (D11):
  - seção "Desenvolvimento" com requisitos e os alvos `build/*`, `test/unit`, `test/integration`, `test/vuln` e `test/coverage`;
  - nota sobre `--amount` obrigatório e erros em stderr com código 1;
  - correção de "portifolio", "disponivel" e "o preços".

  Verificar executando cada comando documentado na seção, exceto `test/integration` se não houver rede

## 7. Verificação final

- [x] 7.1 Rodar `go vet ./...`, `go test -race ./...` e `make test/vuln` e verificar que os três passam sem avisos nem vulnerabilidades alcançáveis
- [x] 7.2 Rodar o binário localmente (`go build -o /tmp/trader . && /tmp/trader stock get`) e confirmar o código de saída 1 com erro em stderr. Confirmar também que `/tmp/trader version` imprime "development" com código 0
- [x] 7.3 Rodar `openspec validate complete-unit-tests-and-fix-bugs --strict` e verificar que a change continua válida, com cada cenário das 4 specs correspondendo a ao menos um teste dos grupos 2 a 5

## Workflow follow-up

- Rodar `make test/integration` com acesso à rede antes de arquivar, para confirmar que os XPaths ainda batem com o StatusInvest e que as dependências atualizadas não mudaram o parse.
- Após o merge, conferir a primeira execução do workflow de CI e disparar um release de teste via `workflow_dispatch`.
- Arquivar a change com `/opsx:archive`, criando as specs principais em `openspec/specs/`.
