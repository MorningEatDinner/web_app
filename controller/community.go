package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaorui/web_app/logic"
	"go.uber.org/zap"
	"strconv"
)

func CommunityHandler(ctx *gin.Context) {
	// 查询到所有的数据(id, community_name)
	data, err := logic.GetCommunityList() // 从表中获取数据
	if err != nil {
		zap.L().Error("logic.GetCommunityList failed.", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	ResponseSuccess(ctx, data)
}

func CommunityDetailHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	communityID, err := strconv.ParseInt(idStr, 10, 64) // 10进制， int64类型
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	data, err := logic.GetCommunityDetail(communityID)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	ResponseSuccess(ctx, data)
}
