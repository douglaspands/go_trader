# reit-quote Specification

## Purpose
Obter os dados cadastrais e a cotação atual de fundos de investimento imobiliário (FIIs) a partir das páginas públicas do StatusInvest, por ticker individual ou por lista de tickers.

## Requirements

### Requirement: Consulta de FII por ticker
O sistema SHALL consultar a página `<base>/fundos-imobiliarios/<ticker em minúsculas>` do StatusInvest e devolver um FII com: ticker em maiúsculas, nome, administrador, CNPJ, segmento, preço atual, tipo "REIT", moeda BRL (sinal "R$", descrição "Brazilian Real"), URL consultada como origem e o instante da captura.

#### Scenario: FII encontrado
- **WHEN** o FII "mxrf11" é consultado e a página responde com status 200 e os dados do fundo
- **THEN** o sistema devolve o FII com ticker "MXRF11", tipo "REIT", moeda "R$ BRL", nome, administrador, CNPJ e segmento lidos da página, e origem `<base>/fundos-imobiliarios/mxrf11`

### Requirement: Validação do ticker de FII
The system SHALL reject with an error, without making any request, a ticker that is empty, longer than 12 characters, or contains characters other than ASCII letters and digits.

#### Scenario: Ticker vazio
- **WHEN** the REIT "" is queried
- **THEN** the system returns an invalid ticker error, no REIT, and no request is made

#### Scenario: Ticker com caractere inválido
- **WHEN** the REIT "MX RF11" or "%zz" is queried
- **THEN** the system returns an invalid ticker error, without a process failure and without a request

#### Scenario: Ticker longo demais
- **WHEN** a REIT ticker with 13 characters is queried
- **THEN** the system returns an invalid ticker error and no request is made

### Requirement: Campos ausentes do FII
Se algum campo (nome, administrador, CNPJ, segmento ou preço) não existir na página, o sistema SHALL devolver o FII com esse campo vazio (ou preço 0), sem falha. O preço SHALL ser interpretado no formato brasileiro (ponto como separador de milhar e vírgula como decimal), desconsiderando espaços nas extremidades.

#### Scenario: Página sem segmento
- **WHEN** a página do FII não exibe o segmento
- **THEN** o FII é devolvido com o segmento vazio e os demais campos preenchidos

#### Scenario: Preço no formato brasileiro
- **WHEN** a página exibe o valor atual "1.050,75"
- **THEN** o preço do FII é 1050.75

### Requirement: FII inexistente
O sistema SHALL devolver erro, e nenhum FII, quando a página exibir a mensagem "Não encontramos o que você está procurando". O provedor responde com status 200 nesse caso, então a detecção SHALL ser feita pelo conteúdo.

#### Scenario: Ticker inexistente
- **WHEN** o FII "XXXX00" é consultado e a página responde com status 200 e a mensagem de conteúdo não encontrado
- **THEN** o sistema devolve um erro que menciona a URL consultada

### Requirement: Falha na consulta de FII
O sistema SHALL devolver erro, e nenhum FII, quando a página responder com status diferente de 200, quando a requisição falhar (falha de rede ou tempo esgotado), quando a resposta declarar compactação gzip com conteúdo inválido ou quando o conteúdo não puder ser interpretado como HTML. Um FII com dados vazios e preço 0 nunca SHALL ser devolvido como resultado de uma falha de requisição.

#### Scenario: Provedor indisponível
- **WHEN** o FII "MXRF11" é consultado e a requisição falha ou responde com status 500
- **THEN** o sistema devolve um erro e nenhum FII

#### Scenario: Gzip inválido
- **WHEN** a página do FII responde com `Content-Encoding: gzip` e um corpo que não é gzip válido
- **THEN** o sistema devolve um erro, nenhum FII, e o processo não é encerrado de forma anormal

### Requirement: Requisição ao provedor de FIIs
REIT requests SHALL follow the same rules as stock requests: browser headers, decompression of compressed responses, a 60-second timeout, the 10 MiB body limit and the redirect policy. Stock and REIT requests SHALL share the limit of 4 requests in flight within one command run.

#### Scenario: Resposta gzip
- **WHEN** the provider responds to a REIT query with `Content-Encoding: gzip`
- **THEN** the content is decompressed before it is parsed

#### Scenario: Resposta grande demais
- **WHEN** the provider responds to a REIT query with status 200 and a body larger than 10 MiB
- **THEN** the query ends with an error and no REIT is returned

#### Scenario: Limite compartilhado com ações
- **WHEN** one command run queries 3 stocks and 3 REITs
- **THEN** at no moment are more than 4 requests to the provider in flight

### Requirement: Consulta de FIIs em lista
The system SHALL query every ticker in the list, concurrently within the request limit, and return the REITs obtained successfully in the order the tickers were given. For each ticker that fails, the system SHALL report the ticker and the reason for the failure. A failed ticker SHALL NOT stop the others.

#### Scenario: Lista com ticker inexistente
- **WHEN** the list "MXRF11", "XXXX00", "HGLG11" is queried and "XXXX00" does not exist
- **THEN** the system returns "MXRF11" and "HGLG11", in that order, and reports "XXXX00" with the reason of its failure
