package mysql

import (
	"database/sql"

	"github.com/xiaorui/web_app/models"
	"go.uber.org/zap"
)

func GetCommunityList() (communities []*models.Community, err error) {
	// sqlStr := `select community_id, community_name from community`
	// err = db.Select(&communities, sqlStr)
	res := DB.Model(&models.Community{}).Select("community_id, community_name").Find(&communities)
	if res.Error == sql.ErrNoRows {
		zap.L().Error("there is no comminity in db")
		err = nil
	}
	return
}

func GetCommunityDetailByID(communityID int64) (community *models.Community, err error) {
	community = new(models.Community)
	// sqlStr := `select community_id, community_name, introduction, create_time from community where id=?`
	// err = db.Get(community, sqlStr, communityID)
	result := DB.Where("community_id = ?", communityID).First(community)
	if result.Error == sql.ErrNoRows {
		zap.L().Error("there is no this id in db")
		err = ErrorInvalidID
	}
	return
}

// CheckCommunityExist: 确认社区名是否存在
func CheckCommunityExist(name string) error {
	var count int64
	err := DB.Model(&models.Community{}).Where("community_name = ?", name).Count(&count).Error
	if count > 0 {
		return ErrorCommunityExist
	}
	return err
}

// InsertCommunity: 插入新的社区
func InsertCommunity(comm *models.Community) error {
	return DB.Create(comm).Error
}
