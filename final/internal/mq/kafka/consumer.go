package kafka

import (
	"github.com/Shopify/sarama"
)

type MsgSave interface {
	Save(msgByte []byte)
	Close()
}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready    chan bool
	MsgQueue chan []byte
	Save     MsgSave
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	consumer.writeFile()
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	consumer.Save.Close()
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		//log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		consumer.MsgQueue <- message.Value
		session.MarkMessage(message, "")
	}

	return nil
}

func (consumer *Consumer) writeFile() {
	go func() {
		for {
			consumer.Save.Save(<-consumer.MsgQueue)
		}
	}()
}
