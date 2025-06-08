# GoTrader [trader.exe]

É uma ferramenta projetada para auxiliar investidores de todos os níveis a tomar decisões informadas e otimizar seus portfólios. Com uma interface intuitiva e acesso a dados em tempo real, você terá as informações necessárias para acompanhar o mercado e analisar oportunidades com confiança.

**Funcionalidades Principais:**

*   **Preços de Cotas:** Acompanhe os preços de cotas de ações de diversas empresas, obtendo informações em tempo real.
*   **Preços de FIIs (Fundos de Investimento Imobiliário):** Monitore os preços dos FIIs mais populares, permitindo que você invista em imóveis de forma diversificada.

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
 Currency    BRL R$                                  
 Price       10.87                                   
 CapturedAt  2025-06-08 20:04:59                     
 Origin      https://statusinvest.com.br/acoes/itsa3 
```

2. Execute a aplicação para listar o preços das ações:

```sh
# Windows
.\trader.exe stock list ITSA3 BBDC3 VALE3 ABEV3 PETR4 WEGE3 IGTA3 B3SA3
# Linux/Mac
./trader stock list ITSA3 BBDC3 VALE3 ABEV3 PETR4 WEGE3 IGTA3 B3SA3

# Output:
 TICKER  NAME                                  DOCUMENT            CURRENCY  PRICE  CAPTURED AT         
 ITSA3   ITAUSA                                61.532.644/0001-15  BRL R$    10.87  2025-06-08 20:05:19 
 BBDC3   BRADESCO                              60.746.948/0001-12  BRL R$    13.75  2025-06-08 20:05:19 
 VALE3   VALE                                  33.592.510/0001-54  BRL R$    52.90  2025-06-08 20:05:19 
 ABEV3   AMBEV                                 07.526.557/0001-00  BRL R$    14.05  2025-06-08 20:05:19 
 PETR4   PETROBRAS                             33.000.167/0001-01  BRL R$    29.59  2025-06-08 20:05:20 
 WEGE3   WEG                                   84.429.695/0001-11  BRL R$    42.57  2025-06-08 20:05:20 
 IGTA3   IGUATEMI EMPRESA DE SHOPPING CENTERS  51.218.147/0001-93  BRL R$    33.00  2025-06-08 20:05:20 
 B3SA3   B3                                    09.346.601/0001-25  BRL R$    13.55  2025-06-08 20:05:20 
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
 Currency    BRL R$                                                 
 Price       9.41                                                   
 CapturedAt  2025-06-08 20:04:04                                    
 Origin      https://statusinvest.com.br/fundos-imobiliarios/mxrf11    
```

4. Execute a aplicação para listar o preços dos FIIs:

```sh
# Windows
.\trader.exe reit list MXRF11 XPML11 GARE11 HGLG11 VGHF11
# Linux/Mac
./trader reit list MXRF11 XPML11 GARE11 HGLG11 VGHF11

# Output:
 TICKER  NAME                  DOCUMENT            CURRENCY   PRICE  CAPTURED AT         
 MXRF11  Maxi Renda            97.521.225/0001-25  BRL R$      9.41  2025-06-08 20:04:25 
 XPML11  XP Malls              28.757.546/0001-00  BRL R$    103.21  2025-06-08 20:04:25 
 GARE11  Guardian Real Estate  37.295.919/0001-60  BRL R$      8.73  2025-06-08 20:04:25 
 HGLG11  CGHG Logística        11.728.688/0001-47  BRL R$    156.20  2025-06-08 20:04:25 
 VGHF11  VALORA HEDGE FUND     36.771.692/0001-19  BRL R$      7.70  2025-06-08 20:04:26 
```
