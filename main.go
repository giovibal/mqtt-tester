package main

import (
	"flag"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
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
	c1 := make(chan time.Duration, 1)
	go func() {
		//time.Sleep(time.Second * 2)
		duration := test(url, keepalive, username, password, topic)
		c1 <- duration
	}()
	//c1 := test(url, keepalive, username, password, topic)
	select {
	case res := <-c1:
		millis := (res.Seconds()*1000)
		fmt.Printf("%.0f\n",millis)
	case <-time.After(time.Second * 5):
		fmt.Println("-1")
	}
}

func test(url string, keepalive string, username string, password string, topic string) time.Duration {
	var start time.Time
	start = time.Now()

	var duration time.Duration
	duration = -1

	//var brokerUrl = strings.Join([]string{"tcp://", host, ":", port}, "")
	var brokerUrl = url

	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	opts := MQTT.NewClientOptions().AddBroker(brokerUrl)
	opts.SetClientID("mqtt-tester")

	durationChan := make(chan time.Duration, 1)
	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		durationChan <- time.Since(start)
	})
	if username != "" && password != "" {
		opts.SetUsername(username)
		opts.SetPassword(password)
	}
	if keepalive != "" {
		keepaliveDuration, err := time.ParseDuration(keepalive + "s")
		if err != nil {
			panic(err)
		}
		opts.SetKeepAlive(keepaliveDuration)
	}

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		durationChan <- time.Duration(-1)
	}

	//subscribe to the topic /go-mqtt/sample and request messages to be delivered
	//at a maximum qos of zero, wait for the receipt to confirm the subscription
	if token := c.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		durationChan <- time.Duration(-1)
	}
	token := c.Publish(topic, 0, false, "check msg")
	token.Wait()

	duration = <-durationChan

	//unsubscribe
	if token := c.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		//fmt.Println("-1")
		//os.Exit(0)
		//durationMillis = -1
		durationChan <- time.Duration(-1)
	}

	defer c.Disconnect(250)

	return duration
}
