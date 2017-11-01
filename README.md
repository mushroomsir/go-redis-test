
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
Parallel,Operation,QPS
100,SET,35676.05
100,GET,46210.70
100,RPUSH100,3897.68
100,LRANGE100,26759.41
```

## Licenses

All source code is licensed under the [MIT License](https://github.com/mushroomsir/go-redis-test/blob/master/LICENSE).
