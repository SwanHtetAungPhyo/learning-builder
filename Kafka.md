### Kafka Imagei


/opt/homebrew/opt/kafka/bin/kafka-consumer-groups --bootstrap-server localhost:9092 --describe --group test-group


### Confluent local 

confluent local kafka start


docker run -d -p 4222:4222 --name nats-server nats