package helper

import "testing"

func TestConvertPriceToFloat64(t *testing.T) {
	type args struct {
		price uint64
	}

	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Large number",
			args: args{1000},
			want: 10,
		},
		{
			name: "Medium number",
			args: args{120},
			want: 1.2,
		},
		{
			name: "Small number",
			args: args{123},
			want: 1.23,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertPriceToFloat64(tt.args.price); got != tt.want {
				t.Errorf("ConvertPriceToFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertPriceToUint64(t *testing.T) {
	type args struct {
		price float64
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			name: "Large number",
			args: args{10},
			want: 1000,
		},
		{
			name: "Medium number",
			args: args{1.2},
			want: 120,
		},
		{
			name: "Small number",
			args: args{1.23},
			want: 123,
		},
		{
			name: "Very small number",
			args: args{1.2323},
			want: 123,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertPriceToUint64(tt.args.price); got != tt.want {
				t.Errorf("ConvertPriceToUint64() = %v, want %v", got, tt.want)
			}
		})
	}
}
