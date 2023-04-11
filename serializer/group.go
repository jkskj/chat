package serializer

import "chat/model"

type Group struct {
	ID      uint   `json:"id" form:"id" example:"1"`                    // 群聊ID
	Name    string `json:"user_name" form:"user_name" example:"FanOne"` // 群聊名
	Creator User   `json:"creator"`                                     //创建者
}

// BuildGroup 序列化群聊
func BuildGroup(group model.Group) Group {
	var user model.User
	model.DB.Where("id=?", group.Creator).First(&user)
	creator := BuildUser(user)
	return Group{
		ID:      group.ID,
		Name:    group.Name,
		Creator: creator,
	}
}
func BuildGroups(items []model.Group) (Groups []Group) {
	for _, item := range items {
		group := BuildGroup(item)
		Groups = append(Groups, group)
	}
	return Groups
}
