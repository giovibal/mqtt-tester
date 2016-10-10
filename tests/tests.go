package tests

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"time"
	"sync/atomic"
)

func TracePublishDuration(url string, keepalive string, username string, password string, topic string) time.Duration {
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
		durationChannel <- time.Duration(0)
	}
	defer c.Disconnect(250)

	//subscribe to the topic /go-mqtt/sample and request messages to be delivered
	//at a maximum qos of zero, wait for the receipt to confirm the subscription
	if token := c.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		durationChannel <- time.Duration(0)
	}
	defer c.Unsubscribe(topic)

	if token := c.Publish(topic, 0, false, "check msg"); token.Wait() && token.Error() != nil {
		durationChannel <- time.Duration(0)
	}
	return <-durationChannel
}


func CountMessages(url string, keepalive string, username string, password string, topic string, duration time.Duration) int64 {
	var count int64;
	var brokerUrl = url

	count = 0

	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	opts := MQTT.NewClientOptions().AddBroker(brokerUrl)
	opts.SetClientID("mqtt-tester-count")
	opts.SetConnectTimeout(time.Duration(1 * time.Second))

	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		count = atomic.AddInt64(&count, 1)
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
		count = 0
	}
	defer c.Disconnect(250)

	//subscribe to the topic /go-mqtt/sample and request messages to be delivered
	//at a maximum qos of zero, wait for the receipt to confirm the subscription
	if token := c.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		count = 0
	}
	defer c.Unsubscribe(topic)

	select {
	case <-time.After(duration):
		return count
	}

}