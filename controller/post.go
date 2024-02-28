package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xiaorui/web_app/logic"
	"github.com/xiaorui/web_app/models"
	"go.uber.org/zap"
)

func CreatePostHandler(ctx *gin.Context) {
	// 1. 进行参数校验
	p := new(models.Post)
	err := ctx.ShouldBindJSON(p)
	if err != nil {
		zap.L().Error("CreatePostHandler ctx.ShouldBindJSON failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	//这里还需要获取当前发帖子的用户id
	userID, err := getCurrentUser(ctx)
	if err != nil {
		ResponseError(ctx, CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	// 2. 进行业务处理， 也就是说创建一个post
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost failed", zap.Error(err))
		ResponseError(ctx, CodeServerBusy) // 不要将太多的后端错误暴露给前端
	}

	// 3. 返回数据
	ResponseSuccess(ctx, nil)
}

func GetPostHandler(ctx *gin.Context) {
	//1. 进行参数校验
	id := ctx.Param("id") // 记住， 这里返回的值都是string类型的
	pid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		zap.L().Error("GetPostHandler strconv.ParseInt(pid) failed.", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	//2. 处理业务逻辑， 也就是从数据库中获取数据
	data, err := logic.GetPostByID(pid)
	if err != nil {
		zap.L().Error("GetPostHandler  logic.GetPostByID failed.", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	//3. 返回数据
	ResponseSuccess(ctx, data)
}

func GetPostListHandler(ctx *gin.Context) {
	//1. 进行参数的获取和校验
	//Get是从ctx中拿到的， 不是从地址中拿到的， Query才是从地址中的？=拿到的
	pageNum, pageSize := getPageInfo(ctx)

	//2. 处理业务的逻辑
	data, err := logic.GetPostList(pageNum, pageSize)
	if err != nil {
		zap.L().Error("GetPostListHandler logic.GetPostList failed.", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}

	//3. 返回数据
	ResponseSuccess(ctx, data)
}

// 这个函数能够实现更加高级的功能，即根据分数或者时间对post进行排序
// 1. 获取参数进行参数校验
// 2. 实现业务逻辑， 即从redis中根据参数获取postid
// 3. 从数据库中根据postid获取post的详细数据

// 现在希望加入的是统计每个帖子已经投票的分数
func GetPostListHandler2(ctx *gin.Context) {
	//1. 进行参数的获取和校验
	//Get是从ctx中拿到的， 不是从地址中拿到的， Query才是从地址中的？=拿到的
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}
	if err := ctx.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 ctx.ShouldBindQuery failed.", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}

	//2. 处理业务的逻辑
	data, err := logic.GetPostList2(p)
	if err != nil {
		zap.L().Error("GetPostListHandler logic.GetPostList failed.", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}

	//3. 返回数据
	ResponseSuccess(ctx, data)
}

// 就是说现在想要实现的功能是多一个按照community分类的post list， 也就是说相比上面那个逻辑， 现在多加上了community， 还是要保留order的
func GetCommunityPostListHandler(ctx *gin.Context) {
	// 1. 进行参数校验
	p := &models.ParamPostList{
		Page:        1,
		Size:        10,
		Order:       models.OrderTime,
		CommunityID: 1,
	}
	if err := ctx.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetCommunityPostListHandler ctx.ShouldBindQuery failed.", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	// 2. 处理业务逻辑
	data, err := logic.GetCommunityPostList(p)
	if err != nil {
		zap.L().Error("GetCommunityPostListHandler logic.GetCommunityPostList failed.", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(ctx, data)
}

// TODO:: 现在想要完成的是对上面两个接口进行融合
func GetPostListHandler0(ctx *gin.Context) {
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}

	if err := ctx.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler0 ctx.ShouldBindQuery failed.", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	// zap.L().Info("param", zap.Any("param", p))

	//处理业务逻辑
	data, err := logic.GetPostList0(p)
	if err != nil {
		zap.L().Error("GetPostListHandler0 logic.GetCommunityPostList failed.", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(ctx, data)
}
