# Spec Delta

## Purpose

Definir a interface de linha de comando do executável `trader`: os comandos disponíveis, seus argumentos e flags, o formato das tabelas de saída, os canais de saída, os códigos de saída e as mensagens de erro exibidas ao investidor.

## ADDED Requirements

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
Resultados (tabelas, CSV, versão, ajuda solicitada) SHALL ser escritos em stdout, com código de saída 0. Toda mensagem de erro SHALL ser escrita em stderr, e o processo SHALL encerrar com código de saída 1. Isso vale tanto para erros de uso quanto para ativo(s) não encontrado(s).

#### Scenario: Sucesso
- **WHEN** um comando conclui e exibe seu resultado
- **THEN** o resultado está em stdout, stderr está vazio e o código de saída é 0

#### Scenario: Ativo não encontrado
- **WHEN** um comando não obtém nenhum dos ativos pedidos
- **THEN** a mensagem de erro está em stderr, stdout está vazio e o código de saída é 1

### Requirement: Validação de argumentos e flags
O CLI SHALL tratar como erro de uso um número de argumentos inválido ou uma flag obrigatória não informada. `get` exige exatamente 1 ticker. `list` e os `purchase-balance` de `stock` e `reit` exigem ao menos 1 ticker. Todo `purchase-balance` exige `--amount`/`-a`.

#### Scenario: get sem ticker
- **WHEN** o usuário executa `trader stock get`
- **THEN** o CLI escreve um erro de argumentos em stderr e encerra com código 1

#### Scenario: purchase-balance sem valor
- **WHEN** o usuário executa `trader reit purchase-balance MXRF11` sem `--amount`
- **THEN** o CLI escreve um erro de flag obrigatória em stderr e encerra com código 1

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
