//go:generate protoc --go_out=. topic.proto

package published

const (
	Topic     = "pim_pub"
	Partition = 0
)
