package main

import (
	"flag"
	"fmt"
	"github.com/GruppoFilippetti/mqtt-tester/tests"
	"time"
)

func main() {

	var url string
	var keepalive string
	var username string
	var password string
	var topic string

	flag.StringVar(&url, "url", "tcp://localhost:1883", "mqtt server url")
	flag.StringVar(&keepalive, "keepalive", "120", "keepalive, in seconds")
	flag.StringVar(&username, "username", "", "username")
	flag.StringVar(&password, "password", "", "password")
	flag.StringVar(&topic, "topic", "mqtt-tester", "mqtt topic")
	flag.Parse()

	// TIMEOUT 10 sec
	count := make(chan int64, 1)
	go func() {
		count <- tests.CountMessages(url, keepalive, username, password, topic, 1 * time.Second)
	}()
	select {
	case ret := <-count:
		fmt.Println(ret)
	case <-time.After(10 * time.Second):
		fmt.Println("0")
	}
}
