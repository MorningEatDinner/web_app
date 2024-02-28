package models

type User struct {
	UserID   int64  `json:"user_id,string" gorm:"column:user_id"`
	Username string `json:"username" gorm:"column:username"`
	Password string `json:"password" gorm:"column:password"`
	Token    string // 注意：此字段没有json或gorm标签
}
