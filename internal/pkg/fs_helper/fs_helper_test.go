package fshelper

import (
	"reflect"
	"testing"
)

func TestGetPeriodHeaderOrderIndex(t *testing.T) {
	type args struct {
		headers []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "test1",
			args: args{
				headers: []string{"期末金额", "年初金额"},
			},
			want: []string{"年初金额", "期末金额"},
		},
		{
			name: "test2",
			args: args{
				headers: []string{"本期金额", "上期金额"},
			},
			want: []string{"上期金额", "本期金额"},
		},
		{
			name: "test3",
			args: args{
				headers: []string{"期末数", "期初数"},
			},
			want: []string{"期初数", "期末数"},
		},
		{
			name: "test4",
			args: args{
				headers: []string{"期末金额", "年初金额"},
			},
			want: []string{"年初金额", "期末金额"},
		},
		{
			name: "test5",
			args: args{
				headers: []string{"2019年", "2018年（调整后）", "2018年（调整前）", "2017年"},
			},
			want: []string{"2017年", "2018年（调整前）", "2018年（调整后）", "2019年"},
		},
		{
			name: "test6",
			args: args{
				headers: []string{"本年累计数", "上年同期数"},
			},
			want: []string{"上年同期数", "本年累计数"},
		},
		{
			name: "test7",
			args: args{
				headers: []string{"前年同期数", "本年累计数", "去年同期数"},
			},
			want: []string{"前年同期数", "去年同期数", "本年累计数"},
		},
		{
			name: "test8",
			args: args{
				headers: []string{"", "下半年", "上半年"},
			},
			want: []string{"", "上半年", "下半年"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPeriodHeaderOrderIndex(tt.args.headers); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPeriodHeaderOrderIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}
