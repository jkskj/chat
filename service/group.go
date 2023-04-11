package service

import (
	"chat/model"
	"chat/pkg/e"
	"chat/serializer"
	"github.com/jinzhu/gorm"
)

type GroupService struct {
	Gid    uint   `form:"gid" json:"gid"`
	Name   string `form:"name" json:"name"`
	Avatar string `form:"avatar" json:"avatar"`
}

// Get 获取群聊列表
func (service *GroupService) Get(id uint) serializer.Response {
	var joins []model.Join
	var groups []model.Group
	code := e.SUCCESS
	var count int
	err := model.DB.Where("uid=?", id).First(&joins).Count(&count).Error
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	var group model.Group
	for i := 0; i < len(joins); i++ {
		model.DB.Where("id=?", joins[i].Gid).First(&group)
		groups = append(groups, group)
	}
	return serializer.BuildListResponse(serializer.BuildGroups(groups), uint(count))
}

// Make 创建群聊
func (service *GroupService) Make(id uint) serializer.Response {
	var group model.Group
	code := e.SUCCESS
	group.Name = service.Name
	group.Creator = id
	group.Avatar = service.Avatar
	err := model.DB.Create(&group).Error
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	var join model.Join
	join.Gid = group.ID
	join.Uid = id
	err1 := model.DB.Create(&join).Error
	if err1 != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

// Join 加入群聊
func (service *GroupService) Join(id uint) serializer.Response {
	var join model.Join
	code := e.SUCCESS
	var count int
	var group model.Group
	err := model.DB.First(&group, service.Gid).Count(&count).Error
	if err != nil {
		//群聊是否存在
		if gorm.IsRecordNotFoundError(err) {
			code = e.ErrorNotExistGroup
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
	count = 0
	model.DB.Where("gid=?", service.Gid).Where("uid=?", id).First(&join).Count(&count)
	if count == 1 {
		code = e.ErrorHaveJoinGroup
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	join.Gid = service.Gid
	join.Uid = id
	err = model.DB.Create(&join).Error
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	group.MemberNum++
	model.DB.Save(&group)
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

// Message 获取聊天记录
func (service *GroupService) Message() serializer.Response {
	var msg []model.GroupMessage
	code := e.SUCCESS
	count := 0
	err := model.DB.Where("group_id=?", service.Gid).Find(&msg).Count(&count).Error
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.BuildListResponse(serializer.BuildGroupMessages(msg), uint(count))
}
