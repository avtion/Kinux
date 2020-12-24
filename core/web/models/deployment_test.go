package models

import (
	"Kinux/tools/bytesconv"
	"context"
	"testing"
)

const (
	__testCentosDpRaw = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: centos-vnc
  namespace: default
#  labels:
#    app: ubuntu-os-vnc
spec:
  replicas: 1
  selector:
    matchLabels:
      app: centos-vnc
  template:
    metadata:
      labels:
        app: centos-vnc
    spec:
      nodeSelector:
        cpu: amd64
      containers:
        - name: centos
          image: consol/centos-xfce-vnc
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: 80
          volumeMounts:
            - mountPath: /dev/shm
              name: shm
      volumes:
        - name: shm
          hostPath:
            # 内存
            path: /dev/shm
`
	__testUbuntuDpRaw = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ubuntu-os-vnc
  namespace: default
#  labels:
#    app: ubuntu-os-vnc
spec:
  replicas: 1
  selector:
    matchLabels:
      app: ubuntu-os-vnc
  template:
    metadata:
      labels:
        app: ubuntu-os-vnc
    spec:
      nodeSelector:
        cpu: amd64
      containers:
        - name: ubuntu
          image: dorowu/ubuntu-desktop-lxde-vnc
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: 80
          volumeMounts:
            - mountPath: /dev/shm
              name: shm
      volumes:
        - name: shm
          hostPath:
            # 内存
            path: /dev/shm
`
)

func TestCrateDeployment(t *testing.T) {
	centosDp := bytesconv.StringToBytes(__testCentosDpRaw)
	ubuntuDp := bytesconv.StringToBytes(__testUbuntuDpRaw)
	type args struct {
		ctx  context.Context
		name string
		raw  []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "centos",
			args: args{
				ctx:  context.Background(),
				name: "centos-vnc",
				raw:  centosDp,
			},
			wantErr: false,
		},
		{
			name: "ubuntu",
			args: args{
				ctx:  context.Background(),
				name: "ubuntu-vnc",
				raw:  ubuntuDp,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotId, err := CrateOrUpdateDeployment(tt.args.ctx, tt.args.name, tt.args.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("CrateOrUpdateDeployment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(gotId)
		})
	}
}

func TestListDeployment(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
		page *PageBuilder
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
				name: "",
				page: NewPageBuilder(0, 10),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := ListDeployment(tt.args.ctx, tt.args.name, tt.args.page)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListDeployment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(gotRes)
		})
	}
}

func TestGetDeployment(t *testing.T) {
	type args struct {
		ctx context.Context
		id  uint
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
				id:  2,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := GetDeployment(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDeployment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(gotRes)
		})
	}
}
