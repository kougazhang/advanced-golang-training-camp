package kafka

import (
	"context"
	"log"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/webmin7761/go-school/homework/final/internal/conf"
)

type KafkaClient struct {
	topic         string
	produceClient sarama.Client
	consumerGroup sarama.ConsumerGroup
	consumerFunc  func(ctx context.Context, msg string, err error) error
}

func NewKafkaClient(conf *conf.MessageQueue) *KafkaClient {

	brokers := strings.Split(conf.Connect.Source, ",")

	version, err := sarama.ParseKafkaVersion(conf.Connect.Driver)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	config := sarama.NewConfig()
	config.Version = version
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	client, err := sarama.NewConsumerGroup(brokers, conf.Connect.Group, config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	pc := sarama.NewConfig()
	pc.Producer.Return.Successes = true
	produceClient, err := sarama.NewClient(brokers, pc)
	if err != nil {
		log.Fatalf("unable to create kafka client: %q", err)
	}

	return &KafkaClient{
		topic:         conf.Connect.Topic,
		consumerGroup: client,
		produceClient: produceClient,
	}
}

func (k *KafkaClient) Produce(msg string) error {

	producer, err := sarama.NewSyncProducerFromClient(k.produceClient)

	if err != nil {
		log.Fatalf("unable to create kafka producer: %q", err)
	}
	defer producer.Close()

	_, _, err = producer.SendMessage(&sarama.ProducerMessage{Topic: k.topic, Key: nil, Value: sarama.StringEncoder(msg)})
	if err != nil {
		log.Fatalf("unable to produce message: %q", err)
	}
	return nil
}

func (k *KafkaClient) Consume(consume func(ctx context.Context, msg string, err error) error) error {

	k.consumerFunc = consume

	consumer := Consumer{
		ready:    make(chan bool),
		MsgQueue: make(chan []byte, 1000),
		Save:     k,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumer session will need to be
		// recreated to get the new claims
		if err := k.consumerGroup.Consume(ctx, strings.Split(k.topic, ","), &consumer); err != nil {
			log.Panicf("Error from consumer: %v", err)
		}
		// check if context was cancelled, signaling that the consumer should stop
		if ctx.Err() != nil {
			return ctx.Err()
		}
		consumer.ready = make(chan bool)
	}

	<-consumer.ready

	return nil
}

func (k *KafkaClient) Close() {

}

func (k *KafkaClient) Save(msgByte []byte) {
	if len(msgByte) == 0 {
		return
	}

	err := k.consumerFunc(context.Background(), string(msgByte), nil)
	if err != nil {
		panic(err)
	}
}
