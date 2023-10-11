package model

import "mybili/utils"

//执行数据迁移

func migration() {
	// 自动迁移模式
	err := DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{})
	if err != nil {
		utils.Logger.Errorln(err)
	}
	err = DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&Video{})
	if err != nil {
		utils.Logger.Errorln(err)
	}
	err = DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&Comment{})
	if err != nil {
		utils.Logger.Errorln(err)
	}
}
