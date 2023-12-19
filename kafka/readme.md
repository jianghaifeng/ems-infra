```
kafka-console-producer.sh --producer.config /opt/bitnami/kafka/config/producer.properties --bootstrap-server 127.0.0.1:9094 --topic test
kafka-console-consumer.sh --consumer.config /opt/bitnami/kafka/config/consumer.properties --bootstrap-server 127.0.0.1:9094 --topic test --from-beginning
```