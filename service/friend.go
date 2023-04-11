package service

import (
	"chat/model"
	"chat/pkg/e"
	"chat/serializer"
	"github.com/jinzhu/gorm"
)

type FriendService struct {
	Uid    uint  `form:"uid" json:"uid" `
	IsPass int64 `form:"is_pass" json:"is_pass"`
}

// Get 获取好友列表
func (service *FriendService) Get(id uint) serializer.Response {
	var user model.User
	code := e.SUCCESS
	err := model.DB.Where("id=?", id).First(&user).Error
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	var friends []model.Friend
	var users []model.User
	err1 := model.DB.Where("is_pass=1").Where("user_one_id=? or user_two_id=?", id, id).Find(&friends).Error
	if err1 != nil {
		if gorm.IsRecordNotFoundError(err) {
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				//Error : err.Error(),
			}
		}
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	var fid uint
	var friend model.User
	for i := 0; i < len(friends); i++ {
		if friends[i].UserOneID == id {
			fid = friends[i].UserTwoID
		} else {
			fid = friends[i].UserOneID
		}
		model.DB.Find(&friend, fid)
		users = append(users, friend)
	}
	return serializer.BuildListResponse(serializer.BuildUsers(users), uint(len(friends)))
}

// Make 添加好友
func (service *FriendService) Make(id uint) serializer.Response {
	var user model.User
	code := e.SUCCESS
	err := model.DB.First(&user, service.Uid).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			code = e.ErrorNotExistUser
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				//Error : err.Error(),
			}
		}
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	var friend model.Friend
	var count1, count2 int
	model.DB.Where("user_one_id=?", service.Uid).Where("user_two_id=?", id).First(&friend).Count(&count1)
	model.DB.Where("user_one_id=?", id).Where("user_two_id=?", service.Uid).First(&friend).Count(&count2)
	if count1 == 1 || count2 == 1 {
		code = e.ErrorExistFriendship
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			//Error : err.Error(),
		}
	}
	friend.UserOneID = id
	friend.UserTwoID = service.Uid
	friend.IsPass = 0
	err1 := model.DB.Create(&friend).Error
	if err1 != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err1.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

// GetApplication 获取好友申请
func (service *FriendService) GetApplication(id uint) serializer.Response {
	var user model.User
	code := e.SUCCESS
	err := model.DB.Where("id=?", id).First(&user).Error
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	var friends []model.Friend
	var users []model.User
	err1 := model.DB.Where("is_pass=0").Where("user_two_id=?", id).Find(&friends).Error
	if err1 != nil {
		if gorm.IsRecordNotFoundError(err) {
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				//Error : err.Error(),
			}
		}
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	for i := 0; i < len(friends); i++ {
		model.DB.Where("id=?", friends[i].UserOneID).First(&user)
		users = append(users, user)
	}
	return serializer.BuildListResponse(serializer.BuildUsers(users), uint(len(friends)))
}

// Reply 回复申请
func (service *FriendService) Reply(id uint) serializer.Response {
	code := e.SUCCESS
	var friend model.Friend
	err := model.DB.Where("user_two_id=?", id).Where("user_one_id=?", service.Uid).First(&friend).Error
	if err != nil {
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	if service.IsPass == 1 {
		friend.IsPass = service.IsPass
		model.DB.Save(&friend)
	} else {
		model.DB.Delete(&friend)
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

// Message 获取聊天记录
func (service *FriendService) Message(id uint) serializer.Response {
	var msg []model.SingleMessage
	code := e.SUCCESS
	count := 0
	err := model.DB.Where("(to_id=? AND from_id=?) or(to_id=? AND from_id=?)", id, service.Uid, service.Uid, id).Find(&msg).Count(&count).Error
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.BuildListResponse(serializer.BuildSingleMessages(msg), uint(count))
}
