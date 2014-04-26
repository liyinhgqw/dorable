package main

import (
	"github.com/liyinhgqw/doracle"
	"log"
	"sync/atomic"
	"time"
)

func main() {
	arr := make([]bool, 100000000)
	client, err := doracle.NewClient([]string{"localhost:4001", "localhost:4002", "localhost:4003"})
	if err != nil {
		log.Fatalln(err)
	}
	var count int64 = 0
	t := time.Now().UnixNano()

	for i := 0; i < 100000; i++ {
		go func(i int) {
			for j := 0; j < 1000000; j++ {
				if ts, err := client.TS(); err != nil {
					log.Fatalln("ts error")
				} else {
					if arr[ts] {
						log.Fatalln("dup")
					}
					arr[ts] = true
					atomic.AddInt64(&count, 1)
					if atomic.LoadInt64(&count)%10000 == 0 {
						eps := time.Now().UnixNano() - t
						qps := atomic.LoadInt64(&count) * int64(time.Second) / eps
						log.Println("qps = ", qps)
					}
				}
			}
		}(i)
	}

	time.Sleep(100 * time.Second)
}
