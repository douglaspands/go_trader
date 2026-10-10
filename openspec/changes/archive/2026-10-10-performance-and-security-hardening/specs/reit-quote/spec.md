# Spec Delta

## MODIFIED Requirements

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
