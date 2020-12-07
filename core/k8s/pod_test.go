package k8s

import (
	"context"
	"testing"
)

func TestGetPods(t *testing.T) {
	type args struct {
		ctx     context.Context
		account string
		job     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				ctx:     context.Background(),
				account: "avtion",
				job:     "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotP, err := GetPods(tt.args.ctx, tt.args.account, tt.args.job)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPods() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(gotP)
		})
	}
}
