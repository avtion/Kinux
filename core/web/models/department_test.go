package models

import (
	"context"
	"testing"
)

func TestNewDepartment(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
		opts []DepartmentOpt
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
				name: "test2",
				opts: []DepartmentOpt{
					DepartmentNsOpt("default", "test"),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, err := NewDepartment(tt.args.ctx, tt.args.name, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDepartment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(gotD)
		})
	}
}
