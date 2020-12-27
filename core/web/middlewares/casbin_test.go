package middlewares

import (
	"Kinux/core/web/models"
	"github.com/spf13/cast"
	"net/http"
	"testing"
)

func Test_initCasbinRule(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enforcer, err := newEnforcer(&gormAdapter{DB: models.GetGlobalDB()})
			if err != nil {
				t.Fatal(err)
			}
			if err = enforcer.SavePolicy(); err != nil {
				t.Fatal(err)
			}
			if err := initCasbinRoles(enforcer); err != nil {
				t.Fatal(err)
			}
			ok, err := enforcer.Enforce(cast.ToString(models.RoleAnonymous), "/", "GET")
			if err != nil {
				t.Fatal(err)
			}
			if !ok {
				t.Log("no permission")
			}
			ok, err = enforcer.Enforce(cast.ToString(models.RoleAnonymous), "/123", "POST")
			if err != nil {
				t.Fatal(err)
			}
			if !ok {
				t.Log("no permission")
			}
		})
	}
}

// 新增一个全局允许的路由规则用于测试
func Test_TurnOffCasbin(t *testing.T) {
	enforcer, err := newEnforcer(&gormAdapter{DB: models.GetGlobalDB()})
	if err != nil {
		t.Fatal(err)
	}
	if err = enforcer.SavePolicy(); err != nil {
		t.Fatal(err)
	}
	if _, err = enforcer.AddPolicy(cast.ToString(models.RoleNormalAccount), "*",
		http.MethodGet,
		http.MethodPost,
		http.MethodOptions,
	); err != nil {
		t.Fatal(err)
	}
}
