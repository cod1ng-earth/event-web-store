//go:generate protoc --go_out=. topic.proto

package checkout

const (
	Topic     = "checkout_pub"
	Partition = 0
)
