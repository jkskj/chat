package api

import (
	"chat/pkg/e"
	"chat/service"
	"github.com/gin-gonic/gin"
)

// UserRegister 用户注册
func UserRegister(c *gin.Context) {
	var userRegister service.UserService
	//绑定结构体
	code := e.SUCCESS
	err := c.ShouldBind(&userRegister)
	if err == nil {
		res := userRegister.Register()
		c.JSON(code, res)
	} else {
		code = e.InvalidParams
		c.JSON(code, err)
	}
}

// UserLogin 用户登录
func UserLogin(c *gin.Context) {
	var userLogin service.UserService
	code := e.SUCCESS
	err := c.ShouldBind(&userLogin)
	if err == nil {
		res := userLogin.Login()
		c.JSON(code, res)
	} else {
		code = e.InvalidParams
		c.JSON(code, err)
	}
}
