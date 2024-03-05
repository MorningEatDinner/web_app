package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaorui/web_app/dao/mysql"
	"github.com/xiaorui/web_app/logic"
	"github.com/xiaorui/web_app/models"
	"github.com/xiaorui/web_app/pkg/file"
	"go.uber.org/zap"
)

// CurrentUser: 获得当前用户信息
func CurrentUser(ctx *gin.Context) {
	userID, err := getCurrentUser(ctx)
	if err != nil {
		if err == ErrorUserNotLogin {
			ResponseError(ctx, CodeNeedLogin)
			return
		}
		ResponseError(ctx, CodeServerBusy)
		return
	}
	user, err := mysql.GetUserByID(userID)
	if err != nil {
		ResponseError(ctx, CodeServerBusy)
		return
	}
	ResponseSuccess(ctx, user)
}

// UpdateProfile:更改用户个人信息
func UpdateProfile(ctx *gin.Context) {
	// 1. 进行参数验证
	p := new(models.ParamUpdateProfile)
	if ok := Validate(ctx, p, ValidateUpdateProfile); !ok {
		return
	}

	// 2. 处理业务逻辑， 修改用户信息
	userID, err := getCurrentUser(ctx)
	if err != nil {
		ResponseError(ctx, CodeServerBusy)
		return
	}
	if user, err := logic.UpdateProfile(p, userID); err != nil {
		zap.L().Error("logic.UpdateProfile", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
	} else {
		// 3. 返回响应
		ResponseSuccess(ctx, user)
	}
}

// UpdateEmail: 更改邮箱
func UpdateEmail(ctx *gin.Context) {
	// 1. 接受参数: 验证码+新的邮箱
	p := new(models.ParamUpdateEmail)
	if ok := Validate(ctx, p, ValidateUpdateEmail); !ok {
		ResponseError(ctx, CodeServerBusy)
		return
	}
	// 2. 处理业务逻辑：更改邮箱
	userID, err := getCurrentUser(ctx)
	if err != nil {
		ResponseError(ctx, CodeServerBusy)
		return
	}
	if user, err := logic.UpdateEmail(p, userID); err != nil {
		if err == mysql.ErrorEmailExist {
			ResponseError(ctx, CodeEmailExist)
			return
		}
		ResponseError(ctx, CodeServerBusy)
		return
	} else {
		// 3. 返回响应
		ResponseSuccess(ctx, user)
	}
}

// UpdatePhone: 更改手机号码
func UpdatePhone(ctx *gin.Context) {
	// 1. 进行参数验证
	p := new(models.ParamUpdatePhone)
	if ok := Validate(ctx, p, ValidateUpdatePhone); !ok {
		return
	}

	// 2. 过去当前yonghu
	userID, err := getCurrentUser(ctx)
	if err != nil {
		ResponseError(ctx, CodeServerBusy)
		return
	}
	//3. 进行业务处理
	if user, err := logic.UpdatePhone(p, userID); err != nil {
		if err == mysql.ErrorPhoneExist {
			ResponseError(ctx, CodePhoneExist)
			return
		}
		ResponseError(ctx, CodeServerBusy)
		return
	} else {
		// 返回成功相应
		ResponseSuccess(ctx, user)
	}

}

// UpdatePassword: 更新用户密码
func UpdatePassword(ctx *gin.Context) {
	// 1. 进行参数验证
	p := new(models.ParamUpdatePassword)
	if ok := Validate(ctx, p, ValidateUpdatePassword); !ok {
		return
	}

	// 2. 处理业务逻辑
	// 获取用户
	userID, err := getCurrentUser(ctx)
	if err != nil {
		ResponseError(ctx, CodeServerBusy)
		return
	}
	// 进行业务逻辑
	if err = logic.UpdatePassword(p, userID); err != nil {
		if err == mysql.ErrorPasswordInvalid {
			ResponseError(ctx, CodeInvalidPassword)
			return
		}
		ResponseError(ctx, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(ctx, nil)
}

// UpdateAvatar：更新用户头像
func UpdateAvatar(ctx *gin.Context) {
	// 1. 验证请求
	p := new(models.ParamUpdateAvatar)
	if ok := Validate(ctx, p, ValidateUpdateAvatar); !ok {
		return
	}

	// 2. 处理业务：上传头像
	userID, err := getCurrentUser(ctx)
	if err != nil {
		ResponseError(ctx, CodeServerBusy)
		return
	}
	if user, err := file.SaveUploadAvatar(ctx, p.Avatar, userID); err != nil {
		ResponseError(ctx, CodeServerBusy)
		return
	} else {
		//3. 返回响应
		ResponseSuccess(ctx, user)
	}
}
