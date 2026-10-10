# 📈 GoTrader

> **Assistente de Linha de Comando (CLI) para investidores da B3:** consulte cotações atualizadas e calcule o balanceamento ideal de compras para seus aportes mensais em Ações e Fundos Imobiliários (FIIs).

---

## 💡 Por que o GoTrader?

Ao realizar aportes mensais no mercado de capitais (B3), investidores frequentemente se deparam com o desafio de distribuir um orçamento fixo (por exemplo, R$ 1.000,00) entre diversos ativos selecionados. Fazer esse cálculo manualmente ou depender de planilhas complexas costuma ser demorado e gera dúvidas sobre quantas cotas inteiras comprar de cada ativo sem estourar o orçamento e minimizando o saldo que sobra parado na corretora.

O **GoTrader** resolve isso instantaneamente:
- **Automatiza a matemática do aporte:** divide o capital pretendido igualmente entre os ativos escolhidos, reserva orçamento para ativos com cotação mais alta e usa as sobras em rodadas adicionais para comprar as cotas mais baratas.
- **Evita planilhas manuais:** basta informar o valor total a ser investido e os tickers de interesse.
- **Garante decisões com preços reais:** consulta cotações atualizadas no momento da execução para garantir que seu cálculo reflita o mercado.

---

## 🎯 Para Quem É?

- **Investidores que realizam aportes periódicos:** quem investe mensalmente em uma carteira recomendada ou própria de ações e FIIs e quer saber exatamente quantas cotas comprar hoje.
- **Investidores focados em privacidade e segurança:** quem não quer vincular contas da corretora, expor senhas bancárias ou subir seus dados financeiros para serviços na nuvem.
- **Usuários de terminal e desenvolvedores:** quem busca agilidade e praticidade para integrar consultas ao seu fluxo de trabalho diário no terminal (Windows, Linux ou macOS).

---

## ✨ Funcionalidades Principais

- 🔍 **Cotações Detalhadas e Listagens:** Consulte o preço individual de uma ação ou FII, ou liste múltiplos ativos em uma tabela formatada com um único comando.
- ⚖️ **Balanceamento Automático de Compras:** Calcule a distribuição ótima do seu aporte para uma cesta de ações (`stock purchase-balance`), FIIs (`reit purchase-balance`) ou ambos simultaneamente (`security purchase-balance`).
- 💰 **Otimização de Sobra (Restante):** Algoritmo inteligente que prioriza compras igualitárias e redistribui o saldo restante em cotas adicionais dos ativos mais acessíveis, garantindo que o menor valor possível fique parado sem render.
- 📊 **Exportação para CSV:** Exporte listagens e planos de balanceamento diretamente em formato CSV (`--csv`) com proteção automática contra injeção de fórmulas para abrir com segurança no Excel ou Google Sheets.
- 🎨 **Visualização Clara e Configurável:** Tabelas com suporte a cores para facilitar a leitura no terminal e opção `--no-color` para scripts ou ambientes que não suportam cores.


---

## 🚀 Instalação e Início Rápido

O **GoTrader** é distribuído como um único arquivo executável (binário compilado), sem necessidade de instalar dependências externas, runtimes ou bibliotecas adicionais.

### 1. Baixar o executável

Acesse a página de lançamentos e faça o download da versão compatível com seu sistema:
👉 **[Baixar última versão em Releases](https://github.com/douglaspands/go_trader/releases)**

| Sistema Operacional | Arquivo | Como abrir o terminal |
| :--- | :--- | :--- |
| **Windows** | `trader.exe` | Pressione `Win + R`, digite `powershell` ou `cmd` e tecle `Enter` |
| **Linux** | `trader` | Abra seu emulador de terminal preferido (`bash`, `zsh`) |
| **macOS** | `trader` | Pressione `Cmd + Espaço`, digite `Terminal` e tecle `Enter` |

---

### 2. Como executar no seu sistema

#### 🪟 No Windows (PowerShell ou Prompt de Comando)
1. Navegue até a pasta onde você baixou o `trader.exe` (ex: pasta `Downloads`):
   ```powershell
   cd ~/Downloads
   ```
2. Verifique se o executável está funcionando:
   ```powershell
   .\trader.exe version
   ```

#### 🐧 No Linux
1. Navegue até a pasta do arquivo baixado e conceda permissão de execução:
   ```bash
   cd ~/Downloads
   chmod +x trader
   ```
2. Teste a execução:
   ```bash
   ./trader version
   ```
> **Dica:** Para usar o comando de qualquer diretório sem o prefixo `./`, mova o binário para seu PATH local:
> ```bash
> sudo mv trader /usr/local/bin/
> trader version
> ```

#### 🍎 No macOS
1. Navegue até a pasta do arquivo baixado e dê permissão de execução:
   ```bash
   cd ~/Downloads
   chmod +x trader
   ```
2. Teste a execução:
   ```bash
   ./trader version
   ```
> *(No primeiro uso no macOS, caso o sistema bloqueie a execução de desenvolvedores não identificados, acesse `Preferências do Sistema > Segurança e Privacidade` e clique em "Permitir mesmo assim" ou execute `xattr -d com.apple.quarantine trader`).*


---

## 📖 Exemplos Práticos de Uso

O **GoTrader** oferece comandos diretos para acompanhar preços e calcular seus aportes mensais. Veja a seguir os principais cenários de uso:

### 1. Consultar Cotação Individual (`get`)
Obtenha dados cadastrais e cotação em tempo real de uma ação ou fundo imobiliário:

#### Ação:
```sh
# Windows
.\trader.exe stock get ITSA3

# Linux/macOS
./trader stock get ITSA3
```
Saída:
```text
 FIELD       VALUE                                   
 Ticker      ITSA3                                   
 Name        ITAUSA                                  
 Document    61.532.644/0001-15                      
 Currency    R$ BRL                                  
 Price       10.87                                   
 CapturedAt  2025-06-08 22:59:37                     
 Origin      https://statusinvest.com.br/acoes/itsa3 
```

#### Fundo Imobiliário (FII):
```sh
# Windows
.\trader.exe reit get MXRF11

# Linux/macOS
./trader reit get MXRF11
```
Saída:
```text
 FIELD       VALUE                                                  
 Ticker      MXRF11                                                 
 Name        Maxi Renda                                             
 Admin       BTG PACTUAL SERVIÇOS FINANCEIROS S/A DTVM              
 Document    97.521.225/0001-25                                     
 Segment     Híbrido                                                
 Currency    R$ BRL                                                 
 Price       9.41                                                   
 CapturedAt  2025-06-08 23:02:13                                    
 Origin      https://statusinvest.com.br/fundos-imobiliarios/mxrf11 
```

---

### 2. Listar Cotações de Múltiplos Ativos (`list`)
Acompanhe simultaneamente os preços de uma lista de ativos em uma tabela organizada:

#### Lista de Ações:
```sh
# Windows
.\trader.exe stock list ITSA3 BBDC3 VALE3 ABEV3 PETR4 WEGE3 IGTA3 B3SA3

# Linux/macOS
./trader stock list ITSA3 BBDC3 VALE3 ABEV3 PETR4 WEGE3 IGTA3 B3SA3
```
Saída:
```text
 TICKER  NAME                                  DOCUMENT            PRICE  CURRENCY  CAPTURED AT         
  ITSA3  ITAUSA                                61.532.644/0001-15  10.87  R$ BRL    2025-06-08 23:03:44 
  BBDC3  BRADESCO                              60.746.948/0001-12  13.75  R$ BRL    2025-06-08 23:03:44 
  VALE3  VALE                                  33.592.510/0001-54  52.90  R$ BRL    2025-06-08 23:03:45 
  ABEV3  AMBEV                                 07.526.557/0001-00  14.05  R$ BRL    2025-06-08 23:03:45 
  PETR4  PETROBRAS                             33.000.167/0001-01  29.59  R$ BRL    2025-06-08 23:03:45 
  WEGE3  WEG                                   84.429.695/0001-11  42.57  R$ BRL    2025-06-08 23:03:45 
  IGTA3  IGUATEMI EMPRESA DE SHOPPING CENTERS  51.218.147/0001-93  33.00  R$ BRL    2025-06-08 23:03:46 
  B3SA3  B3                                    09.346.601/0001-25  13.55  R$ BRL    2025-06-08 23:03:46 
```

#### Lista de FIIs:
```sh
# Windows
.\trader.exe reit list MXRF11 XPML11 GARE11 HGLG11 VGHF11

# Linux/macOS
./trader reit list MXRF11 XPML11 GARE11 HGLG11 VGHF11
```
Saída:
```text
 TICKER  NAME                  DOCUMENT             PRICE  CURRENCY  CAPTURED AT         
 MXRF11  Maxi Renda            97.521.225/0001-25    9.41  R$ BRL    2025-06-08 23:02:37 
 XPML11  XP Malls              28.757.546/0001-00  103.21  R$ BRL    2025-06-08 23:02:37 
 GARE11  Guardian Real Estate  37.295.919/0001-60    8.73  R$ BRL    2025-06-08 23:02:38 
 HGLG11  CGHG Logística        11.728.688/0001-47  156.20  R$ BRL    2025-06-08 23:02:38 
 VGHF11  VALORA HEDGE FUND     36.771.692/0001-19    7.70  R$ BRL    2025-06-08 23:02:38 
```

---

### 3. Balancear Aporte Mensal em Ações (`stock purchase-balance`)
Informe o orçamento disponível com `--amount` (ou `-a`) e a lista de tickers. O GoTrader calcula a quantidade de cotas de cada empresa:

```sh
# Windows
.\trader.exe stock purchase-balance ITSA3 BBDC3 VALE3 ABEV3 PETR4 WEGE3 IGTA3 B3SA3 --amount 1000

# Linux/macOS
./trader stock purchase-balance ITSA3 BBDC3 VALE3 ABEV3 PETR4 WEGE3 IGTA3 B3SA3 --amount 1000
```
Saída:
```text
 TICKER  PRICE  COUNT   TOTAL  CURRENCY  CAPTURED AT         
  ITSA3  10.87     12  130.44  R$ BRL    2025-06-08 23:13:46 
  BBDC3  13.75     10  137.50  R$ BRL    2025-06-08 23:13:46 
  VALE3  52.90      3  158.70  R$ BRL    2025-06-08 23:13:46 
  ABEV3  14.05      9  126.45  R$ BRL    2025-06-08 23:13:46 
  PETR4  29.59      4  118.36  R$ BRL    2025-06-08 23:13:47 
  WEGE3  42.57      2   85.14  R$ BRL    2025-06-08 23:13:47 
  IGTA3  33.00      3   99.00  R$ BRL    2025-06-08 23:13:47 
  B3SA3  13.55     10  135.50  R$ BRL    2025-06-08 23:13:47 
                   53  991.09  R$ BRL    SPENT AMOUNT        
                         8.91  R$ BRL    REMAINING AMOUNT    
```
> **Destaque:** Com um orçamento de R$ 1.000,00, foram alocados R$ 991,09 em 53 cotas, restando apenas R$ 8,91 — valor inferior à cotação mais barata (ITSA3 a R$ 10,87).

---

### 4. Balancear Aporte Mensal em FIIs (`reit purchase-balance`)
Distribua seu aporte entre Fundos Imobiliários:

```sh
# Windows
.\trader.exe reit purchase-balance MXRF11 XPML11 GARE11 HGLG11 VGHF11 --amount 1000

# Linux/macOS
./trader reit purchase-balance MXRF11 XPML11 GARE11 HGLG11 VGHF11 --amount 1000
```
Saída:
```text
 TICKER   PRICE  COUNT   TOTAL  CURRENCY  CAPTURED AT         
 MXRF11    9.41     23  216.43  R$ BRL    2025-06-08 23:14:08 
 XPML11  103.21      2  206.42  R$ BRL    2025-06-08 23:14:08 
 GARE11    8.73     24  209.52  R$ BRL    2025-06-08 23:14:08 
 HGLG11  156.20      1  156.20  R$ BRL    2025-06-08 23:14:09 
 VGHF11    7.70     27  207.90  R$ BRL    2025-06-08 23:14:09 
                    77  996.47  R$ BRL    SPENT AMOUNT        
                          3.53  R$ BRL    REMAINING AMOUNT    
```

---

### 5. Balanceamento Combinado de Ações e FIIs (`security purchase-balance`)
Faça um aporte unificado entre ações e FIIs utilizando o comando `security purchase-balance`. Separe cada lista por vírgula com as flags `--stocks` (ou `-s`) e `--reits` (ou `-r`):

```sh
# Windows
.\trader.exe security purchase-balance --stocks ITSA3,BBDC3,VALE3,ABEV3,PETR4,WEGE3,IGTA3,B3SA3 --reits MXRF11,XPML11,GARE11,HGLG11,VGHF11 --amount 1000

# Linux/macOS
./trader security purchase-balance --stocks ITSA3,BBDC3,VALE3,ABEV3,PETR4,WEGE3,IGTA3,B3SA3 --reits MXRF11,XPML11,GARE11,HGLG11,VGHF11 --amount 1000
```
Saída:
```text
 TICKER  TYPE    PRICE  COUNT   TOTAL  CURRENCY  CAPTURED AT         
  ITSA3  Stock   10.84      7   75.88  R$ BRL    2025-06-10 09:35:19 
  BBDC3  Stock   13.69      5   68.45  R$ BRL    2025-06-10 09:35:19 
  VALE3  Stock   53.29      1   53.29  R$ BRL    2025-06-10 09:35:19 
  ABEV3  Stock   13.97      5   69.85  R$ BRL    2025-06-10 09:35:19 
  PETR4  Stock   29.17      2   58.34  R$ BRL    2025-06-10 09:35:20 
  WEGE3  Stock   42.84      1   42.84  R$ BRL    2025-06-10 09:35:20 
  IGTA3  Stock   33.00      2   66.00  R$ BRL    2025-06-10 09:35:20 
  B3SA3  Stock   13.16      5   65.80  R$ BRL    2025-06-10 09:35:21 
 MXRF11  REIT     9.31      8   74.48  R$ BRL    2025-06-10 09:35:21 
 GARE11  REIT     8.74      9   78.66  R$ BRL    2025-06-10 09:35:21 
 VGHF11  REIT     7.57     11   83.27  R$ BRL    2025-06-10 09:35:22 
 XPML11  REIT   102.10      1  102.10  R$ BRL    2025-06-10 09:35:21 
 HGLG11  REIT   155.27      1  155.27  R$ BRL    2025-06-10 09:35:21 
                           58  994.23  R$ BRL    SPENT AMOUNT        
                                 5.77  R$ BRL    REMAINING AMOUNT    
```

---

### 6. Exportação para CSV e Opções de Formatação

#### Exportação para Planilhas (`--csv`):
Gere arquivos CSV prontos para importação no Excel, LibreOffice Calc ou Google Planilhas adicionando `--csv` ao comando:

```sh
# Exportar listagem de ações para arquivo
./trader stock list ITSA3 BBDC3 VALE3 --csv > carteira.csv

# Exportar plano de aporte para arquivo
./trader stock purchase-balance ITSA3 BBDC3 VALE3 --amount 1000 --csv > aporte.csv
```
Exemplo de conteúdo gerado no arquivo CSV:
```csv
TICKER,PRICE,COUNT,TOTAL,CURRENCY,CAPTURED AT
ITSA3,10.87,31,336.97,R$ BRL,2025-06-08 23:13:46
BBDC3,13.75,24,330.00,R$ BRL,2025-06-08 23:13:46
VALE3,52.90,6,317.40,R$ BRL,2025-06-08 23:13:46
,,61,984.37,R$ BRL,SPENT AMOUNT
,,,15.63,R$ BRL,REMAINING AMOUNT
```
*(Nota de segurança: o GoTrader aplica sanitização automática em campos de texto no CSV para impedir injeção maliciosa de fórmulas em planilhas).*

#### Saída sem Cores (`--no-color`):
Para uso em pipelines, automações ou terminais que não suportam formatação ANSI, utilize a flag `--no-color`:
```sh
./trader stock get ITSA3 --no-color
```

---

## ❓ Perguntas Frequentes (FAQ)

### 1. De onde vêm as cotações e informações exibidas?
As cotações e dados cadastrais (razão social, CNPJ, segmento, administrador) são consultados diretamente do portal financeiro [StatusInvest](https://statusinvest.com.br) no momento da execução. Não é necessário possuir cadastro nem chave de API.

### 2. O GoTrader é seguro? Ele pede login ou acessa minha corretora?
**Sim, é 100% seguro.** O GoTrader nunca solicita credenciais bancárias, logins ou senhas de corretoras, nem se comunica com sistemas de negociação. Ele é uma ferramenta analítica que roda inteiramente no seu computador, calculando a distribuição matemática das cotas para que você mesmo efetue as compras na sua corretora de preferência.

### 3. Preciso de conexão com a internet para usá-lo?
Sim, é necessária conexão com a internet durante a execução para que o GoTrader consulte os preços públicos atualizados. Os cálculos matemáticos de balanceamento e redistribuição de sobras são executados instantaneamente e localmente na sua máquina.

### 4. O que acontece se uma cotação falhar ou um ticker estiver incorreto?
- **Nas listagens (`list`):** O GoTrader exibe na tabela todos os ativos encontrados com sucesso e emite um aviso (`Warning: <ticker>: <motivo>`) para os que falharam, mantendo o código de saída em `0`.
- **Nos balanceamentos (`purchase-balance`):** Para sua segurança financeira, se **qualquer** cotação falhar, o cálculo é cancelado por completo. O GoTrader não arrisca dividir seu dinheiro apenas entre os ativos que deram certo: ele avisa exatamente quais tickers falharam e encerra com código de saída `1`.

---

## 🛠️ Para Desenvolvedores

Esta seção é destinada a desenvolvedores e mantenedores que desejam compilar, testar, estender ou contribuir com o **GoTrader**.

### Arquitetura do Projeto

O projeto é escrito em [Go](https://go.dev/) e organizado em camadas bem definidas:
- **`main.go`**: Ponto de entrada que inicializa a aplicação chamando `internal/core`.
- **`internal/core`**: Monta o contêiner de injeção de dependências e liga os comandos CLI aos serviços.
- **`cmd/`**: Implementação dos comandos e flags da CLI utilizando a biblioteca [Cobra](https://github.com/spf13/cobra) (`stock`, `reit`, `security`, `version`).
- **`internal/service`**: Regras de negócio e algoritmo de balanceamento de compras (`PurchaseBalanceService`, `StockService`, `ReitService`).
- **`internal/scraping`**: Cliente HTTP e parser de HTML para consulta pública ao StatusInvest.
- **`internal/resource`**: Modelos e estruturas de dados de domínio (ações, FIIs, falhas de cotação).
- **`openspec/specs/`**: Especificações executáveis de comportamento do sistema no formato OpenSpec.

---

### Compilação Local

**Requisitos:** [Go](https://go.dev/dl/) na versão declarada em `go.mod` e utilitário `make`.

```sh
# Compilar binário para Linux (gera ./trader)
make build/linux

# Compilar binário para Windows (gera ./trader.exe)
make build/windows

# Compilar binário para macOS (gera ./trader)
make build/darwin
```

---

### Testes e Qualidade

O repositório adota testes automatizados rigorosos, separando testes unitários rápidos e reproduzíveis de testes de integração com a rede:

```sh
# Testes unitários (offline, sem acesso à rede), com cobertura entre pacotes (gera coverage.out)
make test/unit

# Relatório de cobertura de testes em HTML
make test/coverage

# Testes de integração (consultam o StatusInvest ao vivo - necessitam de conexão com a rede)
make test/integration

# Verificação de vulnerabilidades conhecidas em dependências (govulncheck)
make test/vuln
```

> **Nota sobre os testes de scraping:** Os testes unitários utilizam amostras HTML estáticas salvas em `internal/scraping/testdata/` para garantir execução offline instantânea e independente de instabilidades externas. Os testes de integração (usando a build tag `integration`) servem para detectar alterações no layout HTML do portal StatusInvest.

---

### Erros, Códigos de Saída e Sanitização

O GoTrader segue padrões estritos para fluxos de dados, códigos de saída e proteção contra injeção:

- **Canais de Saída:**
  - Resultados normais (tabelas, CSV e versão) saem em `stdout`, com código de saída `0`.
  - Mensagens de erro saem em `stderr`, com código de saída `1`.
- **Validação de Argumentos:**
  - Comandos `get` exigem exatamente 1 ticker. Comandos `list` e `purchase-balance` exigem pelo menos 1 ticker.
  - Comandos `purchase-balance` exigem a flag `--amount` (`-a`). Sem ela, a execução falha com erro de flag obrigatória.
  - O valor de `--amount` deve ser um número finito, maior que `0` e no máximo `1000000000000` (1e12). Entradas como `inf`, `NaN`, `0` ou números negativos são rejeitadas imediatamente antes de qualquer requisição.
  - Tickers aceitam apenas letras e dígitos, com limite de 12 caracteres (ex: `PETR4`, `MXRF11`). Tickers inválidos são rejeitados antes de requisições à rede.
- **Tratamento de Falhas:**
  - **Listagens (`list`):** se ao menos um ativo for obtido com sucesso, a tabela é gerada em `stdout` e cada ativo que falhou gera um aviso `Warning: <ticker>: <motivo>` em `stderr`. O código de saída continua `0`. Se nenhum ativo for obtido, o erro sai em `stderr` com código `1`.
  - **Balanceamento (`purchase-balance`):** se qualquer ticker falhar na cotação, nenhuma tabela é exibida. O erro detalha cada ticker que falhou em `stderr` e encerra com código `1`, impedindo balanceamentos incorretos com ativos faltantes.
- **Sanitização de Dados e CSV:**
  - Todos os textos obtidos da web passam por remoção de caracteres de controle e sequências de escape ANSI antes de serem exibidos.
  - Na exportação `--csv`, qualquer célula de texto que inicie com `=`, `+`, `-` ou `@` recebe o prefixo `'` para impedir que softwares de planilha executem fórmulas maliciosas (proteção contra CSV Injection).

#### Exemplos de tratamento de erros no terminal:
```sh
# Ativo inexistente em consulta individual
./trader stock get XXXX0 > /dev/null; echo $?
# Error: ticker "XXXX0" not found!   (stderr)
# 1

# Valor de aporte inválido
./trader stock purchase-balance PETR4 --amount inf
# Error: invalid amount +Inf: it must be a finite number greater than 0 and no larger than 1000000000000   (stderr)

# Falha parcial em listagem (aviso em stderr, tabela em stdout, saída 0)
./trader stock list ITSA3 XXXX0
# (tabela com ITSA3 em stdout)
# Warning: XXXX0: <motivo>   (stderr)

# Falha em balanceamento (interrompe execução, saída 1)
./trader stock purchase-balance ITSA3 XXXX0 --amount 1000
# Error: could not get the quote of 1 ticker(s):   (stderr)
# XXXX0: <motivo>
```

---

### Desenvolvimento com Agentes de IA

O repositório possui uma infraestrutura de guardrails locais para agentes de codificação autônomos (como Claude Code e Antigravity `agy`):

- **Limites de Segurança:** Operações dentro do repositório operam de forma autônoma; operações fora do repositório, comandos opacos, `git push` e alterações na configuração do harness exigem confirmação do usuário.
- **Acesso a Segredos:** Acesso a arquivos confidenciais (`.env*`, chaves privadas, credenciais em `~/.ssh` ou `~/.aws`) e comandos destrutivos (`rm -rf`) são bloqueados.
- **Ferramentas Nativas:** Agentes devem utilizar ferramentas nativas do harness para leitura e escrita de arquivos, sendo proibido o uso de leitores/escritores shell (`cat`, `sed`, `tee`, redirecionamentos) e interpretadores externos (`python`, `node`).
- **Chamadas de Shell:** Cada invocação de comando de terminal deve conter uma única instrução (sem `&&`, `;`, `|`).
- **Suíte de Conformidade:**
  ```sh
  # Executada sob demanda para validar o comportamento dos guardrails:
  scripts/harness-test.sh
  ```
- **Arquivos de Instruções:**
  - [AGENTS.md](AGENTS.md): Diretrizes gerais compartilhadas por todos os agentes.
  - [CLAUDE.md](CLAUDE.md): Configurações específicas para Claude Code.
  - [GEMINI.md](GEMINI.md): Configurações específicas para Antigravity CLI.
  - [docs/harness/local-guardrails.md](docs/harness/local-guardrails.md): Detalhamento da arquitetura de segurança e modelos de ameaça.
  - [.agents/agy-settings.example.json](.agents/agy-settings.example.json): Modelo de permissões globais para o Antigravity.
- **Convenção de Idioma:** Este arquivo `README.md` é mantido em Português do Brasil. Todo o restante do repositório (código, comentários, especificações OpenSpec, documentações internas e mensagens de commit) é redigido em Inglês.

