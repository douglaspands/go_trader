# cli Specification

## Purpose
Definir a interface de linha de comando do executável `trader`: os comandos disponíveis, seus argumentos e flags, o formato das tabelas de saída, os canais de saída, os códigos de saída e as mensagens de erro exibidas ao investidor.

## Requirements

### Requirement: Comando raiz e ajuda
Executado sem subcomando, o `trader` SHALL exibir a descrição "Investor Support Tool" e a lista de comandos disponíveis: `reit`, `security`, `stock`, `version`, `completion` e `help`.

#### Scenario: Execução sem argumentos
- **WHEN** o usuário executa `trader`
- **THEN** a ajuda é exibida com os comandos `reit`, `security`, `stock` e `version`, e o código de saída é 0

### Requirement: Versão
O comando `trader version` SHALL imprimir em stdout a versão definida no build, ou "development" quando nenhuma versão foi definida.

#### Scenario: Build local
- **WHEN** o usuário executa `trader version` em um binário compilado sem definir a versão
- **THEN** a saída é "development"

### Requirement: Canais e códigos de saída
Results (tables, CSV, version, requested help) SHALL be written to stdout, with exit code 0. Warnings about tickers that failed in a list that still produced a result SHALL be written to stderr and SHALL NOT change the exit code. Every error message SHALL be written to stderr, and the process SHALL exit with code 1. This applies to usage errors, to asset(s) not found and to quote failures in a purchase balance.

#### Scenario: Sucesso
- **WHEN** a command completes, every requested asset is obtained, and the result is displayed
- **THEN** the result is on stdout, stderr is empty and the exit code is 0

#### Scenario: Ativo não encontrado
- **WHEN** a command obtains none of the requested assets
- **THEN** the error message is on stderr, stdout is empty and the exit code is 1

### Requirement: Validação de argumentos e flags
The CLI SHALL treat as a usage error an invalid number of arguments, a missing required flag, or an invalid amount. `get` requires exactly 1 ticker. `list` and the `purchase-balance` of `stock` and `reit` require at least 1 ticker. Every `purchase-balance` requires `--amount`/`-a`, which SHALL be a finite number greater than 0 and no larger than 1,000,000,000,000. An invalid amount SHALL be rejected before any request is made.

#### Scenario: get sem ticker
- **WHEN** the user runs `trader stock get`
- **THEN** the CLI writes an arguments error to stderr and exits with code 1

#### Scenario: purchase-balance sem valor
- **WHEN** the user runs `trader reit purchase-balance MXRF11` without `--amount`
- **THEN** the CLI writes a required flag error to stderr and exits with code 1

#### Scenario: Valor infinito ou indefinido
- **WHEN** the user runs `trader stock purchase-balance PETR4 --amount inf` or `--amount NaN`
- **THEN** the CLI writes an invalid amount error to stderr, makes no request, and exits with code 1 without hanging

#### Scenario: Valor zero ou negativo
- **WHEN** the user runs `trader security purchase-balance -s PETR4 --amount 0` or `--amount -100`
- **THEN** the CLI writes an invalid amount error to stderr, makes no request, and exits with code 1

#### Scenario: Valor acima do teto
- **WHEN** the user runs `trader reit purchase-balance MXRF11 --amount 1e13`
- **THEN** the CLI writes an invalid amount error to stderr, makes no request, and exits with code 1

### Requirement: Detalhe de um ativo
`trader stock get <ticker>` SHALL exibir uma tabela FIELD/VALUE com as linhas Ticker, Name, Document, Currency, Price, CapturedAt e Origin. `trader reit get <ticker>` SHALL exibir as linhas Ticker, Name, Admin, Document, Segment, Currency, Price, CapturedAt e Origin. Se o ativo não for obtido, o CLI SHALL escrever `Error: ticker "<ticker>" not found!` em stderr e encerrar com código 1.

#### Scenario: Ação encontrada
- **WHEN** o usuário executa `trader stock get ITSA3` e a cotação é obtida
- **THEN** a tabela exibe as linhas de Ticker a Origin com os dados da ação

#### Scenario: FII não encontrado
- **WHEN** o usuário executa `trader reit get XXXX00` e o FII não é obtido
- **THEN** stderr contém `Error: ticker "XXXX00" not found!` e o código de saída é 1

### Requirement: Listagem de ativos
`trader stock list` e `trader reit list` SHALL exibir uma tabela com as colunas TICKER, NAME, DOCUMENT, PRICE, CURRENCY e CAPTURED AT, com uma linha por ativo obtido, na ordem informada. Se nenhum ativo for obtido, o CLI SHALL escrever `Error: tickers not found!` em stderr e encerrar com código 1.

#### Scenario: Lista de ações
- **WHEN** o usuário executa `trader stock list ITSA3 BBDC3` e ambas as cotações são obtidas
- **THEN** a tabela exibe ITSA3 e BBDC3, nessa ordem

#### Scenario: Nenhum ativo obtido
- **WHEN** o usuário executa `trader reit list XXXX00` e nenhum FII é obtido
- **THEN** stderr contém `Error: tickers not found!` e o código de saída é 1

### Requirement: Balanceamento por tipo de ativo
`trader stock purchase-balance` e `trader reit purchase-balance` SHALL exibir uma tabela com as colunas TICKER, PRICE, COUNT, TOTAL, CURRENCY e CAPTURED AT, uma linha por ativo comprado, e dois rodapés: "SPENT AMOUNT" com a quantidade total e o valor gasto, e "REMAINING AMOUNT" com o saldo restante. Se nenhum ativo for comprado, o CLI SHALL escrever `Error: tickers not found!` em stderr e encerrar com código 1.

#### Scenario: Balanceamento de FIIs
- **WHEN** o usuário executa `trader reit purchase-balance MXRF11 HGLG11 --amount 1000` e as cotações são obtidas
- **THEN** a tabela exibe uma linha por FII comprado e os rodapés SPENT AMOUNT e REMAINING AMOUNT

#### Scenario: Valor menor que qualquer preço
- **WHEN** o usuário executa `trader reit purchase-balance HGLG11 --amount 10` e o preço de HGLG11 é maior que 10
- **THEN** stderr contém `Error: tickers not found!` e o código de saída é 1

### Requirement: Balanceamento entre ações e FIIs
`trader security purchase-balance` SHALL receber as ações em `--stocks`/`-s` e os FIIs em `--reits`/`-r`, como listas separadas por vírgula, e exibir a mesma tabela do balanceamento por tipo com a coluna adicional TYPE ("Stock" ou "REIT"). Entradas vazias SHALL ser ignoradas, de modo que uma das duas flags pode ser omitida.

#### Scenario: Somente ações
- **WHEN** o usuário executa `trader security purchase-balance --stocks ITSA3,BBDC3 --amount 1000`
- **THEN** o balanceamento considera apenas ITSA3 e BBDC3, sem consultar nenhum FII

#### Scenario: Ações e FIIs
- **WHEN** o usuário executa `trader security purchase-balance -s ITSA3 -r MXRF11 -a 1000`
- **THEN** a tabela exibe ITSA3 com TYPE "Stock" e MXRF11 com TYPE "REIT"

### Requirement: Formatação de valores
Na saída tabular, preços, totais e saldos SHALL ser exibidos com 2 casas decimais. Datas de captura SHALL usar o formato `AAAA-MM-DD HH:MM:SS`. A moeda SHALL ser exibida como "<sinal> <código>" (por exemplo, "R$ BRL").

#### Scenario: Preço e data formatados
- **WHEN** um ativo com preço 10.5 capturado em 2025-06-08 22:59:37 é exibido
- **THEN** a tabela mostra o preço "10.50", a data "2025-06-08 22:59:37" e a moeda "R$ BRL"

### Requirement: Opções de saída
Os subcomandos `get`, `list` e `purchase-balance` SHALL aceitar `--no-color`, que exibe a tabela sem cores e sem bordas externas, e `--csv`, que exibe o mesmo conteúdo em formato CSV. Sem essas flags, a tabela SHALL ser exibida colorida. Uma opção informada em um comando SHALL valer só para aquela execução.

#### Scenario: Saída CSV
- **WHEN** o usuário executa `trader stock list ITSA3 --csv`
- **THEN** a saída é CSV, com as colunas da listagem no cabeçalho e uma linha para ITSA3

#### Scenario: Opção não vaza para outra execução
- **WHEN** na mesma árvore de comandos `trader stock list ITSA3 --csv` é executado e em seguida `trader reit list MXRF11`, sem `--csv`
- **THEN** a segunda saída é uma tabela, não CSV

### Requirement: Failed tickers in a list
When `trader stock list` or `trader reit list` obtains at least one asset but some tickers fail, the CLI SHALL display the table with the assets obtained and SHALL write to stderr one warning per failed ticker, naming the ticker and the reason. The exit code SHALL be 0.

#### Scenario: Lista com uma falha
- **WHEN** the user runs `trader stock list ITSA3 XXXX0` and only ITSA3 is obtained
- **THEN** stdout has the table with ITSA3, stderr has a warning naming "XXXX0" and its reason, and the exit code is 0

### Requirement: Quote failure in a purchase balance
When the quote of any requested ticker cannot be obtained, `trader stock purchase-balance`, `trader reit purchase-balance` and `trader security purchase-balance` SHALL display no table, SHALL write to stderr an error that names every failed ticker and the reason for each, and SHALL exit with code 1.

#### Scenario: Uma ação sem cotação
- **WHEN** the user runs `trader stock purchase-balance ITSA3 XXXX0 --amount 1000` and the quote of XXXX0 cannot be obtained
- **THEN** stdout is empty, stderr names "XXXX0" and its reason, and the exit code is 1

#### Scenario: Falhas entre ações e FIIs
- **WHEN** the user runs `trader security purchase-balance -s XXXX0 -r MXRF11,XXXX00 -a 1000` and XXXX0 and XXXX00 cannot be quoted
- **THEN** stdout is empty, stderr names both "XXXX0" and "XXXX00" with their reasons, and the exit code is 1

### Requirement: Sanitized output
Before any value is written as a table or CSV cell, the CLI SHALL remove control characters from it, including escape sequences. In CSV output, a cell with text obtained from the provider that starts with `=`, `+`, `-` or `@` SHALL be prefixed with a single quote so spreadsheets do not evaluate it.

#### Scenario: Nome com sequência de escape
- **WHEN** the name of an asset obtained from the provider contains the ESC character followed by `[2J`
- **THEN** the table shows the name without the ESC character and the terminal receives no escape sequence from it

#### Scenario: Nome que parece fórmula no CSV
- **WHEN** the user runs `trader stock list ITSA3 --csv` and the name obtained for ITSA3 is `=HYPERLINK("http://x")`
- **THEN** the NAME cell in the CSV is `'=HYPERLINK("http://x")`
