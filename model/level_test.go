package model

import (
	"reflect"
	"testing"
)

func Test_parseString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []keyVal
	}{
		{
			name: "empty",
			args: args{
				s: "",
			},
			want: nil,
		},
		{
			name: "no key",
			args: args{
				s: "hello",
			},
			want: []keyVal{{"", "hello"}},
		},
		{
			name: "only keys",
			args: args{
				s: `Key:value
OtherKey:other value`,
			},
			want: []keyVal{{"Key", "value\n"}, {"OtherKey", "other value"}},
		},
		{
			name: "multi line",
			args: args{
				s: `Key:value
foo
OtherKey:other value`,
			},
			want: []keyVal{{"Key", "value\nfoo\n"}, {"OtherKey", "other value"}},
		},
		{
			name: "mixed",
			args: args{
				s: `Preamble
value
OtherKey:other value`,
			},
			want: []keyVal{{"", "Preamble\nvalue\n"}, {"OtherKey", "other value"}},
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseString(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseString() = %v, want %v", got, tt.want)
			}
		})
	}
}
