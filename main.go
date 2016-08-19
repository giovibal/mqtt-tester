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

	var brokerUrl = url

	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	opts := MQTT.NewClientOptions().AddBroker(brokerUrl)
	opts.SetClientID("mqtt-tester")
	opts.SetConnectTimeout(time.Duration(1 * time.Second))

	durationChannel := make(chan time.Duration, 1)
	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		durationChannel <- time.Since(start)
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
		durationChannel <- time.Duration(-1)
	}
	defer c.Disconnect(250)

	//subscribe to the topic /go-mqtt/sample and request messages to be delivered
	//at a maximum qos of zero, wait for the receipt to confirm the subscription
	if token := c.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		durationChannel <- time.Duration(-1)
	}
	defer c.Unsubscribe(topic);

	if token := c.Publish(topic, 0, false, "check msg"); token.Wait() && token.Error() != nil {
		durationChannel <- time.Duration(-1)
	}
	return <- durationChannel
}
