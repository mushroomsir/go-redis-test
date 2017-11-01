
## Installation

```sh
go get github.com/mushroomsir/go-redis-test
```

## Usage
```
Usage of ./main:
  -a string
    	redis password
  -c int
    	Pool size (default 100)
  -clean
    	will FlushDb if it's true
  -h string
    	redis host (default "127.0.0.1:6379")
  -n int
    	Specifies the total number of requests (default 100000)
```
### OUTPUT
```
2017/10/31 17:53:50 SET: 41860.17 requests per second
2017/10/31 17:53:52 GET: 53806.09 requests per second
2017/10/31 17:54:10 RPUSH: 5666.57 requests per second
2017/10/31 17:54:17 LRANGE100: 26492.89 requests per second
```

## Licenses

All source code is licensed under the [MIT License](https://github.com/mushroomsir/go-redis-test/blob/master/LICENSE).
