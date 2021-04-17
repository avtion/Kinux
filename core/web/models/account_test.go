package models

import (
	"Kinux/tools"
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
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
							Username: "admin",
							Password: "admin",
							Role:     RoleManager,
						},
						Profile: Profile{
							RealName: "系统管理员",
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

func TestListAccountsWithProfiles(t *testing.T) {
	type args struct {
		ctx     context.Context
		builder *PageBuilder
		filters []AccountFilterFn
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
				builder: NewPageBuilder(1, 10),
				filters: []AccountFilterFn{
					AccountDepartmentFilter(1),
					AccountNameFilter("t"),
					AccountRoleFilter(RoleNormalAccount),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := ListAccountsWithProfiles(tt.args.ctx, tt.args.builder, tt.args.filters...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListAccountsWithProfiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, v := range gotRes {
				t.Log(v)
			}
		})
	}
}

func TestNewExampleAccounts(t *testing.T) {
	// 初始化姓名生成
	rand.Seed(time.Now().Unix())
	var names = make([]string, 0)
	raw, err := os.ReadFile("../../../example_configs/names.txt")
	if err != nil {
		t.Fatal(err)
	}
	names = strings.Split(string(raw), "\n")
	t.Log(names[rand.Intn(len(names))])

	dps, err := ListDepartments(context.Background(), "", nil)
	if err != nil {
		t.Fatal(err)
	}
	var distinctMapping = make(map[string]struct{})
	var accounts = make([]*AccountWithProfile, 0, len(dps)*60)

	for index, dp := range dps {
		for k := range [60]struct{}{} {
			// 生成用户名
			var username = dp.Name[0:2] + fmt.Sprintf("551117%d", index) + fmt.Sprintf("%02d", k)

			// 生成姓名
			var realName = names[rand.Intn(len(names))]
			for {
				if _, isExist := distinctMapping[realName]; !isExist {
					distinctMapping[realName] = struct{}{}
					break
				}
				realName = names[rand.Intn(len(names))]
			}

			// 生成用户并添加
			accounts = append(accounts, &AccountWithProfile{
				Account: Account{
					Username: username,
					Password: username,
					Role:     RoleNormalAccount,
				},
				Profile: Profile{
					Department: dp.ID,
					RealName:   realName,
					AvatarSeed: tools.GetRandomString(6),
				},
			})
		}
	}
	if err = NewAccounts(context.Background(), accounts...); err != nil {
		t.Fatal(err)
	}
	t.Log("示例用户生成成功")
}
