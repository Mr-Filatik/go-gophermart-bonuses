package helper

import "testing"

func TestValidateAlgorithmLuhn(t *testing.T) {
	type args struct {
		number uint64
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Correct number",
			args: args{1234567897},
			want: true,
		},
		{
			name: "Uncorrect number",
			args: args{1234567893},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateAlgorithmLuhn(tt.args.number); got != tt.want {
				t.Errorf("ValidateAlgorithmLuhn() = %v, want %v", got, tt.want)
			}
		})
	}
}
