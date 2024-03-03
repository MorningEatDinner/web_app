package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/xiaorui/web_app/dao/mysql"
	"github.com/xiaorui/web_app/logic"
	"github.com/xiaorui/web_app/models"
	"github.com/xiaorui/web_app/pkg/captcha"
	"go.uber.org/zap"
)

// SignUpHandler: 处理注册登录请求
func SignUpHandler(ctx *gin.Context) {
	//1. 获取参数， 参数校验

	//var p models.ParamSignUp
	p := new(models.ParamSignUp)
	if err := ctx.ShouldBindJSON(p); err != nil {
		//如果出现错误了会进来
		zap.L().Error("Signup with invalid param", zap.Error(err))
		err, ok := err.(validator.ValidationErrors)
		if !ok {
			//ctx.JSON(http.StatusOK, gin.H{
			//	"msg": err.Error(),
			//})
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		//如果是校验器错误
		//ctx.JSON(http.StatusOK, gin.H{
		//	"msg": removeTopStruct(err.Translate(trans)), // 进行错误的翻译
		//})
		ResponseErrorWithMsg(ctx, CodeInvalidParam, removeTopStruct(err.Translate(trans)))
		return
	}

	//手动对请求参数进行业务上的校验， 比如说密码必须满足某些格式等等
	//TODO::使用库进行参数校验而不是手动
	//if len(p.Password) == 0 || len(p.Username) == 0 || len(p.RePassword) == 0 || len(p.Password) != len(p.RePassword) {
	//	zap.L().Error("Signup with invalid param")
	//	ctx.JSON(http.StatusOK, gin.H{
	//		"msg": "请求参数有误",
	//	})
	//	return
	//}

	//2. 业务处理
	//这里有两种错误， 一种是用户已经存在， 一种是其他错误
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("logic.Signup failed", zap.Error(err))
		//ctx.JSON(http.StatusOK, gin.H{
		//	"msg": "注册失败",
		//})
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(ctx, CodeUserExist)
			return
		}
		ResponseError(ctx, CodeServerBusy)
		return
	}
	//3. 返回响应
	ResponseSuccess(ctx, nil)
}

// IsPhoneExist: 判断手机号是否已经被注册
func IsPhoneExist(ctx *gin.Context) {
	//1. 验证参数， 验证手机号是否正确
	p := new(models.ParamPhoneExist)
	if ok := Validate(ctx, p, ValidateSignupPhoneExist); !ok {
		// 如果参数解析失败
		return
	}
	// 2. 处理业务逻辑
	var exist bool
	var err error
	if exist, err = logic.IsPhoneExist(p.Phone); err != nil {
		zap.L().Error("logic.IsPhoneExist failed.. ", zap.Error(err))
		return
	}

	// 3. 返回响应
	JSON(ctx, gin.H{
		"exist": exist,
	})
}

// IsEmailExist: 验证邮箱是否存在
func IsEmailExist(ctx *gin.Context) {
	//1. 验证参数， 验证手机号是否正确
	p := new(models.ParamEmailExist)
	if ok := Validate(ctx, p, ValidateSignupEmailExist); !ok {
		// 如果参数解析失败
		return
	}
	// 2. 处理业务逻辑
	var exist bool
	var err error
	if exist, err = logic.IsEmailExist(p.Email); err != nil {
		zap.L().Error("logic.IsEmailExist failed.. ", zap.Error(err))
		return
	}

	// 3. 返回响应
	JSON(ctx, gin.H{
		"exist": exist,
	})
}

// GetCaptcha: 获取图形验证码
func GetCaptcha(ctx *gin.Context) {
	// 1. 进行参数验证

	// 2. 进行业务处理：生成验证码
	id, b64s, _, err := captcha.NewCaptcha().GenerateCaptcha()
	if err != nil {
		zap.L().Error("GenerateCaptcha failed..", zap.Error(err))
	}

	JSON(ctx, gin.H{
		"captcha_id":    id,
		"captcha_image": b64s,
	})
}

// SendPhoneCode: 给制定号码发送验证码
func SendPhoneCode(ctx *gin.Context) {
	// 1. 验证参数
	p := new(models.ParamPhoneCode)
	if ok := Validate(ctx, p, ValidatePhoneCodeRequest); !ok {
		return
	}
	// 2. 进行业务处理：发送验证码， 保存验证码
	if err := logic.SendPhoneCode(p.Phone); err != nil {
		zap.L().Error("send phone code failed..", zap.Error(err))
		ResponseError(ctx, CodePhoneCodeSendError)
		return
	}
	// 3. 返回响应
	ResponseSuccess(ctx, nil)
}

// SendEmailCode: 给指定邮箱发送验证码
func SendEmailCode(ctx *gin.Context) {
	// 1. 验证参数
	p := new(models.ParamEmailCode)
	if ok := Validate(ctx, p, ValidateEmailCodeRequest); !ok {
		return
	}

	// 2. 业务逻辑：发送验证码给指定邮箱
	if err := logic.SendEmailCode(p.Email); err != nil {
		zap.L().Error("send email code failed..", zap.Error(err))
		ResponseError(ctx, CodeEmailCodeSendError)
		return
	}
	//3. 返回响应
	ResponseSuccess(ctx, nil)
}

// SignupUsingPhone: 使用手机号码进行注册
func SignupUsingPhone(ctx *gin.Context) {
	// 1. 进行参数的验证
	p := new(models.ParamSignupUsingPhone)
	if ok := Validate(ctx, p, ValidateSignupUsingPhone); !ok {
		return
	}
	// 2. 处理业务逻辑 ， 进行注册
	if err := logic.SignupUsingPhone(p); err != nil {
		zap.L().Error("logic.SignupUsingPhone failed..", zap.Error(err))
		if err == mysql.ErrorUserExist {
			ResponseError(ctx, CodeUserExist)
			return
		} else if err == mysql.ErrorPhoneExist {
			ResponseError(ctx, CodePhoneExist)
			return
		} else if err == mysql.ErrorEmailExist {
			ResponseError(ctx, CodeEmailExist)
			return
		}
		ResponseError(ctx, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(ctx, nil)
}
