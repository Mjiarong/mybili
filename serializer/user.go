package serializer

import "mybili/model"

// User 用户序列化器
type User struct {
	ID        uint         `json:"id"`
	UserName  string       `json:"user_name"`
	Nickname  string       `json:"nickname"`
	Status    model.Status `json:"status"`
	AvatarKey string       `json:"avatar_key"`
	Avatar    string       `json:"avatar"`
	CreatedAt int64        `json:"created_at"`
}

// BuildUser 序列化用户
func BuildUser(user model.User) User {
	return User{
		ID:        user.ID,
		UserName:  user.UserName,
		Nickname:  user.Nickname,
		Status:    user.Status,
		AvatarKey: user.Avatar,
		Avatar:    user.AvatarURL(),
		CreatedAt: user.CreatedAt.Unix(),
	}
}

// BuildUserResponse 序列化用户响应
func BuildUserResponse(user model.User, code int, msg string) Response {
	return Response{
		Code: code,
		Data: BuildUser(user),
		Msg:  msg,
	}
}
