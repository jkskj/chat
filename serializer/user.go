package serializer

import "chat/model"

type User struct {
	ID       uint   `json:"id" form:"id" example:"1"`                    // 用户ID
	UserName string `json:"user_name" form:"user_name" example:"FanOne"` // 用户名
	CreateAt int64  `json:"create_at" form:"create_at"`                  // 创建
	Avatar   string `json:"avatar"`
}

// BuildUser 序列化用户
func BuildUser(user model.User) User {
	return User{
		ID:       user.ID,
		UserName: user.UserName,
		CreateAt: user.CreatedAt.Unix(),
		Avatar:   user.Avatar,
	}
}
func BuildUsers(items []model.User) (Users []User) {
	for _, item := range items {
		user := BuildUser(item)
		Users = append(Users, user)
	}
	return Users
}
