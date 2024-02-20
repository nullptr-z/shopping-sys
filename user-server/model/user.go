package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        int32
	CratedAt  time.Time      `gorm:"column:create_time"`
	UpdatedAt time.Time      `gorm:"column:update_time"`
	DeleteAt  gorm.DeletedAt `gorm:"column:delete_time"`
	IsDelete  bool
}

/*
1. 密码
*/
type User struct {
	BaseModel
	// gorm.Model
	Mobile   string     `gorm:"index:idx_mobile;unique;type:varchar(11);not null"`
	Password string     `gorm:"type:varchar(100);not null"`
	NickName string     `gorm:"type:varchar(30)"`
	Birthday *time.Time `gorm:"type:datetime"`
	Gender   string     `gorm:"type:varchar(6);default:man;comment:'man 男, woman 女'"`
	Role     int        `gorm:"default:1;type:int;comment:'1 用户, 2 管理员'"`
}
