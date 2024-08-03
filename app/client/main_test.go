package main

import (
	"reflect"
	"testing"
)

func Test_countEachCall(t *testing.T) {
	type args struct {
		messages []string
	}
	tests := []struct {
		name string
		args args
		want resultCall
	}{
		{
			name: "only one variation with single occurance",
			args: args{
				messages: []string{"aaaa"},
			},
			want: resultCall{
				Detail: map[string]int{
					"aaaa": 1,
				},
			},
		},
		{
			name: "only one variation with multiple occurance",
			args: args{
				messages: []string{"aaaa", "aaaa"},
			},
			want: resultCall{
				Detail: map[string]int{
					"aaaa": 2,
				},
			},
		},
		{
			name: "two variation with single occurance",
			args: args{
				messages: []string{"aaaa", "bbbb"},
			},
			want: resultCall{
				Detail: map[string]int{
					"aaaa": 1,
					"bbbb": 1,
				},
			},
		},
		{
			name: "two variation with multiple occurance",
			args: args{
				messages: []string{"aaaa", "bbbb", "aaaa"},
			},
			want: resultCall{
				Detail: map[string]int{
					"aaaa": 2,
					"bbbb": 1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countEachCall(tt.args.messages); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("countEachCall() = %v, want %v", got, tt.want)
			}
		})
	}
}
