package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"

	"github.com/xiaorui/web_app/models"
)

const secret = "xiaorui.com"

// CheckUserExist: 检查制定用户名是否存在
func CheckUserExist(username string) (err error) {
	var count int64
	// sqlStr := `select count(user_id) from user where username = ?`
	res := DB.Model(&models.User{}).Where("username = ?", username).Count(&count)
	if err = res.Error; err != nil {
		return
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}

// InsertUser: 向数据库中添加一条新的用户数据
func InsertUser(user *models.User) (err error) {
	//在插入数据前， 需要对密码进行加密处理
	user.Password = encryptPassword(user.Password)
	// 将数据实例插入数据表中
	// sqlStr := `insert into user(user_id, username, password) values(?,?,?)`
	// _, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	res := DB.Create(user)
	err = res.Error
	return
}

func encryptPassword(oPasssword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPasssword)))
}

func Login(user *models.User) (err error) {
	//这里进行用户登陆， 也就是根据用户名查询用户数据， 之后验证密码是否相等，完成
	// sqlStr := `select user_id, username, password from user where username=?`
	oPassword := user.Password // 记录原始密码， 后面从数据库中返回的密码会叠加在这个数据上
	// err = db.Get(user, sqlStr, user.Username)
	err = DB.Model(&models.User{}).Where("username = ?", user.Username).First(user).Error
	if err == sql.ErrNoRows {
		//var ErrNoRows = errors.New("sql: no rows in result set")
		//想要特别写出来这个错误，否则则原始错误信息返回
		return ErrorUserNotExist
	}
	if err != nil {
		return
	}
	if encryptPassword(oPassword) != user.Password {
		return ErrorPasswordInvalid
	}
	return
}

func GetUserByID(id int64) (user *models.User, err error) {
	user = new(models.User)
	// sqlStr := `select user_id, username from user where user_id=?`

	// err = db.Get(user, sqlStr, id)
	err = DB.Model(&models.User{}).Where("user_id = ?", id).First(user).Error

	return
}

// IsPhoneExist： 验证手机号是否已经注册了
func IsPhoneExist(phone string) (bool, error) {
	// db.Model(&models.User{}).Where("")
	var count int64
	err := DB.Model(&models.User{}).Where("phone = ?", phone).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		return false, ErrorPhoneExist
	}
	return true, err
}

// IsEmailExist: 验证邮箱是否注册了
func IsEmailExist(email string) (bool, error) {
	var count int64
	err := DB.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		return false, ErrorEmailExist
	}
	return true, err
}
