package simba

import (
	"log"

	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
)

// Consumer fetches messages from kafka and calls the view function to update itself
type Consumer struct {
	doneCh   chan struct{}
	consumer *cluster.Consumer
	reducer  func(msg *sarama.ConsumerMessage) error
}

// NewConsumer constructs a startable Consumer
func NewConsumer(consumer *cluster.Consumer, reducer func(msg *sarama.ConsumerMessage) error) *Consumer {
	return &Consumer{
		consumer: consumer,
		doneCh:   make(chan struct{}, 1),
		reducer:  reducer,
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
		case err := <-c.consumer.Errors():
			log.Panicf("failure from kafka consumer: %s", err)

		case ntf := <-c.consumer.Notifications():
			log.Printf("Rebalanced: %+v\n", ntf)

		case msg := <-c.consumer.Messages():
			err := c.reducer(msg)
			if err != nil {
				log.Panicf("msg processing failed: %s", err)
				c.Stop()
			}

		case <-c.doneCh:
			log.Print("interrupt is detected")
			c.consumer.Close()
			return
		}
	}
}
