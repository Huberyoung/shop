package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"

	"shop_srvs/pkg/setting"
)

func NewMysqlDBEngine(setting *setting.DatabaseSettingS) (*gorm.DB, error) {
	s := "%s:%s@%s(%s:%s)/%s?charset=%s&parseTime=%t&loc=%s"
	dsn := fmt.Sprintf(s, setting.Username, setting.Password, setting.Protocol, setting.Host, setting.Port,
		setting.DBName, setting.Charset, setting.ParseTime, setting.Loc)

	newLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{Colorful: true, LogLevel: logger.Info})
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		return nil, err
	}
	return db.Set("gorm:table_options", "charset="+setting.Charset), nil
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
