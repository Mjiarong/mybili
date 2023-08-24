package serializer

import (
	"mybili/model"
)

// Video 视频序列化器
type Video struct {
	ID          uint    `json:"id"`
	Title       string  `json:"title"`
	Info        string  `json:"info"`
	AvatarURL   string  `json:"avatar_url"`
	VideoURL    string  `json:"video_url"`
	AvatarKey   string  `json:"avatar_key"`
	VideoKey    string  `json:"video_key"`
	View        uint64  `json:"view"`
	User        User    `json:"user"`
	CreatedAt   int64   `json:"created_at"`
	Duration    float32 `json:"duration"`
	CommentsNum uint64  `json:"comments_num"`
	LikesNum    int64   `json:"likes_num"`
}

// BuildVideo 序列化视频
func BuildVideo(item model.Video) Video {
	user, _ := model.GetUserByName(item.Creator)
	return Video{
		ID:          item.ID,
		Title:       item.Title,
		Info:        item.Info,
		Duration:    item.Duration,
		AvatarKey:   item.AvatarKey,
		VideoKey:    item.VideoKey,
		VideoURL:    item.VideoURL(),
		AvatarURL:   item.AvatarURL(),
		User:        BuildUser(user),
		CreatedAt:   item.CreatedAt.Unix(),
		View:        item.View(),
		CommentsNum: item.GetCommentsNum(),
		LikesNum:    item.GetLikesNum(),
	}
}

// BuildVideos 序列化视频列表
func BuildVideos(items []model.Video) (videos []Video) {
	for _, item := range items {
		video := BuildVideo(item)
		videos = append(videos, video)
	}
	return videos
}
