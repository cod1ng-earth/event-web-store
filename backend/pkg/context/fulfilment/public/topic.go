//go:generate protoc --go_out=. topic.proto

package fulfilment

const (
	Topic     = "fulfilment_pub"
	Partition = 0
)
