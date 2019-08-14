package products

import (
	"reflect"
	"testing"

	"git.votum-media.net/event-web-store/event-web-store/backend/pkg/pb"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

func Test_remove(t *testing.T) {
	type args struct {
		slice []*pb.Product
		i     int
	}
	tests := []struct {
		name string
		args args
		want []*pb.Product
	}{
		{
			name: "one element in list",
			args: args{
				slice: []*pb.Product{&pb.Product{Uuid: "1"}},
				i:     0,
			},
			want: []*pb.Product{},
		},
		{
			name: "two element in list",
			args: args{
				slice: []*pb.Product{&pb.Product{Uuid: "1"}, &pb.Product{Uuid: "2"}},
				i:     0,
			},
			want: []*pb.Product{&pb.Product{Uuid: "2"}},
		},
		{
			name: "remove 2nd element in list",
			args: args{
				slice: []*pb.Product{&pb.Product{Uuid: "1"}, &pb.Product{Uuid: "2"}},
				i:     1,
			},
			want: []*pb.Product{&pb.Product{Uuid: "1"}},
		},
		{
			name: "remove middle element",
			args: args{
				slice: []*pb.Product{&pb.Product{Uuid: "1"}, &pb.Product{Uuid: "2"}, &pb.Product{Uuid: "3"}},
				i:     1,
			},
			want: []*pb.Product{&pb.Product{Uuid: "1"}, &pb.Product{Uuid: "3"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := remove(tt.args.slice, tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("remove() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_insert(t *testing.T) {
	type args struct {
		slice []*pb.Product
		i     int
		p     *pb.Product
	}
	tests := []struct {
		name string
		args args
		want []*pb.Product
	}{
		{
			name: "add element to empty list",
			args: args{
				slice: []*pb.Product{},
				i:     0,
				p:     &pb.Product{Uuid: "1"},
			},
			want: []*pb.Product{&pb.Product{Uuid: "1"}},
		},
		{
			name: "add element after one element list",
			args: args{
				slice: []*pb.Product{&pb.Product{Uuid: "1"}},
				i:     1,
				p:     &pb.Product{Uuid: "2"},
			},
			want: []*pb.Product{&pb.Product{Uuid: "1"}, &pb.Product{Uuid: "2"}},
		},
		{
			name: "add element before one element list",
			args: args{
				slice: []*pb.Product{&pb.Product{Uuid: "1"}},
				i:     0,
				p:     &pb.Product{Uuid: "2"},
			},
			want: []*pb.Product{&pb.Product{Uuid: "2"}, &pb.Product{Uuid: "1"}},
		},
		{
			name: "add element between two elements",
			args: args{
				slice: []*pb.Product{&pb.Product{Uuid: "1"}, &pb.Product{Uuid: "2"}},
				i:     1,
				p:     &pb.Product{Uuid: "3"},
			},
			want: []*pb.Product{&pb.Product{Uuid: "1"}, &pb.Product{Uuid: "3"}, &pb.Product{Uuid: "2"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := insert(tt.args.slice, tt.args.i, tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("insert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_modelRealtimeUpdater(t *testing.T) {

	encodeProductUpdate := func(pu *pb.ProductUpdate) *sarama.ConsumerMessage {
		bytes, err := proto.Marshal(pu)
		if err != nil {
			t.Errorf("failed to marshal product update: %v", err)
		}
		return &sarama.ConsumerMessage{
			Offset: 0,
			Value:  bytes,
		}
	}

	type args struct {
		msg *sarama.ConsumerMessage
		m   *model
	}
	tests := []struct {
		name string
		args args
		want *model
	}{
		{
			name: "empty model",
			args: args{
				msg: encodeProductUpdate(&pb.ProductUpdate{
					New: &pb.Product{Uuid: "123-456", Title: "Name123", Price: 1.23},
				}),
				m: newModel(),
			},
			want: &model{
				products: map[string]*pb.Product{
					"123-456": &pb.Product{Uuid: "123-456", Title: "Name123", Price: 1.23},
				},
				sortedByPrice: []*pb.Product{&pb.Product{Uuid: "123-456", Title: "Name123", Price: 1.23}},
				sortedByTitle: []*pb.Product{&pb.Product{Uuid: "123-456", Title: "Name123", Price: 1.23}},
				sortedByUUID:  []*pb.Product{&pb.Product{Uuid: "123-456", Title: "Name123", Price: 1.23}},
			},
		},
		{
			name: "change name product",
			args: args{
				msg: encodeProductUpdate(&pb.ProductUpdate{
					Old: &pb.Product{Uuid: "123-456", Title: "Name123", Price: 1.23},
					New: &pb.Product{Uuid: "123-456", Title: "Name124", Price: 1.23},
				}),
				m: &model{
					products: map[string]*pb.Product{
						"123-456": &pb.Product{Uuid: "123-456", Title: "Name123", Price: 1.23},
					},
					sortedByPrice: []*pb.Product{&pb.Product{Uuid: "123-456", Title: "Name123", Price: 1.23}},
					sortedByTitle: []*pb.Product{&pb.Product{Uuid: "123-456", Title: "Name123", Price: 1.23}},
					sortedByUUID:  []*pb.Product{&pb.Product{Uuid: "123-456", Title: "Name123", Price: 1.23}},
				},
			},
			want: &model{
				products: map[string]*pb.Product{
					"123-456": &pb.Product{Uuid: "123-456", Title: "Name124", Price: 1.23},
				},
				sortedByPrice: []*pb.Product{&pb.Product{Uuid: "123-456", Title: "Name124", Price: 1.23}},
				sortedByTitle: []*pb.Product{&pb.Product{Uuid: "123-456", Title: "Name124", Price: 1.23}},
				sortedByUUID:  []*pb.Product{&pb.Product{Uuid: "123-456", Title: "Name124", Price: 1.23}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modelRealtimeUpdater(tt.args.msg, tt.args.m)
			got := tt.args.m
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("modelRealtimeUpdater() = %s, want %s", got, tt.want)
			}
		})
	}
}
