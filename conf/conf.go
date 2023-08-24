package conf

import (
	"github.com/joho/godotenv"
	"mybili/cache"
	"mybili/model"
	"mybili/task"
	"mybili/utils"
	"os"
)

// Init 初始化配置项
func Init() {
	// 从本地读取环境变量
	godotenv.Load()

	// 日志模块
	utils.Init()
	// 设置日志级别
	/*	utils.BuildLogger(os.Getenv("LOG_LEVEL"))

		// 读取翻译文件
		if err := LoadLocales("conf/locales/zh-cn.yaml"); err != nil {
			utils.Log().Panic("翻译文件加载失败", err)
		}*/

	// 连接数据库
	model.Database(os.Getenv("MYSQL_DSN"))
	cache.Redis()

	//启动定时任务
	task.CronJob()
}
