# Usage

## Counter

Count messages per second

    mqtt-counter -url tcp://iot.eclise.org -topic "#"
 
## Latency

Measures milliseconds between the publish of a sample message 
and receive throught subscribe  

    mqtt-pubsub-latency -url tcp://iot.eclise.org