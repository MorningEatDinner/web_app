package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaorui/web_app/dao/mysql"
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
