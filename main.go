package main

import (
	"flag"
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
		PoolSize: poolSize,
		Addr:     host,
		Password: password,
	})
	if clean {
		client.FlushDb()
	}
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
						log.Println(err)
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
	log.Printf("SET: %.2f requests per second\n", float64(requests)/time.Since(start).Seconds())
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
						log.Println(err)
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
	log.Printf("GET: %.2f requests per second\n", float64(requests)/time.Since(start).Seconds())
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
						log.Println(err)
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
	log.Printf("RPUSH: %.2f requests per second\n", float64(requests)/time.Since(start).Seconds())
	// for i := 0; i < requests; i++ {
	// 	err := client.Expire("redis:list:benchmark:test:"+strconv.Itoa(i), time.Minute).Err()
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// }
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
					err := client.LRange("redis:list:benchmark:test:"+m, 1, 100).Err()
					if err != nil {
						log.Println(err)
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
	log.Printf("LRANGE100: %.2f requests per second\n", float64(requests)/time.Since(start).Seconds())
}
