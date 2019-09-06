// Code generated by go-bindata.
// sources:
// context.go.tpl
// DO NOT EDIT!

package simba

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)
type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _contextGoTpl = []byte(`// Code generated by simba. DO NOT EDIT.

package {{ .Name }}

import (
	"fmt"
	"log"
	"sync"
{{ if eq .ReadLock "wait-free" }}
	"sync/atomic"
	"time"
{{ end }}

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"

{{ range .Bridges }}
	{{ .Name }} "{{ .PkgPath }}/public"
{{ end }}
)

const (
	Topic     = "{{ .Name }}_internal"
	Partition = 0
)

type context struct {
	doneCh    chan struct{}
{{ range .Bridges }}
	doneCh{{ .Name | title }}    chan struct{}
{{ end }}
	client    sarama.Client
	consumer  sarama.Consumer
	producer  sarama.SyncProducer

	batchOffset int64

{{ if eq .ReadLock "exclusive" }}
	model  *model
	lock   *sync.Mutex
{{ else if eq .ReadLock "parallel" }}
	model  *model
	lock   *sync.RWMutex
{{ else if eq .ReadLock "wait-free" }}
	readerAChanged *sync.Cond
	readerBChanged *sync.Cond
	aIsReading     bool
	readersA       int32
	readersB       int32
	modelA         *model
	modelB         *model

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
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		log.Panicf("failed to setup kafka producer: %s", err)
	}

	batchOffset, err := client.GetOffset(Topic, Partition, sarama.OffsetNewest)
	if err != nil {
		log.Printf("failed to get last offset for topic %v partition %v", Topic, Partition)
	}
	batchOffset--

	context := context{
		doneCh:    make(chan struct{}, 1),
	{{ range .Bridges }}
		doneCh{{ .Name | title }}: make(chan struct{}, 1),
	{{ end }}
		client:    client,
		consumer:  consumer,
		producer:  producer,

		batchOffset: batchOffset,

	{{ if eq .ReadLock "exclusive" }}
		model:  newModel(),
		lock:   &sync.Mutex{},
	{{ else if eq .ReadLock "parallel" }}
		model:  newModel(),
		lock:   &sync.RWMutex{},
	{{ else if eq .ReadLock "wait-free" }}
		readerAChanged: sync.NewCond(&sync.Mutex{}),
		readerBChanged: sync.NewCond(&sync.Mutex{}),
		aIsReading:     true,
		readersA:		0,
		readersB:		0,
		modelA:         newModel(),
		modelB:         newModel(),

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

func (c *context) await(offset int64) {
	if offset == -1 {
		return
	}
	if c.offset >= offset {
		return
	}
	c.offsetChanged.L.Lock()
	for c.offset < offset {
		c.offsetChanged.Wait()
	}
	c.offsetChanged.L.Unlock()
}

func (c *context) AwaitLastOffset() {
	c.await(c.batchOffset)
}

func (c *context) updateLoop(writes <-chan *sarama.ConsumerMessage) {

	{{ if eq .ReadLock "wait-free" }}
	writesRedo := make(chan *sarama.ConsumerMessage, 32768)
	{{ end }}

	for {
	{{ if or (eq .ReadLock "exclusive") (eq .ReadLock "parallel") }}
		for msg := range writes {
			applyChange(msg, c.model, c)
		}
	{{ else if eq .ReadLock "wait-free" }}
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

		for len(writesRedo) > 0 {
			msg, ok := <-writesRedo
			if ok {
				applyChange(msg, model, c)
			}
		}

		msg, ok := <-writes
		if !ok {
			return
		}
		writesRedo <- msg
		start := time.Now()
		applyChange(msg, model, c)
		offset := msg.Offset
		writeDelay += time.Since(start)

		for writeDelay < 10*time.Millisecond && len(writes) > 0 && len(writesRedo) < 32768 {
			msg, ok := <-writes
			if !ok {
				return
			}
			writesRedo <- msg
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

//	log.Printf("applying message with offset %v", msg.Offset)

{{ if or (eq .ReadLock "exclusive") (eq .ReadLock "parallel") }}
	c.lock.Lock()
	defer c.lock.Unlock()
	defer func() {
		c.offset = msg.Offset
		c.offsetChanged.Broadcast()
	}()
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

{{ range .Bridges }}
func (c *context) bridge{{ .Name | title }}() {

	c.AwaitLastOffset()

	model, free := c.read()
	{{ .Name }}Offset := model.get{{ .Name | title }}Offset()
	free()
	partition, err := c.consumer.ConsumePartition({{ .Name }}.Topic, 0, {{ .Name }}Offset)
	if err != nil {
		log.Panicf("failed to setup kafka partition: %s", err)
	}

	log.Printf("starting bridge %v", Topic)

	for {
		select {
		case err := <-partition.Errors():
			log.Printf("failure from kafka consumer: %s", err)

		case msg := <-partition.Messages():
//			log.Printf("recieved message with offset %v", msg.Offset)

			cc := {{ .Name }}.TopicMessage{}
			err := proto.Unmarshal(msg.Value, &cc)
			if err != nil {
				log.Fatalf("failed to unmarshal kafka massage %s/%d: %v", Topic, msg.Offset, err)
			}

			switch x := cc.GetMessages().(type) {

			{{$bridge := .}}
			{{ range .MessageNames }}
			case *{{ $bridge.Name }}.TopicMessage_{{ . | title }}:
				if err := translate{{ $bridge.Name | title }}{{ . | title }}(c, model, msg.Offset, cc.Get{{ . | title }}()); err != nil {
					log.Fatalf("failed to translate kafka message $bridge.Name/%v: %s", msg.Offset, err)
				}
			{{ end }}

			case nil:
				panic(fmt.Sprintf("context message is empty"))

			default:
				panic(fmt.Sprintf("unexpected type %T in oneof", x))
			}

		}
	}
}
{{ end }}

func (c *context) Start() {

	log.Printf("starting context %v", Topic)

	writes := make(chan *sarama.ConsumerMessage, 32768)
	go c.updateLoop(writes)

	partition, err := c.consumer.ConsumePartition(Topic, 0, 0)
	if err != nil {
		log.Panicf("failed to setup kafka partition: %s", err)
	}

{{ range .Bridges }}
		go c.bridge{{ .Name | title }}()
{{ end }}

	for {
		select {
		case err := <-partition.Errors():
			log.Printf("failure from kafka consumer: %s", err)

		case msg := <-partition.Messages():
//			log.Printf("recieved message with offset %v", msg.Offset)
			writes <- msg

		case <-c.doneCh:
			log.Print("interrupt is detected")
			{{ range .Bridges }}
				c.doneCh{{ .Name | title }} <- struct{}{}
			{{ end }}
			if err := partition.Close(); err != nil {
				log.Panicf("failed to close kafka partition: %s", err)
			}
			close(writes)
			if err := c.consumer.Close(); err != nil {
				log.Panicf("failed to close kafka consumer: %s", err)
			}
			if err := c.producer.Close(); err != nil {
				log.Panicf("failed to close kafka producer: %s", err)
			}
			if err := c.client.Close(); err != nil {
				log.Panicf("failed to close kafka client: %s", err)
			}
			return
		}
	}
}

func (c *context) read() (*model, func()) {

	c.offsetChanged.L.Lock()
	for c.offset < c.batchOffset {
		c.offsetChanged.Wait()
	}
	c.offsetChanged.L.Unlock()

{{ if eq .ReadLock "exclusive" }}
	c.lock.Lock()
	return c.model, c.lock.Unlock
{{ else if eq .ReadLock "parallel" }}
	c.lock.RLock()
	return c.model, c.lock.RUnlock
{{ else if eq .ReadLock "wait-free" }}
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
	cc := TopicMessage{}
	err := proto.Unmarshal(msg.Value, &cc)
	if err != nil {
		return fmt.Errorf("failed to unmarshal kafka massage %s/%d: %v", Topic, msg.Offset, err)
	}

	switch x := cc.GetMessages().(type) {

	{{ range .MessageNames }}
	case *TopicMessage_{{ . | title }}:
		return batchUpdateModel{{ . | title }}(model, msg.Offset, cc.Get{{ . | title }}())
	{{ end }}

	case nil:
		panic(fmt.Sprintf("context message is empty"))

	default:
		panic(fmt.Sprintf("unexpected type %T in oneof", x))
	}
}
{{ end }}

func updateModel(msg *sarama.ConsumerMessage, model *model) error {
	cc := TopicMessage{}
	err := proto.Unmarshal(msg.Value, &cc)
	if err != nil {
		return fmt.Errorf("failed to unmarshal kafka massage %s/%d: %v", Topic, msg.Offset, err)
	}

	switch x := cc.GetMessages().(type) {

	{{ range .MessageNames }}
	case *TopicMessage_{{ . | title }}:
		return updateModel{{ . | title }}(model, msg.Offset, cc.Get{{ . | title }}())
	{{ end }}

	case nil:
		panic(fmt.Sprintf("context message is empty"))

	default:
		panic(fmt.Sprintf("unexpected type %T in oneof", x))
	}
}

type asyncProducer struct {
	producer sarama.AsyncProducer
	wg       sync.WaitGroup
}

func (c *context) newSyncProducer(f func(error)) (asyncProducer, error) {
	producer, err := sarama.NewAsyncProducerFromClient(c.client)
	if err != nil {
		return asyncProducer{}, fmt.Errorf("failed to create async producer: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for err := range producer.Errors() {
			f(err)
		}
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		for _ = range producer.Successes() {
		}
		wg.Done()
	}()

	return asyncProducer{
		producer: producer,
		wg:       wg,
	}, nil
}

func (p *asyncProducer) Close() {
	p.producer.AsyncClose()
	p.wg.Wait()
}

{{ range .MessageNames }}
func (c *context) log{{ . | title }}(msg *{{ . | title }}) (int32, int64, error) {

	topicMsg := &TopicMessage{
		Messages: &TopicMessage_{{ . | title }}{
			{{ . | title }}: msg,
		},
	}

	bytes, err := proto.Marshal(topicMsg)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to serialize {{ . }} change massage: %v", err)
	}

	producerMsg := &sarama.ProducerMessage{
		Topic: Topic,
		Value: sarama.ByteEncoder(bytes),
	}
	return c.producer.SendMessage(producerMsg)
}

func (p asyncProducer) log{{ . | title }}(msg *{{ . | title }}) error {

	topicMsg := &TopicMessage{
		Messages: &TopicMessage_{{ . | title }}{
			{{ . | title }}: msg,
		},
	}

	bytes, err := proto.Marshal(topicMsg)
	if err != nil {
		return fmt.Errorf("failed to serialize {{ . }} change massage: %v", err)
	}

	producerMsg := &sarama.ProducerMessage{
		Topic: Topic,
		Value: sarama.ByteEncoder(bytes),
	}
	p.producer.Input() <- producerMsg

	return nil
}
{{ end }}
`)

func contextGoTplBytes() ([]byte, error) {
	return _contextGoTpl, nil
}

func contextGoTpl() (*asset, error) {
	bytes, err := contextGoTplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "context.go.tpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"context.go.tpl": contextGoTpl,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"context.go.tpl": &bintree{contextGoTpl, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

