package k8s

import (
	"context"
	appV1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/labels"
	"os"

	"testing"
)

const __testDeploymentConfig = "../../example_configs/ubuntu-vnc-deploymnet.yaml"

func Test_parseDeploymentConfig(t *testing.T) {
	fileRaw, err := os.ReadFile(__testDeploymentConfig)
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
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
				fileRaw: fileRaw,
				strict:  false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDp, err := ParseDeploymentConfig(tt.args.fileRaw, tt.args.strict)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDeploymentConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(gotDp)
		})
	}
}

func TestNewDeployment(t *testing.T) {
	dpData, err := os.ReadFile(__testDeploymentConfig)
	if err != nil {
		t.Fatal(err)
	}
	dp, err := ParseDeploymentConfig(dpData, true)
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx  context.Context
		dp   *appV1.Deployment
		s    labels.Set
		opts []DeploymentOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				ctx:  context.Background(),
				dp:   dp,
				s:    labels.Set{"account-id": "0"},
				opts: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := NewDeployment(tt.args.ctx, tt.args.dp, tt.args.s, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("NewDeployment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteDeployments(t *testing.T) {
	type args struct {
		ctx context.Context
		ns  string
		s   labels.Set
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				ctx: context.Background(),
				ns:  "default",
				s:   labels.Set{"account-id": "1"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteDeployments(tt.args.ctx, tt.args.ns, tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("DeleteDeployments() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
