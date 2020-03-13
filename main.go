package main

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	attempts := 0
	for ;; {
		_, err := client.Ping().Result()
		if err == nil {
			break;
		}

		time.Sleep(time.Second)
		attempts += 1

		if attempts > 3 {
			log.Fatal(err)
		}
	}

	http.HandleFunc("/hit", func(writer http.ResponseWriter, request *http.Request) {
		v, err := client.GetSet("hits", 0).Result()
		if err != nil {
			_, _ = writer.Write([]byte("hits: 0"))
			return
		}

		hits, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			_, _ = writer.Write([]byte("Error: " + err.Error()))
		}
		hits++

		err = client.Set("hits", hits, 0).Err()
		if err != nil {
			_, _ = writer.Write([]byte("Error: " + err.Error()))
		}

		log.Println("hits: ", hits)
		_, _ = writer.Write([]byte(fmt.Sprintf("Hits: %d", hits)))
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
