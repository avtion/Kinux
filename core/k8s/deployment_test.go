package k8s

import (
	"context"
	"testing"
)

func TestNewDeployment(t *testing.T) {
	type args struct {
		ctx       context.Context
		AccountId string
		JobId     string
		options   []DeploymentOption
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				ctx:       context.Background(),
				AccountId: "avtion",
				JobId:     "test",
				options:   nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewDeployment(tt.args.ctx, tt.args.AccountId, tt.args.JobId, tt.args.options...)
		})
	}
}
