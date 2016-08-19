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

	// TIMEOUT 5 sec
	duration := make(chan time.Duration, 1)
	go func() {
		duration <- tests.TracePublishDuration(url, keepalive, username, password, topic)
	}()
	select {
	case ret := <-duration:
		millis := (ret.Seconds() * 1000)
		fmt.Printf("%.0f\n", millis)
	case <-time.After(time.Second * 5):
		fmt.Println("-1")
	}

}
