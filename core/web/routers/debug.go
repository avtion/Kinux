package routers

import (
	"Kinux/core/web/middlewares"
	"Kinux/core/web/models"
	"Kinux/core/web/msg"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
)

func init() {
	addRouters(debugRouters())
}

// Debug调试专用路由
// TODO 注释
func debugRouters() initFunc {
	return func(r *gin.Engine) {
		debug := r.Group("/debug")
		{
			debug.GET("/login/*role", func(c *gin.Context) {
				// 支持根据不同的role获取对应权限的JWT
				roleStr := c.Param("role")
				var role = models.RoleAdmin
				switch r := cast.ToUint(roleStr); r {
				case models.RoleAnonymous, models.RoleNormalAccount, models.RoleManager, models.RoleAdmin:
					role = r
				default:
				}
				token, _, err := middlewares.TokenCentral.TokenGenerator(&middlewares.TokenPayload{
					Username: "systemTestAccount",
					Role:     role,
				})
				if err != nil {
					c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed("生成测试密钥失败"))
					return
				}
				c.JSON(http.StatusOK, map[string]string{
					"msg":   "系统测试密钥生成成功",
					"token": middlewares.TokenCentral.TokenHeadName + " " + token,
				})
			})
			debug.Use(middlewares.JsonWebTokenAuth).GET("/auth_test", func(c *gin.Context) {
				_u, isExist := c.Get(middlewares.TokenIdentityKey)
				if !isExist {
					c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed("no account info"))
					return
				}
				u, _ := _u.(*middlewares.TokenPayload)
				if u == nil {
					c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed("account info is nil"))
					return
				}
				c.JSON(http.StatusOK, msg.BuildSuccess(*u))
			})
			debug.Use(middlewares.JsonWebTokenAuth).GET("/casbin_test", func(c *gin.Context) {
				c.JSON(http.StatusOK, msg.BuildSuccess("ok"))
			})
		}
	}
}
