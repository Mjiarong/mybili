package cache

import (
	"fmt"
	"strconv"
)

const (
	// DailyRankKey 每日排行
	DailyRankKey = "rank:daily"
)

// VideoViewKey 视频点击数的key
// view:video:1 -> 100
// view:video:2 -> 150
func VideoViewKey(id uint) string {
	return fmt.Sprintf("view:video:%s", strconv.Itoa(int(id)))
}

// VideoLikesKey 视频点赞数的key
// likes:video:1 -> 100
// likes:video:2 -> 150
func VideoLikesKey(id uint) string {
	return fmt.Sprintf("likes:video:%s", strconv.Itoa(int(id)))
}

// VideoCommentKey 视频评论数的key
// comment:video:1 -> 100
// comment:video:2 -> 150
func VideoCommentKey(id uint) string {
	return fmt.Sprintf("comment:video:%s", strconv.Itoa(int(id)))
}

// CommentLikesKey 评论点赞数的key
// likes:comment:1 -> 100
// likes:comment:2 -> 150
func CommentLikesKey(id uint) string {
	return fmt.Sprintf("likes:comment:%s", strconv.Itoa(int(id)))
}

// CommentDislikesKey 评论点踩数的key
// dislikes:comment:1 -> 100
// dislikes:comment:2 -> 150
func CommentDislikesKey(id uint) string {
	return fmt.Sprintf("dislikes:comment:%s", strconv.Itoa(int(id)))
}
