# Usage

## Counter

Count messages per second

    mqtt-counter -url tcp://localhost:1883 -topic "#"
 
## Latency

Measures milliseconds between the publish of a sample message 
and receive throught subscribe  

    mqtt-pubsub-latency -url tcp://localhost:1883