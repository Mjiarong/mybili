package serializer

import (
	"mybili/utils"
)

// Response 基础序列化器
type Response struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data,omitempty"`
	Msg   string      `json:"msg"`
	Error string      `json:"error,omitempty"`
}

// TrackedErrorResponse 有追踪信息的错误响应
type TrackedErrorResponse struct {
	Response
	TrackID string `json:"track_id"`
}

// CheckToken 检查令牌
func CheckToken(code int, msg string) Response {
	return Response{
		Code: code,
		Msg:  msg,
	}
}

// DataList 基础列表结构
type DataList struct {
	Items interface{} `json:"items"`
	Total uint        `json:"total"`
}

// BuildListResponse 列表构建器
func BuildListResponse(items interface{}, total uint) Response {
	return Response{
		Code: utils.SUCCESS,
		Msg:  utils.GetErrMsg(utils.SUCCESS),
		Data: DataList{
			Items: items,
			Total: total,
		},
	}
}
