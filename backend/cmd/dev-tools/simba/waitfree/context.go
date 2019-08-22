package {{ .Name }}

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Shopify/sarama"
)

const (
	Topic     = "{{ .Name }}"
	Partition = 0
)

type context struct {
	doneCh    chan struct{}
	client    sarama.Client
	consumer  sarama.Consumer
	partition sarama.PartitionConsumer

	batchOffset int64

{{ if eq .Lock "exclusive" }}
	model  *model
	lock   sync.Mutex
	writes chan *sarama.ConsumerMessage
{{ else if eq .Lock "parallel" }}
	model  *model
	lock   sync.RWMutex
	writes chan *sarama.ConsumerMessage
{{ else if eq .Lock "wait-free" }}
	readerAChanged *sync.Cond
	readerBChanged *sync.Cond
	aIsReading     bool
	readersA       int32
	readersB       int32
	modelA         *model
	modelB         *model

	writes     chan *sarama.ConsumerMessage
	writesRedo chan *sarama.ConsumerMessage
{{ end }}

	offset        int64
	offsetChanged *sync.Cond
}

func NewContext(brokers *[]string, cfg *sarama.Config) context {

	client, err := sarama.NewClient(*brokers, cfg)
	if err != nil {
		log.Panicf("failed to setup kafka client: %s", err)
	}
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		log.Panicf("failed to setup kafka consumer: %s", err)
	}
	partition, err := consumer.ConsumePartition(Topic, 0, 0)
	if err != nil {
		log.Panicf("failed to setup kafka partition: %s", err)
	}

	batchOffset, err := client.GetOffset(Topic, Partition, sarama.OffsetNewest)
	if err != nil {
		log.Printf("failed to get last offset for topic %v partition %v", Topic, Partition)
	}
	batchOffset--

	context := context{
		doneCh:    make(chan struct{}, 1),
		client:    client,
		consumer:  consumer,
		partition: partition,

		batchOffset: batchOffset,

	{{ if eq .Lock "exclusive" }}
		model:  newModel(),
		lock:   &sync.Mutex{},
		writes: make(chan *sarama.ConsumerMessage, 32768),
	{{ else if eq .Lock "parallel" }}
		model:  newModel(),
		lock:   &sync.RWMutex{},
		writes: make(chan *sarama.ConsumerMessage, 32768),
	{{ else if eq .Lock "wait-free" }}
		readerAChanged: sync.NewCond(&sync.Mutex{}),
		readerBChanged: sync.NewCond(&sync.Mutex{}),
		aIsReading:     true,
		readersA:		0,
		readersB:		0,
		modelA:         newModel(),
		modelB:         newModel(),

		writes:     make(chan *sarama.ConsumerMessage, 32768),
		writesRedo: make(chan *sarama.ConsumerMessage, 32768),
	{{ end }}

		offset:        0,
		offsetChanged: sync.NewCond(&sync.Mutex{}),
	}
	return context
}

func (c *context) Stop() {
	c.doneCh <- struct{}{}
}

func (c *context) AwaitLastOffset() {
	c.offsetChanged.L.Lock()
	for c.offset < c.batchOffset {
		c.offsetChanged.Wait()
	}
	c.offsetChanged.L.Unlock()
}

func (c *context) updateLoop() {

	var offset int64
	var start time.Time

	for {
	{{ if eq .Lock "exclusive" }}
		applyChange(msg, c.model, c)
	{{ else if eq .Lock "parallel" }}
		applyChange(msg, c.model, c)
	{{ else if eq .Lock "wait-free" }}
		writeDelay := 0 * time.Second

		model := c.modelA
		readers := &c.readersA
		waiter := c.readerAChanged
		if c.aIsReading {
			model = c.modelB
			readers = &c.readersB
			waiter = c.readerBChanged
		}

		waiter.L.Lock()
		if atomic.LoadInt32(readers) != 0 {
			waiter.Wait()
		}
		waiter.L.Unlock()

		for len(c.writesRedo) > 0 {
			msg, ok := <-c.writesRedo
			if ok {
				applyChange(msg, model, c)
			}
		}

		msg, ok := <-c.writes
		if !ok {
			return
		}
		c.writesRedo <- msg
		start = time.Now()
		applyChange(msg, model, c)
		offset = msg.Offset
		writeDelay += time.Since(start)

		for writeDelay < 10*time.Millisecond && len(c.writes) > 0 && len(c.writesRedo) < 32768 {
			msg, ok := <-c.writes
			if !ok {
				return
			}
			c.writesRedo <- msg
			start := time.Now()
			applyChange(msg, model, c)
			offset = msg.Offset
			writeDelay += time.Since(start)
		}

		c.aIsReading = !c.aIsReading
		c.offset = offset
		c.offsetChanged.Broadcast()
	{{ end }}
	}

}

func applyChange(msg *sarama.ConsumerMessage, m *model, c *context) {

{{ if eq .Lock "exclusive" }}
	c.lock.Lock()
	defer c.lock.Unlock()
{{ else if eq .Lock "parallel" }}
	c.lock.Lock()
	defer c.lock.Unlock()
{{ end }}

{{ if .Batch }}
	if msg.Offset < c.batchOffset {
		batchUpdateModel(msg, m)
	} else if msg.Offset == c.batchOffset {
		batchUpdateModel(msg, m)
		batchFinalizeModel(m)
	} else {
		updateModel(msg, m)
	}
{{ else }}
	updateModel(msg, m)
{{ end }}
}

func (c *context) Start() {

	log.Printf("starting context %v", Topic)

	go c.updateLoop()

	for {
		select {
		case err := <-c.partition.Errors():
			log.Printf("failure from kafka consumer: %s", err)

		case msg := <-c.partition.Messages():
			c.writes <- msg

		case <-c.doneCh:
			close(c.writes)
			log.Print("interrupt is detected")
			if err := c.partition.Close(); err != nil {
				log.Panicf("failed to close kafka partition: %s", err)
			}
			if err := c.consumer.Close(); err != nil {
				log.Panicf("failed to close kafka consumer: %s", err)
			}
			if err := c.client.Close(); err != nil {
				log.Panicf("failed to close kafka client: %s", err)
			}
			return
		}
	}
}

func (c *context) read() (*model, func()) {

	if msg.Offset < c.batchOffset {
		c.AwaitLastOffset()
	}

{{ if eq .Lock "exclusive" }}
	c.lock.Lock()
	return model, c.lock.Unlock
{{ else if eq .Lock "parallel" }}
	c.lock.RLock()
	return model, c.lock.Unlock
{{ else if eq .Lock "wait-free" }}
	atomic.AddInt32(&c.readersA, 1)
	atomic.AddInt32(&c.readersB, 1)

	if c.aIsReading {
		atomic.AddInt32(&c.readersB, -1)
		c.readerBChanged.Signal()
		return c.modelA, func() {
			atomic.AddInt32(&c.readersA, -1)
			c.readerAChanged.Signal()
		}
	}

	atomic.AddInt32(&c.readersA, -1)
	c.readerAChanged.Signal()
	return c.modelB, func() {
		atomic.AddInt32(&c.readersB, -1)
		c.readerBChanged.Signal()
	}
{{ end }}
}

{{ if .Batch }}
func batchUpdateModel(msg *sarama.ConsumerMessage, model *model) error {
	cc := {{ .Name | title }}Messages{}
	err := proto.Unmarshal(msg.Value, &cc)
	if err != nil {
		return fmt.Errorf("failed to unmarshal kafka massage %s/%d: %v", Topic, msg.Offset, err)
	}

	defer func() { offset = msg.Offset }()

	switch x := cc.Get{{ .Name | title }}Message().(type) {

	{{ range .MessageNames }}
	case *{{ $.Name | title }}Messages_{{ . | title }}:
		if err := batchUpdateModel{{ . | title }}(cc.Get{{ . | title }}(), offset); err != nil {
			return err
		}
	{{ end }}

	case nil:
		panic(fmt.Sprintf("context message is empty"))

	default:
		panic(fmt.Sprintf("unexpected type %T in oneof", x))
	}
	return nil
}
{{ end }}

func updateModel(msg *sarama.ConsumerMessage, model *model) error {
	cc := CheckoutMessages{}
	err := proto.Unmarshal(msg.Value, &cc)
	if err != nil {
		return fmt.Errorf("failed to unmarshal kafka massage %s/%d: %v", Topic, msg.Offset, err)
	}

	defer func() { offset = msg.Offset }()

	switch x := cc.Get{{ .Name | title }}Message().(type) {

	{{ range .MessageNames }}
	case *{{ $.Name | title }}Messages_{{ . | title }}:
		if err := updateModel{{ . | title }}(cc.Get{{ . | title }}(), offset); err != nil {
			return err
		}
	{{ end }}

	case nil:
		panic(fmt.Sprintf("checkout context message is empty"))

	default:
		panic(fmt.Sprintf("unexpected type %T in oneof", x))
	}
	return nil
}
