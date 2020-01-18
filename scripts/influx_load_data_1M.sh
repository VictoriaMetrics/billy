curl -X POST http://localhost:8086/query?q=create%20database%20benchmark

./billy -startdate=2019-01-01 -enddate=2019-12-31 -startkey=1 -endkey=1000000 -sink='http://localhost:8086/write?db=benchmark' -format=influx -blocks-per-request=10
