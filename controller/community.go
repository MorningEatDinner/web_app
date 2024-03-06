package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xiaorui/web_app/dao/mysql"
	"github.com/xiaorui/web_app/logic"
	"github.com/xiaorui/web_app/models"
	"go.uber.org/zap"
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

// CreateNewCommunity: 创建新的社区
func CreateNewCommunity(ctx *gin.Context) {
	p := new(models.ParamCreateNewCommunity)
	if ok := Validate(ctx, p, ValidateNewCommunity); !ok {
		return
	}

	// 处理业务： 创建新的社区
	if err := logic.CreateNewCommunity(p); err != nil {
		if err == mysql.ErrorCommunityExist {
			ResponseError(ctx, CodeCommunityExist)
			return
		}
		ResponseError(ctx, CodeServerBusy)
		return
	}

	ResponseSuccess(ctx, nil)
}
