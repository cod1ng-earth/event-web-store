package simba

import (
	"log"

	"github.com/Shopify/sarama"
)

// Consumer fetches messages from kafka and calls the view function to update itself
type Consumer struct {
	doneCh    chan struct{}
	consumer  sarama.Consumer
	partition sarama.PartitionConsumer
	reducer   func(msg *sarama.ConsumerMessage) error
}

// NewConsumer constructs a startable Consumer
func NewConsumer(consumer sarama.Consumer, partition sarama.PartitionConsumer, reducer func(msg *sarama.ConsumerMessage) error) *Consumer {
	return &Consumer{
		consumer:  consumer,
		partition: partition,
		doneCh:    make(chan struct{}, 1),
		reducer:   reducer,
	}
}

// Stop ends eventloop
func (c *Consumer) Stop() {
	c.doneCh <- struct{}{}
}

// Start listens for events from kafka
func (c *Consumer) Start() {

	for {
		select {
		case err := <-c.partition.Errors():
			log.Printf("failure from kafka consumer: %s", err)

		case msg := <-c.partition.Messages():
			err := c.reducer(msg)
			if err != nil {
				log.Panicf("msg processing failed: %s", err)
				c.Stop()
			}

		case <-c.doneCh:
			log.Print("interrupt is detected")
			if err := c.consumer.Close(); err != nil {
				log.Panicf("failed to close kafka consumer: %s", err)
			}
			return
		}
	}
}
