# ABC Trading

Is a simple platform to get trading statistics. 

It consist of two endpoints.
- `/stats` - lets you get stats for given `symbol` passed as symbol URL param and 10^k last tranding prices passed as `k` URL param.
- `/addBatch` - lets you too add batch of trading prices of given symbol. Ordered from oldes to newest.

## How to run 

To run the application use make target `make run` by default it will expose port 8080. You can
configure port or how much data should be colleted for symbol by environment variables. Look at `internal/config/config.go` to 
see what can be changed.

## Stats 

Example request:
``` bash
curl -X GET  http://localhost:8080/stats/?symbol=AAPL&k=3
```

Will return statistics for last 1000 trading prices for symbol "AAPL".

Here is example response
``` bash 
λ ~> curl -X GET  http://localhost:8080/stats/\?symbol\=A\&k\=3
{"min":0,"max":9,"last":3,"avg":3.923076923076923,"var":7.609467455621303}
```

Response consist of:
- min - Lowest price
- max - Highest price
- last - Latest price seen
- avg - Average price
- var - Variance of prices

## AddBatch

Example request 
``` bash
curl -X POST -d '{"symbol":"AAPL","values":[1,2,3,4,5,6,7,8,9,0]}' http://localhost:8080/add_batch/
```

It will add 10 new prices for symbol AAPL. After that next GET request will include statistics for this values.


