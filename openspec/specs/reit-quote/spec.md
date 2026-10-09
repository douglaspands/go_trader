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
O sistema SHALL rejeitar com erro, sem fazer nenhuma requisição, um ticker vazio ou que contenha caracteres fora de letras ASCII e dígitos.

#### Scenario: Ticker vazio
- **WHEN** o FII "" é consultado
- **THEN** o sistema devolve um erro de ticker inválido, nenhum FII e nenhuma requisição é feita

#### Scenario: Ticker com caractere inválido
- **WHEN** o FII "MX RF11" ou "%zz" é consultado
- **THEN** o sistema devolve um erro de ticker inválido, sem falha do processo e sem requisição

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
As requisições de FIIs SHALL seguir as mesmas regras das requisições de ações: cabeçalhos de navegador, aceitação de resposta gzip e tempo limite de 60 segundos.

#### Scenario: Resposta gzip
- **WHEN** o provedor responde à consulta de um FII com `Content-Encoding: gzip`
- **THEN** o conteúdo é descompactado antes de ser interpretado

### Requirement: Consulta de FIIs em lista
O sistema SHALL consultar cada ticker da lista, na ordem informada, e devolver apenas os FIIs obtidos com sucesso. Tickers com falha SHALL ser omitidos sem interromper os demais.

#### Scenario: Lista com ticker inexistente
- **WHEN** a lista "MXRF11", "XXXX00", "HGLG11" é consultada e "XXXX00" não existe
- **THEN** o sistema devolve "MXRF11" e "HGLG11", nessa ordem
