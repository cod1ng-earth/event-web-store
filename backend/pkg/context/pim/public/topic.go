//go:generate protoc --go_out=. topic.proto

package pim

const (
	Topic     = "pim_pub"
	Partition = 0
)
