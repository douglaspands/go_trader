# Spec Delta

## Purpose

Obter os dados cadastrais e a cotação atual de ações brasileiras a partir das páginas públicas do StatusInvest, por ticker individual ou por lista de tickers.

## ADDED Requirements

### Requirement: Consulta de ação por ticker
O sistema SHALL consultar a página `<base>/acoes/<ticker em minúsculas>` do StatusInvest e devolver uma ação com: ticker em maiúsculas, nome, CNPJ, preço atual, tipo "Stock", moeda BRL (sinal "R$", descrição "Brazilian Real"), URL consultada como origem e o instante da captura.

#### Scenario: Ação encontrada
- **WHEN** a ação "petr4" é consultada e a página responde com status 200 e os dados da empresa
- **THEN** o sistema devolve a ação com ticker "PETR4", tipo "Stock", moeda "R$ BRL" e origem `<base>/acoes/petr4`

### Requirement: Validação do ticker de ação
O sistema SHALL rejeitar com erro, sem fazer nenhuma requisição, um ticker vazio ou que contenha caracteres fora de letras ASCII e dígitos.

#### Scenario: Ticker vazio
- **WHEN** a ação "" é consultada
- **THEN** o sistema devolve um erro de ticker inválido e nenhuma requisição é feita

#### Scenario: Ticker com caractere inválido
- **WHEN** a ação "PE%ZZ" ou "../x" é consultada
- **THEN** o sistema devolve um erro de ticker inválido, sem falha do processo e sem requisição

### Requirement: Nome da empresa
O nome da ação SHALL ser todo o texto após o primeiro `-` do atributo de título do cabeçalho principal da página, sem espaços nas extremidades. Se o cabeçalho não existir ou o título não contiver `-`, o nome SHALL ficar vazio e a consulta SHALL ser concluída sem falha.

#### Scenario: Título no formato padrão
- **WHEN** o título do cabeçalho é "PETR4 - PETROBRAS"
- **THEN** o nome da ação é "PETROBRAS"

#### Scenario: Nome com hífen
- **WHEN** o título do cabeçalho é "XYZW3 - FOO-BAR S.A."
- **THEN** o nome da ação é "FOO-BAR S.A."

#### Scenario: Título sem separador
- **WHEN** o título do cabeçalho é "PETROBRAS"
- **THEN** a consulta é concluída sem erro e o nome da ação fica vazio

### Requirement: Interpretação do preço
O preço SHALL ser lido do valor atual do ativo exibido na página, no formato brasileiro (ponto como separador de milhar e vírgula como decimal), desconsiderando espaços nas extremidades. Se o valor não existir ou não puder ser interpretado, o preço SHALL ser 0.

#### Scenario: Preço com milhar
- **WHEN** a página exibe o valor atual "1.234,56"
- **THEN** o preço da ação é 1234.56

#### Scenario: Preço com espaços
- **WHEN** a página exibe o valor atual " 38,15 " com espaços e quebras de linha ao redor
- **THEN** o preço da ação é 38.15

#### Scenario: Preço ausente
- **WHEN** a página não exibe o valor atual do ativo
- **THEN** o preço da ação é 0

### Requirement: Falha na consulta de ação
O sistema SHALL devolver erro, e nenhuma ação, quando a página responder com status diferente de 200, quando a requisição falhar (falha de rede ou tempo esgotado), quando a resposta declarar compactação gzip com conteúdo inválido ou quando o conteúdo não puder ser interpretado como HTML. Nenhuma dessas situações SHALL encerrar o processo de forma anormal.

#### Scenario: Ticker inexistente
- **WHEN** a ação "XXXX0" é consultada e a página responde com status 404
- **THEN** o sistema devolve um erro que menciona o status e a URL consultada

#### Scenario: Gzip inválido
- **WHEN** a página responde com `Content-Encoding: gzip` e um corpo que não é gzip válido
- **THEN** o sistema devolve um erro e nenhuma ação

### Requirement: Requisição ao provedor
As requisições ao StatusInvest SHALL usar cabeçalhos de navegador (User-Agent escolhido aleatoriamente de uma lista, idioma pt-BR, Referer do Google), SHALL aceitar resposta compactada com gzip e SHALL ter tempo limite de 60 segundos.

#### Scenario: Resposta gzip
- **WHEN** o provedor responde com `Content-Encoding: gzip`
- **THEN** o conteúdo é descompactado antes de ser interpretado

#### Scenario: Tempo limite
- **WHEN** o provedor não responde dentro do tempo limite configurado
- **THEN** a consulta é encerrada com erro

### Requirement: Consulta de ações em lista
O sistema SHALL consultar cada ticker da lista, na ordem informada, e devolver apenas as ações obtidas com sucesso. Tickers com falha SHALL ser omitidos sem interromper os demais.

#### Scenario: Lista com ticker inválido
- **WHEN** a lista "PETR4", "XXXX0", "VALE3" é consultada e "XXXX0" falha
- **THEN** o sistema devolve "PETR4" e "VALE3", nessa ordem
