package errmsg

const(
	SUCCESS = 0
	ERROR = -1
	VALIDATION_ERROR=100
	UNMARSHAL_TYPE_ERROR=102
	PARAM_ERROR=103

	//code = 1000... 用户模块的错误
	PASSWORD_ENTERED_DIFFERENT=1001
	NICKNAME_OCCUPIED=1002
	USERNAME_REGISTERED=1003
	PASSWORD_ENCRYPT_FAILED=1004
	REGISTER_FAILED=1005
	ACCOUNT_OR_PASSWORD_INCORRECT=1006

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
	UNMARSHAL_TYPE_ERROR:"JSON类型不匹配",
	PARAM_ERROR:"参数错误",

	PASSWORD_ENTERED_DIFFERENT:	"两次输入的密码不相同",
	NICKNAME_OCCUPIED:"昵称被占用",
	USERNAME_REGISTERED:"用户名已经注册",
	PASSWORD_ENCRYPT_FAILED:"密码加密失败",
	REGISTER_FAILED:"注册失败",
	ACCOUNT_OR_PASSWORD_INCORRECT:"账号或密码错误",

	DB_CONNECT_FAILED:      "数据库连接错误",

	ERROR_VIDEO_CREATING:      "视频创建失败",
	ERROR_VIDEO_NOEXIST:      "视频不存在",
	ERROR_VIDEO_SAVE_FAILED:      "视频保存失败",
	ERROR_VIDEO_DELETE_FAILED:      "视频删除失败",
}

func GetErrMsg(code int)string{
	return CodeMsg[code]
}