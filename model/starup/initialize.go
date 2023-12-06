package starup

import (
	"shop_srvs/global"
	"shop_srvs/model/database"
	"shop_srvs/pkg/setting"
)

// SetUpSetting 读取一些配置文件
func SetUpSetting() error {
	newSetting, err := setting.NewSetting()
	if err != nil {
		return err
	}
	err = newSetting.ReadSection("DataBase", &global.DatabaseSetting)
	if err != nil {
		return err
	}
	return nil
}

// SetUpDBEngine 启动数据库
func SetUpDBEngine() error {
	engine, err := database.NewMysqlDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}
	global.DBEngine = engine
	return nil
}
