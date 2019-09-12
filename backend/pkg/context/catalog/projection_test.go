package catalog

import (
	"reflect"
	"testing"
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

	type args struct {
		msg *Product
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
				msg: &Product{Id: "123-456", Name: "Name123", Price: 123},
				m:   newModel(),
			},
			want: &model{
				products: map[string]*Product{
					"123-456": &Product{Id: "123-456", Name: "Name123", Price: 123},
				},
				sortedByPrice: []*Product{&Product{Id: "123-456", Name: "Name123", Price: 123}},
				sortedByName:  []*Product{&Product{Id: "123-456", Name: "Name123", Price: 123}},
				sortedByUUID:  []*Product{&Product{Id: "123-456", Name: "Name123", Price: 123}},
			},
		},
		{
			name: "change name product",
			args: args{
				msg: &Product{Id: "123-456", Name: "Name124", Price: 123},
				m: &model{
					products: map[string]*Product{
						"123-456": &Product{Id: "123-456", Name: "Name123", Price: 123},
					},
					sortedByPrice: []*Product{&Product{Id: "123-456", Name: "Name123", Price: 123}},
					sortedByName:  []*Product{&Product{Id: "123-456", Name: "Name123", Price: 123}},
					sortedByUUID:  []*Product{&Product{Id: "123-456", Name: "Name123", Price: 123}},
				},
			},
			want: &model{
				products: map[string]*Product{
					"123-456": &Product{Id: "123-456", Name: "Name124", Price: 123},
				},
				sortedByPrice: []*Product{&Product{Id: "123-456", Name: "Name124", Price: 123}},
				sortedByName:  []*Product{&Product{Id: "123-456", Name: "Name124", Price: 123}},
				sortedByUUID:  []*Product{&Product{Id: "123-456", Name: "Name124", Price: 123}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateModelProduct(tt.args.m, 50, tt.args.msg)
			got := tt.args.m
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("modelRealtimeUpdater() = %v, want %v", got, tt.want)
			}
		})
	}
}
