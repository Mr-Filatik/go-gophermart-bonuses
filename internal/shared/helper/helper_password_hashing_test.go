package helper

import "testing"

func TestPasswordHashing(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "Standart password",
			args:    args{"password"},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := GeneratePasswordHash(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratePasswordHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got := ComparePasswordHashes(hash, tt.args.password); got != tt.want {
				t.Errorf("ComparePasswordHashes() = %v, want %v", got, tt.want)
			}
		})
	}
}
