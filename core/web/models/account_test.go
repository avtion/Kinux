package models

import (
	"Kinux/tools"
	"context"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestPasswordEncode(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				password: tools.GetRandomString(32),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pw := &Password{
				Raw: tt.args.password,
			}
			got, err := pw.Encode()
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if ok, err := pw.Verify(got); err != nil || !ok {
				t.Fatal(got, err)
			}
			t.Log("密码算法测试成功", got)
		})
	}
}

func TestNewAccounts(t *testing.T) {
	type args struct {
		ctx context.Context
		acs []*AccountWithProfile
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
				acs: []*AccountWithProfile{
					{
						Account: Account{
							Username: "test",
							Password: "test",
							Role:     RoleNormalAccount,
						},
						Profile: Profile{
							RealName: "test",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := NewAccounts(tt.args.ctx, tt.args.acs...); (err != nil) != tt.wantErr {
				t.Errorf("NewAccounts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
