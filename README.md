# GoTrader [trader.exe]

É uma ferramenta projetada para auxiliar investidores de todos os níveis a tomar decisões informadas e otimizar seus portfólios. Com uma interface intuitiva e acesso a dados em tempo real, você terá as informações necessárias para acompanhar o mercado e analisar oportunidades com confiança.

**Funcionalidades Principais:**

*   **Preços de Cotas:** Acompanhe os preços de cotas de ações de diversas empresas, obtendo informações em tempo real.
*   **Preços de FIIs (Fundos de Investimento Imobiliário):** Monitore os preços dos FIIs mais populares, permitindo que você invista em imóveis de forma diversificada.
*   **Balanceamento de Portfólio:** Dado o valor que deseja investir, e quais tickers você deseja comprar, ele calcula a quantidade de cada ativo que deve comprar.

## Download

Existem versões para Windows, Linux e Mac: [Trader/Releases](https://github.com/douglaspands/go_trader/releases)

## Uso

Essa aplicação executa através do `shell` disponivel no seu sistema operacional.

1. Execute a aplicação para obter o preço da ação pelo ticker:

```sh
# Windows
.\trader.exe stock get ITSA3
# Linux/Mac
./trader stock get ITSA3

# Output:
 FIELD       VALUE                                   
 Ticker      ITSA3                                   
 Name        ITAUSA                                  
 Document    61.532.644/0001-15                      
 Currency    R$ BRL                                  
 Price       10.87                                   
 CapturedAt  2025-06-08 22:59:37                     
 Origin      https://statusinvest.com.br/acoes/itsa3 
```

2. Execute a aplicação para listar o preços das ações:

```sh
# Windows
.\trader.exe stock list ITSA3 BBDC3 VALE3 ABEV3 PETR4 WEGE3 IGTA3 B3SA3
# Linux/Mac
./trader stock list ITSA3 BBDC3 VALE3 ABEV3 PETR4 WEGE3 IGTA3 B3SA3

# Output:
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

3. Execute a aplicação para obter o preço do FII pelo ticker:

```sh
# Windows
.\trader.exe reit get MXRF11
# Linux/Mac
./trader reit get MXRF11

# Output:
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

4. Execute a aplicação para listar o preços dos FIIs:

```sh
# Windows
.\trader.exe reit list MXRF11 XPML11 GARE11 HGLG11 VGHF11
# Linux/Mac
./trader reit list MXRF11 XPML11 GARE11 HGLG11 VGHF11

# Output:
 TICKER  NAME                  DOCUMENT             PRICE  CURRENCY  CAPTURED AT         
 MXRF11  Maxi Renda            97.521.225/0001-25    9.41  R$ BRL    2025-06-08 23:02:37 
 XPML11  XP Malls              28.757.546/0001-00  103.21  R$ BRL    2025-06-08 23:02:37 
 GARE11  Guardian Real Estate  37.295.919/0001-60    8.73  R$ BRL    2025-06-08 23:02:38 
 HGLG11  CGHG Logística        11.728.688/0001-47  156.20  R$ BRL    2025-06-08 23:02:38 
 VGHF11  VALORA HEDGE FUND     36.771.692/0001-19    7.70  R$ BRL    2025-06-08 23:02:38 
```

5. Balancear portifolio de ações:

```sh
# Windows
.\trader.exe stock purchase-balance ITSA3 BBDC3 VALE3 ABEV3 PETR4 WEGE3 IGTA3 B3SA3 --amount 1000
# Linux/Mac
./trader stock purchase-balance ITSA3 BBDC3 VALE3 ABEV3 PETR4 WEGE3 IGTA3 B3SA3 --amount 1000

# Output:
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

6. Balancear portifolio de FIIs:

```sh
# Windows
.\trader.exe reit purchase-balance MXRF11 XPML11 GARE11 HGLG11 VGHF11 --amount 1000
# Linux/Mac
./trader reit purchase-balance MXRF11 XPML11 GARE11 HGLG11 VGHF11 --amount 1000

# Output:
 TICKER   PRICE  COUNT   TOTAL  CURRENCY  CAPTURED AT         
 MXRF11    9.41     23  216.43  R$ BRL    2025-06-08 23:14:08 
 XPML11  103.21      2  206.42  R$ BRL    2025-06-08 23:14:08 
 GARE11    8.73     24  209.52  R$ BRL    2025-06-08 23:14:08 
 HGLG11  156.20      1  156.20  R$ BRL    2025-06-08 23:14:09 
 VGHF11    7.70     27  207.90  R$ BRL    2025-06-08 23:14:09 
                    77  996.47  R$ BRL    SPENT AMOUNT        
                          3.53  R$ BRL    REMAINING AMOUNT    
```

7. Balancear portifolio entre ações e FIIs:

```sh
# Windows
.\trader.exe security purchase-balance --stocks ITSA3,BBDC3,VALE3,ABEV3,PETR4,WEGE3,IGTA3,B3SA3 --reits MXRF11,XPML11,GARE11,HGLG11,VGHF11 --amount 1000
# Linux/Mac
./trader security purchase-balance --stocks ITSA3,BBDC3,VALE3,ABEV3,PETR4,WEGE3,IGTA3,B3SA3 --reits MXRF11,XPML11,GARE11,HGLG11,VGHF11 --amount 1000

# Output:
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
