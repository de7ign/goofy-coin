# Goofy Coin
A simple cryptocurrency

## Rules
- Only Goofy can create new coins and it belongs to Goofy
- Whoever owns a coin can spend/pass it to other participants

## Data Structure
- Goofy creates a coin
```
 ___________________________
|    signed by PK(goofy)    |
|___________________________|
| Createcoin[uniqueCoinID]  |
|___________________________|
```

- if Goofy pass the coin to alice then,
```
 ___________________________
|    signed by PK(goofy)    | 
|___________________________|
|   Pay to PK(alice) : H()--|---+
|___________________________|   |
                                |
                                |
                 _______________|
                |    ___________________________
                +-->|    signed by PK(goofy)    |
                    |___________________________|
                    | Createcoin[uniqueCoinID]  |
                    |___________________________|
```

- Similarly if Alice pass the coin to Bob then,
```
 ___________________________
|    signed by PK(alice)    | 
|___________________________|
|   Pay to PK(bob) : H()----|----+
|___________________________|    |
                                 |
                                 |
               __________________|
              |      ___________________________
              +---->|    signed by PK(goofy)    | 
                    |___________________________|
                    |   Pay to PK(alice) : H()--|---+
                    |___________________________|   |
                                                    |
                                                    |
                                 ___________________|
                                |        ___________________________
                                +------>|    signed by PK(goofy)    |
                                        |___________________________|
                                        | Createcoin[uniqueCoinID]  |
                                        |___________________________|
```


## Problems
- Central authority
  
   only Goofy can make coins

- Double Spending attack
   
   Participants can do double spending


## Author
Nihal Murmu - [nihalmurmu](https://github.com/nihalmurmu)

## License
This project is licensed under the MIT License - see the [LICENSE.md](https://github.com/nihalmurmu/goofy-coin/blob/master/LICENSE) file for details
 