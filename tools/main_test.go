package tools

import "testing"

func TestGetRandomString(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				n: 6,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(GetRandomString(tt.args.n))
		})
	}
}
