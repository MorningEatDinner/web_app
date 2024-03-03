package logic

import (
	"errors"
	"fmt"
	"github.com/xiaorui/web_app/dao/mysql"
	"github.com/xiaorui/web_app/dao/redis"
	"github.com/xiaorui/web_app/models"
	"github.com/xiaorui/web_app/pkg/helpers"
	"github.com/xiaorui/web_app/pkg/jwt"
	"github.com/xiaorui/web_app/pkg/mail"
	"github.com/xiaorui/web_app/pkg/sms"
	"github.com/xiaorui/web_app/pkg/snowflake"
	"github.com/xiaorui/web_app/settings"
	"go.uber.org/zap"
)

// 存放业务逻辑的代码
func SignUp(p *models.ParamSignUp) (err error) {
	//1. 判断用户是否存在
	if err = mysql.CheckUserExist(p.Username); err != nil {
		return err
	}

	//2. 生成uid
	userID := snowflake.GenID()
	//3. 密码加密

	//构造数据实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}

	//4. 保存进数据库

	err = mysql.InsertUser(user)

	//这里还可以有很多其他的数据操作， 比如对于redis进行操作

	return
}

func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	if err = mysql.Login(user); err != nil {
		return nil, err
	}
	//如果登录成功
	//生成JWT
	//return jwt.GenToken(user.UserID, user.Username)
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return nil, err
	}
	user.Token = token
	return
}

// LoginUsingPhoneWithCode: 使用手机+验证码的形式进行登陆
func LoginUsingPhoneWithCode(p *models.ParamLoginUsingPhoneWithCode) (*models.User, string, error) {
	user := &models.User{
		Phone: p.Phone,
	}
	// 进行登陆操作
	if err := mysql.LoginUsingPhoneWithCode(user); err != nil {
		return nil, "", err
	}

	// 如果登录成功：即确实有这个用户存在
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

// IsPhoneExist：返回输入手机号码是否存在数据表中
func IsPhoneExist(phone string) (bool, error) {
	return mysql.IsPhoneExist(phone)
}

// IsEmailExist：返回输入邮箱是否存在数据表中
func IsEmailExist(phone string) (bool, error) {
	return mysql.IsEmailExist(phone)
}

// SendPhoneCode: 发送短信验证码
func SendPhoneCode(phone string) error {
	// 1. 生成验证码
	code := helpers.GenerateRandomCode()

	// 2. 将验证码保存到redis中
	if err := redis.SetVerifyCode(phone, code); err != nil {
		zap.L().Error("SetVerifyCode failed...", zap.Error(err))
		return err
	}
	// 3. 发送短信给手机
	if ok := sms.NewSms().Send(phone, code); !ok {
		return errors.New("短信发送失败")
	}
	return nil
}

// SendEmailCode: 发送邮箱验证码
func SendEmailCode(email string) error {
	// 1. 生成验证码
	code := helpers.GenerateRandomCode()

	// 2. 将验证码存放到redis中
	if err := redis.SetVerifyCode(email, code); err != nil {
		zap.L().Error("SetVerifyCode failed...", zap.Error(err))
		return err
	}
	// 3. 发送验证码给邮箱;
	ok := mail.NewMailer().Send(
		mail.Email{
			From: mail.From{
				settings.Conf.EmailConfig.FromConfig.Address,
				settings.Conf.EmailConfig.FromConfig.Name,
			},
			To:      []string{email},
			Subject: "email 验证码",
			HTML:    []byte(fmt.Sprintf("<h1>您的 Email 验证码是 %v </h1>", code)),
		},
	)
	if !ok {
		return errors.New("邮箱验证码发送失败")
	}
	return nil
}

// SignupUsingPhone：处理手机注册登陆逻辑
func SignupUsingPhone(p *models.ParamSignupUsingPhone) (err error) {
	// 1. 判断用户是否存在
	if err = mysql.CheckUserExist(p.Name); err != nil {
		return err
	}

	// 判断手机号码是否存在
	if _, err = mysql.IsPhoneExist(p.Phone); err != nil {
		return err
	}
	// 2. 生成uid
	userID := snowflake.GenID()

	// 3. 构造用户实例
	user := &models.User{
		UserID:   userID,
		Username: p.Name,
		Password: p.Password,
		Phone:    p.Phone,
	}

	// 4. 保存到数据库
	err = mysql.InsertUser(user)

	return
}

// SignUpUsingEmail: 进行使用邮箱进行注册的业务
func SignUpUsingEmail(p *models.ParamSignUpUsingEmail) (err error) {
	// 1. 验证用户名是否存在
	if err = mysql.CheckUserExist(p.Name); err != nil {
		return
	}
	// 2. 验证邮箱是否已经注册
	if _, err = mysql.IsEmailExist(p.Email); err != nil {
		return
	}
	// 3. 生成uid
	userID := snowflake.GenID()

	// 4. 构造用户实例
	_user := models.User{
		UserID:   userID,
		Username: p.Name,
		Email:    p.Email,
		Password: p.Password,
	}

	// 5. 保存到数据库中
	return mysql.InsertUser(&_user)

}
