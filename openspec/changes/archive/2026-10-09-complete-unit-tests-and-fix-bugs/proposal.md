# Proposal

## Why

O projeto precisa de uma base saneada antes do próximo ciclo de desenvolvimento. Hoje:

- Não há nenhuma especificação registrada.
- A cobertura de testes é desigual, e os testes de scraping dependem da rede real.
- A exploração confirmou defeitos que derrubam ou travam o programa:
  - panic e loop infinito no balanceamento de compra;
  - panic com resposta gzip inválida e com ticker contendo caractere inválido;
  - FII "fantasma" com preço 0 quando a rede falha ou o ticker é vazio.
- O `govulncheck` aponta 9 vulnerabilidades conhecidas, alcançáveis pelo código, nas bibliotecas que interpretam o HTML de terceiros.
- O pipeline de release usa uma versão do Go diferente da declarada, não roda testes e falha quando é disparado manualmente.

Esta change registra o comportamento esperado, corrige todos esses problemas e deixa testes determinísticos e CI protegendo o resultado.

## What Changes

**Especificação**
- Registrar a especificação de base do CLI: cotação de ações, cotação de FIIs, balanceamento de compra e interface de linha de comando.

**Balanceamento de compra**
- Corrigir o panic quando há ativos que não puderam ser comprados.
- Corrigir o loop infinito quando há ativo com preço menor ou igual a zero (esses ativos passam a ser descartados).

**Scraping (ações e FIIs)**
- O scraping de FII passa a propagar erros de rede/HTTP, como o de ações, e a usar a URL base configurada na instância.
- Tickers vazios ou com caracteres fora de `[A-Za-z0-9]` são rejeitados com erro, sem requisição.
- Resposta gzip inválida, falha na criação da requisição e falha no parse do HTML passam a gerar erro em vez de panic.
- O corpo da resposta passa a ser fechado também quando o status é diferente de 200.
- O texto do preço é normalizado (espaços nas extremidades) antes da conversão.
- O nome da ação passa a ser todo o texto após o primeiro `-` do título, sem truncar nomes com hífen e sem panic quando não há `-`.
- O timeout de scraping passa a ser uma duração real (60 s), sem a multiplicação implícita por segundos.

**CLI**
- **BREAKING** (leve): `--amount` passa a ser obrigatório em todos os `purchase-balance`. Antes, a ausência resultava em "tickers not found".
- **BREAKING** (leve): as mensagens `Error: ...` de ticker não encontrado passam a sair em stderr, com código de saída 1. Antes saíam em stdout, com código 0.
- `security purchase-balance` ignora entradas vazias de `--stocks`/`--reits`.
- As flags `--no-color`/`--csv` deixam de ser variáveis globais e passam a pertencer a cada comando.
- A saída dos comandos passa a ser injetável, e a definição do código de saída sai de dentro da aplicação para o `main`, permitindo testá-la.

**Dependências e segurança**
- Atualizar `golang.org/x/net` para >= v0.55.0 e `github.com/antchfx/xpath` para >= v1.3.6, junto com `htmlquery` e as demais dependências diretas.
- Rodar `govulncheck` no CI.

**Testes**
- Completar os testes unitários de `cmd`, `core`, `resource`, `service`, `scraping` (HTML de fixture servido localmente) e `tools`.
- Os testes contra a rede real passam a ser testes de integração, executados sob demanda.
- Medir a cobertura entre pacotes (`-coverpkg=./...`).

**Repositório e CI**
- Workflow de release:
  - usar a versão do Go do `go.mod`;
  - rodar os testes antes dos builds;
  - publicar os artefatos com `gh release upload`, o que funciona também no disparo manual;
  - atualizar as actions.
- Novo workflow de CI em push e pull request: `go vet`, testes unitários e `govulncheck`.
- `.gitignore`: trocar o padrão `trader*`, que ignora qualquer arquivo do repositório começando com "trader", por padrões ancorados na raiz.
- Remover a função não utilizada `tools.ToSnakeCase` e seus testes.
- Preencher `LICENSE` com a licença MIT.
- README: adicionar uma seção de desenvolvimento (build, testes, integração, cobertura), corrigir erros de digitação e documentar o comportamento de erro do CLI.

## Capabilities

### New Capabilities
- `stock-quote`: Obtenção dos dados e da cotação de ações por ticker (individual e em lista) a partir do StatusInvest, incluindo validação de ticker e tratamento de falhas da requisição.
- `reit-quote`: Obtenção dos dados e da cotação de FIIs por ticker (individual e em lista) a partir do StatusInvest, incluindo validação de ticker, detecção de ticker inexistente e tratamento de falhas da requisição.
- `purchase-balance`: Algoritmo que distribui um valor de investimento entre ativos, calculando a quantidade de cada um, o total gasto e o saldo restante.
- `cli`: Comandos, argumentos, flags, formato de saída, canais de saída e códigos de saída do executável `trader` (`stock`, `reit`, `security`, `version`).

### Modified Capabilities
<!-- Nenhuma: não há specs existentes em openspec/specs/. -->

## Impact

- **Código:**
  - `internal/service/purchasebalance.go`
  - `internal/scraping/{common,stock,reit}.go`
  - `internal/config/config.go`
  - `internal/common/table.go`
  - `internal/core/app.go`
  - `internal/tools/string.go`
  - `cmd/*.go`
  - `main.go`
- **Testes:** novos ou ampliados em todos os pacotes. Fixtures HTML em `internal/scraping/testdata/`. Os testes de rede atuais ficam sob a build tag `integration`.
- **Build e CI:**
  - `Makefile` (cobertura entre pacotes, alvos de integração e de vulnerabilidades)
  - `.github/workflows/build_and_release.yaml`
  - novo `.github/workflows/ci.yaml`
  - `.gitignore`
- **Dependências:** atualização de `golang.org/x/net`, `github.com/antchfx/xpath` e `github.com/antchfx/htmlquery`, e revisão das demais dependências diretas. Nenhuma dependência nova.
- **Documentação e legal:** `README.md`, `LICENSE` (MIT).
- **Usuário:** `--amount` passa a ser obrigatório, e erros de ticker não encontrado passam a sair em stderr com código 1.
