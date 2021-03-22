package controllers

import (
	"Kinux/core/web/middlewares"
	"Kinux/core/web/models"
	"Kinux/core/web/msg"
	"Kinux/core/web/services"
	"Kinux/tools"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// 账号
type Account struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 登陆入口
func LoginAccount(c *gin.Context) {
	var _ac = new(Account)
	if err := c.ShouldBindJSON(_ac); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	// 登陆鉴权
	ac, err := services.LoginAccount(c, _ac.Username, _ac.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	// 获取密钥Token
	token, ttl, err := middlewares.TokenCentral.TokenGenerator(&middlewares.TokenPayload{
		Username: ac.Username,
		Role:     ac.Role,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	// 用户真实姓名
	var realName, avatarSeed string
	if profile, _err := ac.GetProfile(c); _err == nil {
		realName = profile.RealName
		avatarSeed = profile.AvatarSeed
	}

	// 用户班级
	dp, _ := ac.GetDepartment(c)

	c.JSON(http.StatusOK, msg.BuildSuccess(map[string]string{
		"msg":        "🛫️登陆成功",
		"token":      token,
		"ttl":        strconv.FormatInt(ttl.Unix(), 10),
		"username":   ac.Username,
		"realName":   realName,
		"role":       models.RoleTranslator(ac.Role),
		"department": dp.Name,
		"avatarSeed": avatarSeed,
		"dpID":       cast.ToString(dp.ID),
	}))
	return
}

// 更新用户的头像种子
func UpdateAccountAvatarSeed(c *gin.Context) {
	ac, err := services.GetAccountFromCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	seed := tools.GetRandomString(6)
	if err = ac.UpdateAvatarSeed(c, seed); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(seed))
}

// 更新用户密码
func UpdatePassword(c *gin.Context) {
	passwordSetter := &struct {
		Old string `json:"old"`
		New string `json:"new"`
	}{}
	if err := c.ShouldBindJSON(passwordSetter); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	ac, err := services.GetAccountFromCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	// 校验旧密码
	if err = ac.Verify(c, passwordSetter.Old); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed("旧密码错误"))
		return
	}
	if err = ac.UpdatePassword(c, passwordSetter.New); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	// 发送鉴权失败响应要求客户端重新登录
	c.JSON(http.StatusOK, msg.Build(msg.CodeJWTAuthFailed, "密码修改成功，请重新登录"))
}

// 获取用户资料
func ListAccounts(c *gin.Context) {
	params := &struct {
		Name       string
		Department uint
		Role       uint
		Page, Size uint
	}{
		Name:       c.DefaultQuery("name", ""),
		Department: cast.ToUint(c.DefaultQuery("department", "")),
		Role:       cast.ToUint(c.DefaultQuery("role", "0")),
		Page:       cast.ToUint(c.DefaultQuery("page", "1")),
		Size:       cast.ToUint(c.DefaultQuery("size", "10")),
	}
	data, err := models.ListAccountsWithProfiles(c,
		models.NewPageBuilder(int(params.Page), int(params.Size)),
		models.AccountNameFilter(params.Name),
		models.AccountDepartmentFilter(int(params.Department)),
		models.AccountRoleFilter(params.Role),
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	// 转译结果
	type resStruct struct {
		ID           uint   `json:"id"`
		Profile      uint   `json:"profile"`
		Role         string `json:"role"`
		Username     string `json:"username"`
		RealName     string `json:"real_name"`
		Department   string `json:"department"`
		CreatedAt    string `json:"created_at"`
		RoleID       uint   `json:"role_id"`
		DepartmentID uint   `json:"department_id"`
	}
	var res = make([]*resStruct, 0, len(data))
	for _, v := range data {
		res = append(res, &resStruct{
			ID:           v.ID,
			Profile:      v.Profile,
			Role:         models.RoleTranslator(v.Role),
			Username:     v.Username,
			RealName:     v.RealName,
			Department:   v.Department,
			CreatedAt:    v.CreatedAt.Format("2006-01-02 15:04:05"),
			RoleID:       v.Role,
			DepartmentID: v.DepartmentId,
		})
	}

	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

// 统计用户资料
func CountAccounts(c *gin.Context) {
	params := &struct {
		Name       string
		Department uint
		Role       uint
	}{
		Name:       c.DefaultQuery("name", ""),
		Department: cast.ToUint(c.DefaultQuery("department", "")),
		Role:       cast.ToUint(c.DefaultQuery("role", "0")),
	}
	res, err := models.CountAccountsWithProfiles(c,
		models.AccountNameFilter(params.Name),
		models.AccountDepartmentFilter(int(params.Department)),
		models.AccountRoleFilter(params.Role),
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}

// 修改用户资料
func EditAccount(c *gin.Context) {
	params := &struct {
		ID         int    `json:"id"`
		Role       int    `json:"role"`
		Department int    `json:"department"`
		Username   string `json:"username"`
		RealName   string `json:"real_name"`
		Password   string `json:"password"`
	}{}
	if err := c.ShouldBindJSON(params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	ac, err := models.GetAccountByID(c, params.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	// 更新用户信息
	err = (&models.Account{
		Model: gorm.Model{
			ID: ac.ID,
		},
		Username: params.Username,
		Role:     models.RoleIdentify(params.Role),
	}).Update(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	// 更新密码
	if params.Password != "" {
		err = ac.UpdatePassword(c, params.Password)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
			return
		}
	}

	// 更新用户资料
	err = (&models.Profile{
		Model: gorm.Model{
			ID: ac.Profile,
		},
		RealName:   params.RealName,
		Department: uint(params.Department),
	}).Update(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
}

// 删除用户
func DeleteAccount(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	if err := models.DeleteAccount(c, id); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("账号删除成功"))
}

// 新增用户
func AddAccount(c *gin.Context) {
	params := &struct {
		Department int    `json:"department" binding:"required"`
		Role       int    `json:"role" binding:"required"`
		Username   string `json:"username" binding:"required"`
		Password   string `json:"password" binding:"required"`
		AvatarSeed string `json:"avatar_seed"`
		RealName   string `json:"real_name"`
	}{}
	if err := c.ShouldBindJSON(params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}

	if params.AvatarSeed == "" {
		params.AvatarSeed = tools.GetRandomString(6)
	}

	err := models.NewAccounts(c, &models.AccountWithProfile{
		Account: models.Account{
			Username: params.Username,
			Password: params.Password,
			Role:     uint(params.Role),
		},
		Profile: models.Profile{
			RealName:   params.RealName,
			Department: uint(params.Department),
			AvatarSeed: params.AvatarSeed,
		},
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, msg.BuildFailed(err))
		return
	}
	c.JSON(http.StatusOK, msg.BuildSuccess("账号创建成功"))
}
