package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/go-redis/redis/v8"
)

var counter int64

func main() {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Could not ping Redis: %v", err)
	}

	file, err := os.Open("entries.txt")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		counter++
		line := scanner.Text()
		hllKey := "hll_key"
		// Add the line to the HLL.
		_, err := rdb.PFAdd(context.Background(), hllKey, line).Result()
		if err != nil {
			log.Printf("Error adding element to HLL: %v", err)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	hllKey := "hll_key"
	count, err := rdb.PFCount(context.Background(), hllKey).Result()
	if err != nil {
		log.Printf("Error fetching count from HLL: %v", err)
	} else {
		fmt.Printf("Inserted Count: %d\n", counter)
		fmt.Printf("HLL Count: %d\n", count)
		fmt.Printf("Error : %f\n", float64(counter-count)/float64(counter)*100)
	}
	time.Sleep(time.Millisecond * 500)

	// Handle interruption signal to reset HLL.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		_, err := rdb.Del(context.Background(), "hll_key").Result()
		if err != nil {
			log.Printf("Error resetting HLL: %v", err)
		} else {
			fmt.Println("HLL reset")
		}
		os.Exit(0)
	}()
}
