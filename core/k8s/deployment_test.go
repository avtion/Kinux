package k8s

import (
	"context"
	"io/ioutil"

	"testing"
)

const __testDeploymentConfig = "../../example_configs/ubuntu-vnc-deploymnet.yaml"

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

func Test_parseDeploymentConfig(t *testing.T) {
	fileRaw, err := ioutil.ReadFile(__testDeploymentConfig)
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx     context.Context
		fileRaw []byte
		strict  bool
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
				fileRaw: fileRaw,
				strict:  false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDp, err := ParseDeploymentConfig(tt.args.ctx, tt.args.fileRaw, tt.args.strict)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDeploymentConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(gotDp)
		})
	}
}
