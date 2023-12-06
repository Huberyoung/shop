package user

import (
	"time"

	"shop_srvs/model/base"
)

type User struct {
	base.Model
	Mobile   string     `gorm:"index:idx_mobile;unique;type:varchar(11);not null"`
	Password string     `gorm:"type:varchar(200);not null"`
	NickName string     `gorm:"type:varchar(20)"`
	Birthday *time.Time `gorm:"type:datetime"`
	Gender   int        `gorm:"default:1;oneof=0 1;type:int comment '0未知,1男，2女'"`
	Role     int        `gorm:"default:1;type:int comment '1表示普通用户'"`
}
