package api

import (
	"chat/pkg/e"
	"chat/pkg/utils"
	"chat/service"
	"github.com/gin-gonic/gin"
)

func GetFriends(c *gin.Context) {
	var friend service.FriendService
	code := e.SUCCESS
	chaim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	res := friend.Get(chaim.Id)
	c.JSON(code, res)

}

func MakeFriends(c *gin.Context) {
	var friend service.FriendService
	code := e.SUCCESS
	chaim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	err := c.ShouldBind(&friend)
	if err == nil {
		res := friend.Make(chaim.Id)
		c.JSON(code, res)
	} else {
		code = e.InvalidParams
		c.JSON(code, err)
	}
}
func GetApplication(c *gin.Context) {
	var friend service.FriendService
	code := e.SUCCESS
	chaim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	err := c.ShouldBind(&friend)
	if err == nil {
		res := friend.GetApplication(chaim.Id)
		c.JSON(code, res)
	} else {
		code = e.InvalidParams
		c.JSON(code, err)
	}
}
func Reply(c *gin.Context) {
	var friend service.FriendService
	code := e.SUCCESS
	chaim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	err := c.ShouldBind(&friend)
	if err == nil {
		res := friend.Reply(chaim.Id)
		c.JSON(code, res)
	} else {
		code = e.InvalidParams
		c.JSON(code, err)
	}
}
func SingleMessage(c *gin.Context) {
	var group service.FriendService
	code := e.SUCCESS
	err := c.ShouldBind(&group)
	chaim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err == nil {
		res := group.Message(chaim.Id)
		c.JSON(code, res)
	} else {
		code = e.InvalidParams
		c.JSON(code, err)
	}
}
