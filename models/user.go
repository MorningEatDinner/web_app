package models

import "time"

type User struct {
	UserID       int64     `json:"user_id,string" gorm:"column:user_id"`
	Username     string    `json:"username" gorm:"column:username"`
	Password     string    `json:"password" gorm:"column:password"`
	Email        string    `json:"email" gorm:"column:email"`
	Phone        string    `json:"phone" gorm:"column:phone"`
	City         string    `json:"city" gorm:"column:city"`
	Introduction string    `json:"introduction" gorm:"column:introduction"`
	Avatar       string    `json:"avatar" gorm:"column:avatar"`
	Token        string    // 注意：此字段没有json或gorm标签
	CreateTime   time.Time `json:"create_time" gorm:"column:create_time;autoCreateTime"`
	UpdatedTime  time.Time `json:"updated_time" gorm:"column:updated_time;autoUpdateTime"`
}
