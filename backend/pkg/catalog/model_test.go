package catalog

import (
	"reflect"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

func Test_remove(t *testing.T) {
	type args struct {
		slice []*Product
		i     int
	}
	tests := []struct {
		name string
		args args
		want []*Product
	}{
		{
			name: "one element in list",
			args: args{
				slice: []*Product{&Product{Id: "1"}},
				i:     0,
			},
			want: []*Product{},
		},
		{
			name: "two element in list",
			args: args{
				slice: []*Product{&Product{Id: "1"}, &Product{Id: "2"}},
				i:     0,
			},
			want: []*Product{&Product{Id: "2"}},
		},
		{
			name: "remove 2nd element in list",
			args: args{
				slice: []*Product{&Product{Id: "1"}, &Product{Id: "2"}},
				i:     1,
			},
			want: []*Product{&Product{Id: "1"}},
		},
		{
			name: "remove middle element",
			args: args{
				slice: []*Product{&Product{Id: "1"}, &Product{Id: "2"}, &Product{Id: "3"}},
				i:     1,
			},
			want: []*Product{&Product{Id: "1"}, &Product{Id: "3"}},
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
		slice []*Product
		i     int
		p     *Product
	}
	tests := []struct {
		name string
		args args
		want []*Product
	}{
		{
			name: "add element to empty list",
			args: args{
				slice: []*Product{},
				i:     0,
				p:     &Product{Id: "1"},
			},
			want: []*Product{&Product{Id: "1"}},
		},
		{
			name: "add element after one element list",
			args: args{
				slice: []*Product{&Product{Id: "1"}},
				i:     1,
				p:     &Product{Id: "2"},
			},
			want: []*Product{&Product{Id: "1"}, &Product{Id: "2"}},
		},
		{
			name: "add element before one element list",
			args: args{
				slice: []*Product{&Product{Id: "1"}},
				i:     0,
				p:     &Product{Id: "2"},
			},
			want: []*Product{&Product{Id: "2"}, &Product{Id: "1"}},
		},
		{
			name: "add element between two elements",
			args: args{
				slice: []*Product{&Product{Id: "1"}, &Product{Id: "2"}},
				i:     1,
				p:     &Product{Id: "3"},
			},
			want: []*Product{&Product{Id: "1"}, &Product{Id: "3"}, &Product{Id: "2"}},
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

	encodeProduct := func(pu *Product) *sarama.ConsumerMessage {
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
				msg: encodeProduct(&Product{Id: "123-456", Name: "Name123", Price: 1.23}),
				m:   newModel(),
			},
			want: &model{
				products: map[string]*Product{
					"123-456": &Product{Id: "123-456", Name: "Name123", Price: 1.23},
				},
				sortedByPrice: []*Product{&Product{Id: "123-456", Name: "Name123", Price: 1.23}},
				sortedByName:  []*Product{&Product{Id: "123-456", Name: "Name123", Price: 1.23}},
				sortedByUUID:  []*Product{&Product{Id: "123-456", Name: "Name123", Price: 1.23}},
			},
		},
		{
			name: "change name product",
			args: args{
				msg: encodeProduct(&Product{Id: "123-456", Name: "Name124", Price: 1.23}),
				m: &model{
					products: map[string]*Product{
						"123-456": &Product{Id: "123-456", Name: "Name123", Price: 1.23},
					},
					sortedByPrice: []*Product{&Product{Id: "123-456", Name: "Name123", Price: 1.23}},
					sortedByName:  []*Product{&Product{Id: "123-456", Name: "Name123", Price: 1.23}},
					sortedByUUID:  []*Product{&Product{Id: "123-456", Name: "Name123", Price: 1.23}},
				},
			},
			want: &model{
				products: map[string]*Product{
					"123-456": &Product{Id: "123-456", Name: "Name124", Price: 1.23},
				},
				sortedByPrice: []*Product{&Product{Id: "123-456", Name: "Name124", Price: 1.23}},
				sortedByName:  []*Product{&Product{Id: "123-456", Name: "Name124", Price: 1.23}},
				sortedByUUID:  []*Product{&Product{Id: "123-456", Name: "Name124", Price: 1.23}},
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
