# purchase-balance Specification

## Purpose
Distribuir um valor de investimento entre um conjunto de ativos (ações e/ou FIIs), calculando quantas cotas inteiras comprar de cada um, o total gasto e o saldo que sobra.

## Requirements

### Requirement: Resultado do balanceamento
O sistema SHALL devolver, para um valor investido e uma lista de ativos com preço, a lista de ativos comprados com a quantidade inteira de cada um, o valor investido informado, a quantidade total de cotas, o valor gasto (soma de preço x quantidade) e o saldo restante (valor investido menos valor gasto).

#### Scenario: Totais consistentes
- **WHEN** o balanceamento é calculado para ativos de preço 10 e 10 com valor investido 100
- **THEN** cada ativo recebe 5 cotas, a quantidade total é 10, o valor gasto é 100 e o saldo restante é 0

### Requirement: Divisão igualitária inicial
O sistema SHALL dividir o valor investido igualmente entre os ativos válidos e comprar, de cada um, o número inteiro de cotas que cabe na sua parte. Quando a parte de cada ativo for menor que o preço do ativo mais caro, o sistema SHALL comprar uma cota a menos de cada ativo que comporte ao menos duas, reservando saldo para os ativos mais caros.

#### Scenario: Ativo caro comprado com o saldo reservado
- **WHEN** o balanceamento é calculado para "High" a 600 e "Low" a 100 com valor investido 1000
- **THEN** o resultado contém "Low" com 4 cotas e "High" com 1 cota, e o saldo restante é 0

### Requirement: Ativos caros
Um ativo cujo preço é maior que a sua parte da divisão inicial SHALL receber 1 cota se o saldo restante, depois da divisão inicial, for suficiente. Caso contrário, SHALL ficar fora do resultado.

#### Scenario: Ativo caro que não cabe no saldo
- **WHEN** o balanceamento é calculado para "High" a 900 e "Low" a 100 com valor investido 1000
- **THEN** o cálculo termina sem erro, "High" não aparece no resultado, "Low" recebe 10 cotas e o saldo restante é 0

### Requirement: Redistribuição do saldo restante
Depois da divisão inicial e da compra dos ativos caros, o sistema SHALL usar o saldo restante para comprar cotas adicionais dos ativos já presentes no resultado, percorrendo-os do mais barato para o mais caro, em rodadas, até que o saldo restante seja menor que o preço de todos os ativos válidos.

#### Scenario: Sobra menor que qualquer preço
- **WHEN** um balanceamento termina com qualquer conjunto de ativos válidos
- **THEN** o saldo restante é maior ou igual a 0 e menor que o preço de cada ativo válido informado

### Requirement: Ordem do resultado
O resultado SHALL listar primeiro os ativos comprados na divisão inicial, na ordem em que foram informados, e depois os ativos caros comprados com o saldo restante, também na ordem informada.

#### Scenario: Ativo caro aparece no final
- **WHEN** o balanceamento é calculado para "High" a 600 e "Low" a 100, nessa ordem, com valor investido 1000
- **THEN** o resultado lista "Low" antes de "High"

### Requirement: Valor insuficiente ou lista vazia
O sistema SHALL devolver um resultado sem ativos, com o saldo restante igual ao valor investido, quando o valor investido for menor que o preço do ativo mais barato ou quando não houver ativos válidos.

#### Scenario: Valor menor que o ativo mais barato
- **WHEN** o balanceamento é calculado para ativos de preço 100 e 200 com valor investido 50
- **THEN** o resultado não contém ativos, a quantidade total é 0 e o saldo restante é 50

#### Scenario: Lista vazia
- **WHEN** o balanceamento é calculado sem nenhum ativo com valor investido 100
- **THEN** o resultado não contém ativos e o saldo restante é 100

### Requirement: Ativos com preço inválido
O sistema SHALL desconsiderar ativos com preço menor ou igual a zero: eles não entram na divisão, não aparecem no resultado e o cálculo SHALL sempre terminar.

#### Scenario: Ativo com preço zero
- **WHEN** o balanceamento é calculado para "ZERO" a 0 e "A" a 10 com valor investido 100
- **THEN** o cálculo termina, "ZERO" não aparece no resultado, "A" recebe 10 cotas e o saldo restante é 0

### Requirement: Balanceamento por tickers
O sistema SHALL aceitar listas separadas de tickers de ações e de FIIs, obter a cotação de cada um e calcular o balanceamento sobre o conjunto, com as ações antes dos FIIs. Tickers cuja cotação não puder ser obtida SHALL ser ignorados.

#### Scenario: Ações e FIIs combinados
- **WHEN** o balanceamento é pedido para a ação "PETR4" (preço 50) e o FII "HGLG11" (preço 100) com valor investido 1000
- **THEN** o resultado contém "PETR4" com 10 cotas e "HGLG11" com 5 cotas, totalizando 15 cotas

#### Scenario: Ticker sem cotação
- **WHEN** o balanceamento é pedido para tickers cuja cotação não pôde ser obtida
- **THEN** esses tickers não aparecem no resultado e o cálculo usa apenas os demais
