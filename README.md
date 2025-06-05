# GoTrader [trader.exe]

É uma ferramenta projetada para auxiliar investidores de todos os níveis a tomar decisões informadas e otimizar seus portfólios. Com uma interface intuitiva e acesso a dados em tempo real, você terá as informações necessárias para acompanhar o mercado e analisar oportunidades com confiança.

**Funcionalidades Principais:**

*   **Preços de Cotas:** Acompanhe os preços de cotas de ações de diversas empresas, obtendo informações em tempo real.
*   **Preços de FIIs (Fundos de Investimento Imobiliário):** Monitore os preços dos FIIs mais populares, permitindo que você invista em imóveis de forma diversificada.

## Download

Existem versões para Windows, Linux e Mac: [Trader/Releases](https://github.com/douglaspands/go_trader/releases)

## Uso

Essa aplicação executa através do `shell` disponivel no seu sistema operacional.

> Nos exemplos abaixo usaremos o **Windows** de referencia.

1. Execute a aplicação para obter o preço da ação pelo ticker:

```sh
.\trader.exe stock get ITSA3

# Output:
 TICKER       ITSA3                                   
 NAME         ITAUSA                                  
 DOCUMENT     61.532.644/0001-15                      
 PRICE        10.95                                   
 ORIGIN       https://statusinvest.com.br/acoes/itsa3 
 CAPTURED_AT  2025-06-04 21:49:44  
```

2. Execute a aplicação para listar o preços das ações:

```sh
.\trader.exe stock list ITSA3 BBDC3 VALE3 ABEV3 PETR4 WEGE3 IGTA3 B3SA3

# Output:
 TICKER  NAME                                  DOCUMENT            PRICE  CAPTURED_AT         
 ITSA3   ITAUSA                                61.532.644/0001-15  10.95  2025-06-04 21:52:07 
 BBDC3   BRADESCO                              60.746.948/0001-12  14.21  2025-06-04 21:52:07 
 VALE3   VALE                                  33.592.510/0001-54  52.37  2025-06-04 21:52:07 
 ABEV3   AMBEV                                 07.526.557/0001-00  14.10  2025-06-04 21:52:08 
 PETR4   PETROBRAS                             33.000.167/0001-01  29.40  2025-06-04 21:52:08 
 WEGE3   WEG                                   84.429.695/0001-11  43.22  2025-06-04 21:52:08 
 IGTA3   IGUATEMI EMPRESA DE SHOPPING CENTERS  51.218.147/0001-93  33.00  2025-06-04 21:52:08 
 B3SA3   B3                                    09.346.601/0001-25  13.92  2025-06-04 21:52:09 
```

3. Execute a aplicação para obter o preço do FII pelo ticker:

```sh
.\trader.exe reit get MXRF11

# Output:
 TICKER       MXRF11                                                 
 NAME         Maxi Renda                                             
 DOCUMENT     97.521.225/0001-25                                     
 ADMIN        BTG PACTUAL SERVIÇOS FINANCEIROS S/A DTVM              
 SEGMENT      Híbrido                                                
 PRICE        9.35                                                   
 ORIGIN       https://statusinvest.com.br/fundos-imobiliarios/mxrf11 
 CAPTURED_AT  2025-06-04 21:53:18     
```

4. Execute a aplicação para listar o preços dos FIIs:

```sh
.\trader.exe reit list MXRF11 XPML11 GARE11 HGLG11 VGHF11

# Output:
 TICKER  NAME                  DOCUMENT            ADMIN                                      SEGMENT    PRICE   CAPTURED_AT         
 MXRF11  Maxi Renda            97.521.225/0001-25  BTG PACTUAL SERVIÇOS FINANCEIROS S/A DTVM  Híbrido      9.35  2025-06-04 21:58:01 
 XPML11  XP Malls              28.757.546/0001-00  XP INVESTIMENTOS CCTVM S.A.                Shoppings  103.73  2025-06-04 21:58:02 
 GARE11  Guardian Real Estate  37.295.919/0001-60  BANCO DAYCOVAL S.A.                                     8.74  2025-06-04 21:58:02 
 HGLG11  CGHG Logística        11.728.688/0001-47  PLURAL S.A. BANCO MÚLTIPLO                 Logística  156.70  2025-06-04 21:58:02 
 VGHF11  VALORA HEDGE FUND     36.771.692/0001-19  BANCO DAYCOVAL S.A.                                     7.64  2025-06-04 21:58:03 
```
