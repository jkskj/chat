package api

import (
	"chat/pkg/e"
	"chat/pkg/utils"
	"chat/service"
	"github.com/gin-gonic/gin"
)

func GetGroups(c *gin.Context) {
	var group service.GroupService
	code := e.SUCCESS
	chaim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	res := group.Get(chaim.Id)
	c.JSON(code, res)
}

func MakeGroups(c *gin.Context) {
	var group service.GroupService
	code := e.SUCCESS
	chaim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	err := c.ShouldBind(&group)
	if err == nil {
		res := group.Make(chaim.Id)
		c.JSON(code, res)
	} else {
		code = e.InvalidParams
		c.JSON(code, err)
	}
}

func JoinGroups(c *gin.Context) {
	var group service.GroupService
	code := e.SUCCESS
	chaim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	err := c.ShouldBind(&group)
	if err == nil {
		res := group.Join(chaim.Id)
		c.JSON(code, res)
	} else {
		code = e.InvalidParams
		c.JSON(code, err)
	}
}

func GroupMessage(c *gin.Context) {
	var group service.GroupService
	code := e.SUCCESS
	err := c.ShouldBind(&group)
	if err == nil {
		res := group.Message()
		c.JSON(code, res)
	} else {
		code = e.InvalidParams
		c.JSON(code, err)
	}
}
