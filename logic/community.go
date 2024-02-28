package logic

import (
	"github.com/xiaorui/web_app/dao/mysql"
	"github.com/xiaorui/web_app/models"
)

func GetCommunityList() ([]*models.Community, error) {
	//其实这个业务就只有一个内容， 它不需要去判断用户是否存在， 是否登录在前面的中间件已经完成了
	return mysql.GetCommunityList()
}

func GetCommunityDetail(communityID int64) (*models.Community, error) {
	//这个函数处理的逻辑就是根据id去数据库中查询数据
	return mysql.GetCommunityDetailByID(communityID)
}
