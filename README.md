# Trade App test

###Made with:

1. Go
2. Svelte


----------------------------------------
### Example for computate earn and losts:

After the three buys we have the following Lot Table:

|Lot|Shrs|Price|Extensn 
|---|----|-----|--------
1|12|100|1200
2|17|99|1683
3|3|103|309

Tot 32 at cost  3192

Sell 9 at 101
1. take 9 from lot 1
2. realized P&L on this trade = 9(101-100) = 9

Updated lot table:

|Lot|Shrs|Price|Extensn
|---|----|-----|-
1|3|100|300
2|17|99|1683
3|3|103|309

Total 23 at cost  2292

1. Sell 4 at 105
2. Take 3 from lot 1 and 1 from lot 2
3. Realized P&L on this trade = 3(105-100)+1(105-99)= 21

Updated lot table:

|Lot|Shrs|Price|Extensn
|---|----|-----|-
|1|0|100|0
|2|16|99|1584
|3|3|103|309

Total 19 at cost  1893

- mark to market of final holdings
- value =19*99 = 1881
- cost of final holdings = 0*100+16*99+3*103 = 1893
- unrealized P&L = 1881-1893 = -12
- cumulative realized P&L = 9+21=30
- total P&L = 30-12 = 18