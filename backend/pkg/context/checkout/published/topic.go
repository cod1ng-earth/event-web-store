//go:generate protoc --go_out=. topic.proto

package published

const (
	Topic     = "checkout_pub"
	Partition = 0
)
