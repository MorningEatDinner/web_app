package models

const (
	OrderTime  = "time"
	OrderScore = "score"
)

type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ParamVoteData struct {
	PostID    int64 `json:"post_id,string" binding:"required"`
	Direction int8  `json:"direction,string" binding:"oneof=0 1 -1"` // required会把一些零值给看做没有值， 比如0对于int
}

//type ParamPostList struct {
//	Page  int64  `json:"page" form:"page"`
//	Size  int64  `json:"size" form:"size"`
//	Order string `json:"order" form:"order"`
//}
//
//type ParamCommunityPostList struct {
//	*ParamPostList
//	CommunityID int64 `json:"community_id" form:"community_id"`
//}

type ParamPostList struct {
	CommunityID int64  `json:"community_id" form:"community_id"`
	Page        int64  `json:"page" form:"page"`
	Size        int64  `json:"size" form:"size"`
	Order       string `json:"order" form:"order"`
}

type ParamPhoneExist struct {
	Phone string `json:"phone,omitempty" valid:"phone"`
}
