package serializer

import (
	"mybili/model"
)

// Video 视频评论序列化器
type Comment struct {
	ID            uint      `json:"id"`
	Content       string    `json:"content"`
	UserId        uint      `json:"user_id""`
	UserName      string    `json:"user_name"`
	Nickname      string    `json:"nickname"`
	UserAvatarUrl string    `json:"user_avatar_url"`
	VideoId       uint      `json:"video_id"`
	ParentId      uint      `json:"parent_id"`
	ReplyUserName string    `json:"reply_user_name"`
	Likes         int64     `json:"likes"`
	CreatedAt     int64     `json:"created_at"`
	Children      []Comment `json:"children"`
	ChildrenNum   uint      `json:"children_num"`
	Liked         bool      `json:"liked"`    //用户是否对该评论有点赞
	Disliked      bool      `json:"disliked"` //用户是否对该评论有点赞
}

// BuildComment 评论序列化视频
func BuildComment(item model.Comment) Comment {
	return Comment{
		ID:            item.ID,
		Content:       item.Content,
		UserId:        item.UserID,
		UserName:      item.UserName,
		Nickname:      item.Nickname,
		UserAvatarUrl: item.AvatarURL(),
		VideoId:       item.VideoID,
		ParentId:      item.ParentId,
		ReplyUserName: item.ReplyUserName,
		CreatedAt:     item.CreatedAt.Unix(),
		Likes:         item.Likes(),
		Children:      nil,
		ChildrenNum:   0,
		Liked:         false,
		Disliked:      false,
	}
}

func BuildComments(items []model.Comment, UserId uint) (comments []Comment, total uint) {
	//cmap := make(map[uint]*Comment, 10)
	indexMap := make(map[uint]int, 10)
	rootIndex := 0
	for _, item := range items {
		total += 1
		if item.ParentId == 0 {
			comment := BuildComment(item)
			comment.Liked = item.Liked(UserId)
			comment.Disliked = item.Disliked(UserId)
			indexMap[item.ID] = rootIndex
			comments = append(comments, comment)
			rootIndex += 1
		}
	}

	for _, item := range items {
		if item.ParentId > 0 {
			index, ok := indexMap[item.ParentId]
			if !ok {
				break
			}
			comment := BuildComment(item)
			comment.Liked = item.Liked(UserId)
			comment.Disliked = item.Disliked(UserId)
			comments[index].Children = append(comments[index].Children, comment)
			comments[index].ChildrenNum += 1
		}
	}
	return
}
