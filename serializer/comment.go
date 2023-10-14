package serializer

import (
	"mybili/model"
)

// Video 视频评论序列化器
type Comment struct {
	ID            uint   `json:"id"`
	Content       string `json:"content"`
	UserId        uint   `json:"user_id""`
	VideoId       uint   `json:"video_id"`
	ParentId      uint   `json:"parent_id"`
	ReplyUserName string `json:"reply_user_name"`
	CreatedAt     int64  `json:"created_at"`
	/*	UserName      string    `json:"user_name"`
		Nickname      string    `json:"nickname"`
		UserAvatarUrl string    `json:"user_avatar_url"`
		Likes         int64     `json:"likes"`
		Children      []Comment `json:"children"`
		ChildrenNum   uint      `json:"children_num"`
		Liked         bool      `json:"liked"`    //用户是否对该评论有点赞
		Disliked      bool      `json:"disliked"` //用户是否对该评论有点赞*/
}

// BuildComment 评论序列化视频
func BuildComment(item model.Comment) Comment {
	return Comment{
		ID:            item.ID,
		Content:       item.Content,
		UserId:        item.UserID,
		VideoId:       item.VideoID,
		ParentId:      item.ParentId,
		ReplyUserName: item.ReplyUserName,
		CreatedAt:     item.CreatedAt.Unix(),
	}
}

type derivedComment struct {
	//结构体成员必须以大写开头
	Comment
	UserName      string           `json:"user_name"`
	Nickname      string           `json:"nickname"`
	UserAvatarUrl string           `json:"user_avatar_url"`
	Likes         int64            `json:"likes"`
	Liked         bool             `json:"liked"` //用户是否对该评论有点赞
	Disliked      bool             `json:"disliked"`
	Children      []derivedComment `json:"children"`
	ChildrenNum   uint             `json:"children_num"`
}

// BuildDerivedComment 与用户信息联合查询时评论序列化视频
func BuildDerivedComment(item model.DerivedComment) derivedComment {
	return derivedComment{
		Comment: Comment{
			item.ID,
			item.Content,
			item.UserID,
			item.VideoID,
			item.ParentId,
			item.ReplyUserName,
			item.CreatedAt.Unix(),
		},
		UserName:      item.UserName,
		Nickname:      item.Nickname,
		UserAvatarUrl: item.AvatarURL(),
		Likes:         0,
		Children:      nil,
		ChildrenNum:   0,
		Liked:         false,
		Disliked:      false,
	}
}

func BuildComments(items []model.DerivedComment, UserId uint) (comments []derivedComment, total uint) {
	//cmap := make(map[uint]*Comment, 10)
	indexMap := make(map[uint]int, 10)
	rootIndex := 0
	for _, item := range items {
		total += 1
		if item.ParentId == 0 {
			comment := BuildDerivedComment(item)
			comment.Likes = item.Likes()
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
			comment := BuildDerivedComment(item)
			comment.Liked = item.Liked(UserId)
			comment.Disliked = item.Disliked(UserId)
			comments[index].Children = append(comments[index].Children, comment)
			comments[index].ChildrenNum += 1
		}
	}
	return
}
