package main

import (
	"flag"
	"github.com/GruppoFilippetti/mqtt-tester/tests"
	"time"
	"os"
	"syscall"
	"fmt"
	"os/signal"
)

func main() {

	var url string
	var keepalive string
	var username string
	var password string
	var topic string
	var threshold time.Duration

	flag.StringVar(&url, "url", "tcp://localhost:1883", "mqtt server url")
	flag.StringVar(&keepalive, "keepalive", "15", "keepalive, in seconds")
	flag.StringVar(&username, "username", "", "username")
	flag.StringVar(&password, "password", "", "password")
	flag.StringVar(&topic, "topic", "#", "mqtt topic")
	flag.DurationVar(&threshold, "threshold", 1*time.Second, "duration in seconds. " +
		"Es. '300ms', '-1.5h' or '2h45m'). " +
		"Valid time units are 'ns', 'us' (or 'µs'), 'ms', 's', 'm', 'h'.")
	flag.Parse()

	err := tests.TraceSnapshotVsReceivedTimestamp(url, keepalive, username, password, topic,
		func(topic string, uri string, creationTime time.Time, receivedTime time.Time) {
			var duration time.Duration
			if receivedTime.After(creationTime) {
				duration = receivedTime.Sub(creationTime)
			} else {
				duration = creationTime.Sub(receivedTime)
			}

			if duration > threshold {
				traceEvent(duration, creationTime, receivedTime, topic, uri)
			}
		})
	if err!=nil {
		panic(err)
	}
	// capture ctrl+c
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	select {
	case <-c:
		fmt.Println("Shutting down ...")
		os.Exit(0)
	}
}

func traceEvent(duration time.Duration, creationTime time.Time, receivedTime time.Time,  topic string, uri string) {
	fmt.Printf("ALERT - Received after %v seconds (Creation date: %v -> Received date: %v) on topic %v (uri %v)\n",
		duration.Seconds(), creationTime, receivedTime, topic, uri)
}
