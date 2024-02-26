package logic

import (
	"github.com/xiaorui/web_app/dao/mysql"
	"github.com/xiaorui/web_app/models"
	"github.com/xiaorui/web_app/pkg/jwt"
	"github.com/xiaorui/web_app/pkg/snowflake"
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
