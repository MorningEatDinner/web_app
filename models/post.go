package models

import "time"

type Post struct {
	ID          int64     `json:"id" gorm:"column:post_id"`
	AuthorID    int64     `json:"author_id" gorm:"column:author_id"`
	CommunityID int64     `json:"community_id" gorm:"column:community_id;not null"`
	Status      int32     `json:"status" gorm:"column:status"`
	Title       string    `json:"title" gorm:"column:title;not null"`
	Content     string    `json:"content" gorm:"column:content;not null"`
	CreateTime  time.Time `json:"-" gorm:"column:create_time;autoCreateTime"`
	UpdatedTime time.Time `json:"-" gorm:"column:updated_time;autoUpdateTime"`
}

// 帖子详情结构的结构体 设置api接口专用的模型
type ApiPostDetail struct {
	AuthorName string `json:"author_name"`
	//VoteNum    int64  `json:"vote_num"`
	*Post
	*Community `json:"community"`
}
type ApiPostDetail2 struct {
	AuthorName string `json:"author_name"`
	VoteNum    int64  `json:"vote_num"`
	*Post
	*Community `json:"community"`
}
