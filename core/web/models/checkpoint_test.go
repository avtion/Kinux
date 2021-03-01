package models

import (
	"context"
	"gorm.io/gorm"
	"testing"
)

func TestCheckpoint_CreateOrUpdate(t *testing.T) {
	type fields struct {
		Model  gorm.Model
		Name   string
		Desc   string
		In     string
		Out    string
		Method uint
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test",
			fields: fields{
				Name:   "whoami",
				Desc:   "whoami命令测试",
				In:     "whoami",
				Out:    "",
				Method: MethodExec,
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Checkpoint{
				Model:  tt.fields.Model,
				Name:   tt.fields.Name,
				Desc:   tt.fields.Desc,
				In:     tt.fields.In,
				Out:    tt.fields.Out,
				Method: tt.fields.Method,
			}
			if err := c.Create(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuickListCheckpoint(t *testing.T) {
	type args struct {
		ctx    context.Context
		name   string
		method CheckpointMethod
	}
	tests := []struct {
		name    string
		args    args
		wantRes []*checkpointQuickResult
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				ctx:    context.Background(),
				name:   "",
				method: 0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := QuickListCheckpoint(tt.args.ctx, tt.args.name, tt.args.method)
			if (err != nil) != tt.wantErr {
				t.Errorf("QuickListCheckpoint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, v := range gotRes {
				t.Log(v)
			}
		})
	}
}
