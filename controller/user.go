package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/xiaorui/web_app/dao/mysql"
	"github.com/xiaorui/web_app/logic"
	"github.com/xiaorui/web_app/models"
	"go.uber.org/zap"
)

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

func LoginHandler(ctx *gin.Context) {
	//1. 进行参数校验
	p := new(models.ParamLogin)
	if err := ctx.ShouldBindJSON(p); err != nil {
		//如果发生获取参数发生了错误
		zap.L().Error("Signup with invalid param", zap.Error(err))
		err, ok := err.(validator.ValidationErrors)
		if !ok {
			//如果不是验证器错误
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		//如果是验证器错误, 就是这里捕获错误之后进行翻译返回
		//ctx.JSON(http.StatusOK, gin.H{
		//	"msg": removeTopStruct(err.Translate(trans)),
		//})
		ResponseErrorWithMsg(ctx, CodeInvalidParam, removeTopStruct(err.Translate(trans)))
	}

	//2. 业务处理
	user, err := logic.Login(p) // 登录之后获取一个token
	if err != nil {
		zap.L().Error("Login with invalid data...", zap.String("username", p.Username), zap.Error(err))
		//ctx.JSON(http.StatusOK, gin.H{
		//	"msg": "登录失败",
		//})
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(ctx, CodeUserNotExist)
			return
		} else if errors.Is(err, mysql.ErrorPasswordInvalid) {
			ResponseError(ctx, CodeInvalidPassword)
			return
		}
		ResponseError(ctx, CodeServerBusy)
		return
	}

	//3. 返回响应
	ResponseSuccess(ctx, gin.H{
		"user_id":   fmt.Sprintf("%d", user.UserID),
		"user_name": user.Username,
		"token":     user.Token,
	})
}
