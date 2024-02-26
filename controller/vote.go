package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/xiaorui/web_app/logic"
	"github.com/xiaorui/web_app/models"
	"go.uber.org/zap"
)

func PostVoteHandler(ctx *gin.Context) {
	//1. 获取参数进行参数校验
	p := new(models.ParamVoteData)
	if err := ctx.ShouldBindJSON(p); err != nil {
		//这里还希望将错误进行转换为中文
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		errData := removeTopStruct(errs.Translate(trans)) // 如果存在错误转为中文进行输出
		ResponseErrorWithMsg(ctx, CodeInvalidParam, errData)
		return
	}
	userID, err := getCurrentUser(ctx)
	if err != nil {
		ResponseError(ctx, CodeNeedLogin)
		return
	}

	//2. 处理业务逻辑

	if err := logic.VoteForPost(userID, p); err != nil {
		zap.L().Error("PostVoteHandler logic.VoteForPost failed.", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}

	//3. 返回响应
	ResponseSuccess(ctx, nil)
}
