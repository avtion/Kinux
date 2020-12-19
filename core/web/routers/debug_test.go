package routers

import (
	"Kinux/core/web/middlewares"
	"Kinux/core/web/models"
	"Kinux/tools/bytesconv"
	"github.com/spf13/cast"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCasbinRouter(t *testing.T) {
	const route = "/debug/casbin_test"
	var testRouter = NewRouters()

	// 写入Casbin规则
	e := middlewares.GetGlobalEnforcer()
	if e == nil {
		t.Fatal(middlewares.ErrEnforcerNil)
	}
	_, err := e.AddPolicy(cast.ToString(models.RoleNormalAccount), route, http.MethodGet)
	if err != nil {
		t.Fatal(err)
	}

	// 生成JWT密钥, 分别是管理员和匿名游客
	token1, _, err := middlewares.TokenCentral.TokenGenerator(&middlewares.TokenPayload{
		Username: "systemTestAccount",
		Role:     models.RoleAdmin,
	})
	if err != nil {
		t.Fatal(err)
	}
	token2, _, err := middlewares.TokenCentral.TokenGenerator(&middlewares.TokenPayload{
		Username: "systemTestAccount",
		Role:     models.RoleAnonymous,
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, token := range []string{token1, token2} {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, route, nil)
		req.Header.Set("Authorization", middlewares.TokenCentral.TokenHeadName+" "+token)
		testRouter.ServeHTTP(w, req)
		respData, _ := ioutil.ReadAll(w.Result().Body)
		t.Log(bytesconv.BytesToString(respData))
		_ = w.Result().Body.Close()
	}
}
