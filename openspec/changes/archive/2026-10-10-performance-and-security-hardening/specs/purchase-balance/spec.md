# Spec Delta

## MODIFIED Requirements

### Requirement: Balanceamento por tickers
The system SHALL accept separate lists of stock and REIT tickers, obtain the quote of each one and compute the balance over the whole set, with stocks before REITs. If the quote of any requested ticker cannot be obtained, the system SHALL compute no balance and SHALL return an error that names every failed ticker and the reason for each failure.

#### Scenario: Ações e FIIs combinados
- **WHEN** the balance is requested for the stock "PETR4" (price 50) and the REIT "HGLG11" (price 100) with amount invested 1000
- **THEN** the result contains "PETR4" with 10 units and "HGLG11" with 5 units, 15 units in total

#### Scenario: Ticker sem cotação
- **WHEN** the balance is requested for "PETR4" and "XXXX0" and the quote of "XXXX0" cannot be obtained
- **THEN** no balance is returned, and the error names "XXXX0" and the reason of its failure

#### Scenario: Várias falhas
- **WHEN** the balance is requested for the stock "XXXX0" and the REIT "XXXX00" and neither quote can be obtained
- **THEN** no balance is returned, and the error names both "XXXX0" and "XXXX00", each with its reason
