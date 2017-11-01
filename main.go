package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"gopkg.in/redis.v5"
)

var (
	host     string
	password string
	poolSize int
	requests int
	clean    bool
)

func init() {
	flag.StringVar(&host, "h", "127.0.0.1:6379", "redis host")
	flag.StringVar(&password, "a", "", "redis password")
	flag.IntVar(&poolSize, "c", 100, "Pool size")
	flag.IntVar(&requests, "n", 100000, "Specifies the total number of requests")
	flag.BoolVar(&clean, "clean", false, "will FlushDb if it's true")

}
func main() {

	flag.Parse()
	client := redis.NewClient(&redis.Options{
		PoolSize:    poolSize,
		Addr:        host,
		Password:    password,
		ReadTimeout: 5 * time.Second,
	})
	if clean {
		err := client.FlushDb().Err()
		if err != nil {
			log.Println("FlushDb: ", err.Error())
		} else {
			log.Println("FlushDb Completed")
		}
	}
	fmt.Printf("Parallel,Operation,QPS\n")
	set(client)
	get(client)
	rpush(client)
	lrange(client)

}
func set(client *redis.Client) {
	msg := make(chan string, requests)
	for i := 0; i < requests; i++ {
		msg <- strconv.Itoa(i)
	}
	wg := sync.WaitGroup{}
	wg.Add(poolSize)
	start := time.Now()
	for i := 0; i < poolSize; i++ {
		go func() {
			for {
				select {
				case m := <-msg:
					err := client.Set("redis:benchmark:test:"+m, strings.Repeat("a", 2048), time.Minute).Err()
					if err != nil {
						log.Println("Set: ", err)
					}
				default:
					goto label
				}
			}
		label:
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Printf("%v,SET,%.2f\n", poolSize, float64(requests)/time.Since(start).Seconds())
}
func get(client *redis.Client) {
	msg := make(chan string, requests)
	for i := 0; i < requests; i++ {
		msg <- strconv.Itoa(i)
	}
	wg := sync.WaitGroup{}
	wg.Add(poolSize)
	start := time.Now()
	for i := 0; i < poolSize; i++ {
		go func() {
			for {
				select {
				case m := <-msg:
					err := client.Get("redis:benchmark:test:" + m).Err()
					if err != nil {
						log.Println("Get: ", err)
					}
				default:
					goto label
				}
			}
		label:
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Printf("%v,GET,%.2f\n", poolSize, float64(requests)/time.Since(start).Seconds())
}

func rpush(client *redis.Client) {
	msg := make(chan string, requests)
	for i := 0; i < requests; i++ {
		msg <- strconv.Itoa(i)
	}
	wg := sync.WaitGroup{}
	wg.Add(poolSize)
	start := time.Now()
	for i := 0; i < poolSize; i++ {
		go func() {
			for {
				select {
				case m := <-msg:
					rpushs := make([]interface{}, 100)
					for i := 0; i < 100; i++ {
						rpushs = append(rpushs, strings.Repeat("a", 128))
					}
					err := client.RPush("redis:list:benchmark:test:"+m, rpushs...).Err()
					if err != nil {
						log.Println("RPush100: ", err)
					}
				default:
					goto label
				}
			}
		label:
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Printf("%v,RPUSH100,%.2f\n", poolSize, float64(requests)/time.Since(start).Seconds())
}

func lrange(client *redis.Client) {
	msg := make(chan string, requests)
	for i := 0; i < requests; i++ {
		msg <- strconv.Itoa(i)
	}
	wg := sync.WaitGroup{}
	wg.Add(poolSize)
	start := time.Now()
	for i := 0; i < poolSize; i++ {
		go func() {
			for {
				select {
				case m := <-msg:
					result, err := client.LRange("redis:list:benchmark:test:"+m, 0, 99).Result()
					if err != nil {
						log.Println("LRANGE100: ", err)
					} else {
						if len(result) != 100 {
							log.Println("LRANGE100 not 100, redis:list:benchmark:test:" + m)
						}
					}
				default:
					goto label
				}
			}
		label:
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Printf("%v,LRANGE100,%.2f\n", poolSize, float64(requests)/time.Since(start).Seconds())
}
