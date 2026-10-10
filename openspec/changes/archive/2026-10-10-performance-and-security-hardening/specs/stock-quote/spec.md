# Spec Delta

## MODIFIED Requirements

### Requirement: Validação do ticker de ação
The system SHALL reject with an error, without making any request, a ticker that is empty, longer than 12 characters, or contains characters other than ASCII letters and digits.

#### Scenario: Ticker vazio
- **WHEN** the stock "" is queried
- **THEN** the system returns an invalid ticker error and no request is made

#### Scenario: Ticker com caractere inválido
- **WHEN** the stock "PE%ZZ" or "../x" is queried
- **THEN** the system returns an invalid ticker error, without a process failure and without a request

#### Scenario: Ticker longo demais
- **WHEN** a stock ticker with 13 characters is queried
- **THEN** the system returns an invalid ticker error and no request is made

### Requirement: Requisição ao provedor
Requests to StatusInvest SHALL use browser headers (a User-Agent chosen at random from a list, pt-BR language, Google Referer) and a 60-second timeout. They SHALL advertise only content encodings the system decodes, and a compressed response SHALL be decompressed before parsing. A response body larger than 10 MiB SHALL be an error. Redirects SHALL be followed only to `https` URLs on the provider's host. No more than 4 requests to the provider SHALL be in flight at the same time within one command run.

#### Scenario: Resposta gzip
- **WHEN** the provider responds with `Content-Encoding: gzip`
- **THEN** the content is decompressed before it is parsed

#### Scenario: Tempo limite
- **WHEN** the provider does not respond within the configured timeout
- **THEN** the query ends with an error

#### Scenario: Resposta grande demais
- **WHEN** the provider responds with status 200 and a body larger than 10 MiB
- **THEN** the query ends with an error and no stock is returned

#### Scenario: Redirect para outro host
- **WHEN** the provider responds with a redirect to a host other than the provider's
- **THEN** the redirect is not followed and the query ends with an error

#### Scenario: Redirect para http
- **WHEN** the provider responds with a redirect to an `http` URL on its own host
- **THEN** the redirect is not followed and the query ends with an error

#### Scenario: Limite de requisições simultâneas
- **WHEN** 10 stock tickers are queried in one command run
- **THEN** at no moment are more than 4 requests to the provider in flight

### Requirement: Consulta de ações em lista
The system SHALL query every ticker in the list, concurrently within the request limit, and return the stocks obtained successfully in the order the tickers were given. For each ticker that fails, the system SHALL report the ticker and the reason for the failure. A failed ticker SHALL NOT stop the others.

#### Scenario: Lista com ticker inválido
- **WHEN** the list "PETR4", "XXXX0", "VALE3" is queried and "XXXX0" fails
- **THEN** the system returns "PETR4" and "VALE3", in that order, and reports "XXXX0" with the reason of its failure

#### Scenario: Ordem preservada com respostas fora de ordem
- **WHEN** the list "A1", "B1", "C1" is queried and the provider answers "C1" first and "A1" last
- **THEN** the system returns "A1", "B1" and "C1", in that order
