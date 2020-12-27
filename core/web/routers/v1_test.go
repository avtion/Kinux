package routers

import (
	"Kinux/core/web/middlewares"
	"Kinux/core/web/models"
	"Kinux/tools/bytesconv"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_QueryMissions(t *testing.T) {
	const route = "/v1/mission/"
	// 生成JWT密钥
	token, _, err := middlewares.TokenCentral.TokenGenerator(&middlewares.TokenPayload{
		Username: "test",
		Role:     models.RoleNormalAccount,
	})
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, route, nil)
	req.Header.Set("Authorization", middlewares.TokenCentral.TokenHeadName+" "+token)
	NewRouters().ServeHTTP(w, req)
	respData, _ := ioutil.ReadAll(w.Result().Body)
	t.Log(bytesconv.BytesToString(respData))
	_ = w.Result().Body.Close()
}
