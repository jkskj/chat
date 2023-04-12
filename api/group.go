package api

import (
	"chat/pkg/e"
	"chat/pkg/utils"
	"chat/service"
	"github.com/gin-gonic/gin"
)

// GetGroups 获取群聊列表
func GetGroups(c *gin.Context) {
	var group service.GroupService
	code := e.SUCCESS
	chaim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	res := group.Get(chaim.Id)
	c.JSON(code, res)
}

// MakeGroups 创建群聊
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

// JoinGroups 加入群聊
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

// GroupMessage 获取群聊消息
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
