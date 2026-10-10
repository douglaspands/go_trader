# Spec Delta

## MODIFIED Requirements

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

## ADDED Requirements

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
