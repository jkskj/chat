package serializer

import "chat/model"

type GroupMessage struct {
	Id          int64  `json:"id"`
	GroupId     int64  ` json:"group_id"`
	User        User   `json:"user"`
	MessageData string ` json:"message_data"`
	CreateTime  string `json:"create_time"`
}

type SingleMessage struct {
	Id          int64  `json:"id"`
	ToUser      User   `json:"to_user"`
	FromUser    User   `json:"from_user"`
	MessageData string `json:"message_data"`
	CreateTime  string `json:"create_time"`
}

func BuildGroupMessage(groupMessage model.GroupMessage) GroupMessage {
	var user model.User
	model.DB.First(&user, groupMessage.UserId)
	return GroupMessage{
		Id:          groupMessage.Id,
		User:        BuildUser(user),
		GroupId:     groupMessage.GroupId,
		MessageData: groupMessage.MessageData,
		CreateTime:  groupMessage.CreateTime,
	}
}

func BuildGroupMessages(items []model.GroupMessage) (groupMessages []GroupMessage) {
	for _, item := range items {
		groupMessage := BuildGroupMessage(item)
		groupMessages = append(groupMessages, groupMessage)
	}
	return groupMessages
}

func BuildSingleMessage(singleMessage model.SingleMessage) SingleMessage {
	var fromUser, toUser model.User
	model.DB.First(&fromUser, singleMessage.FromId)
	model.DB.First(&toUser, singleMessage.ToId)
	return SingleMessage{
		Id:          singleMessage.Id,
		FromUser:    BuildUser(fromUser),
		ToUser:      BuildUser(toUser),
		MessageData: singleMessage.MessageData,
		CreateTime:  singleMessage.CreateTime,
	}
}

func BuildSingleMessages(items []model.SingleMessage) (singleMessages []SingleMessage) {
	for _, item := range items {
		singleMessage := BuildSingleMessage(item)
		singleMessages = append(singleMessages, singleMessage)
	}
	return singleMessages
}
