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

type ParamEmailExist struct {
	Email string `json:"email,omitempty" valid:"email"`
}

type ParamPhoneCode struct {
	CaptchaID     string `json:"captcha_id,omitempty" valid:"captcha_id"`
	CaptchaAnswer string `json:"captcha_answer,omitempty" valid:"captcha_answer"`
	Phone         string `json:"phone,omitempty" valid:"phone"`
}

type ParamEmailCode struct {
	CaptchaID     string `json:"captcha_id,omitempty" valid:"captcha_id"`
	CaptchaAnswer string `json:"captcha_answer,omitempty" valid:"captcha_answer"`
	Email         string `json:"email,omitempty" valid:"email"`
}

type ParamSignupUsingPhone struct {
	Phone           string `json:"phone,omitempty" valid:"phone"`
	Code            string `json:"code,omitempty" valid:"verify_code"`
	Name            string `json:"name" valid:"name"`
	Password        string `json:"password" valid:"password"`
	PasswordConfirm string `json:"password_confirm" valid:"password_confirm"`
}

type ParamSignUpUsingEmail struct {
	Email           string `json:"email,omitempty" valid:"email"`
	Code            string `json:"code,omitempty" valid:"verify_code"`
	Name            string `json:"name" valid:"name"`
	Password        string `json:"password" valid:"password"`
	PasswordConfirm string `json:"password_confirm" valid:"password_confirm"`
}

type ParamLoginUsingPhoneWithCode struct {
	Phone string `json:"phone,omitempty" valid:"phone"`
	Code  string `json:"code,omitempty" valid:"verify_code"`
}
