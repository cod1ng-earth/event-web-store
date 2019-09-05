//go:generate protoc --go_out=. topic.proto

package published

const (
	Topic     = "fulfilment_pub"
	Partition = 0
)
