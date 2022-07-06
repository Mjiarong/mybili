package errmsg

const(
	SUCCESS = 0
	ERROR = -1
	//code = 1000... 用户模块的错误

	//code = 2000... 数据库错误
	DB_CONNECT_FAILED= 2001

	//code = 3000... 视频模块的错误
	ERROR_VIDEO_CREATING = 3001
	ERROR_VIDEO_NOEXIST = 3002
	ERROR_VIDEO_SAVE_FAILED = 3003
	ERROR_VIDEO_DELETE_FAILED = 3004

)

var CodeMsg = map[int]string{
	SUCCESS:                  "OK",
	ERROR:                    "FAIL",

	DB_CONNECT_FAILED:      "数据库连接错误",

	ERROR_VIDEO_CREATING:      "视频创建失败",
	ERROR_VIDEO_NOEXIST:      "视频不存在",
	ERROR_VIDEO_SAVE_FAILED:      "视频保存失败",
	ERROR_VIDEO_DELETE_FAILED:      "视频删除失败",
}

func GetErrMsg(code int)string{
	return CodeMsg[code]
}