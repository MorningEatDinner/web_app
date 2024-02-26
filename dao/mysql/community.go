package mysql

import (
	"database/sql"
	"github.com/xiaorui/web_app/models"
	"go.uber.org/zap"
)

func GetCommunityList() (communities []*models.Community, err error) {
	sqlStr := `select community_id, community_name from community`
	err = db.Select(&communities, sqlStr)
	if err == sql.ErrNoRows {
		zap.L().Error("there is no comminity in db")
		err = nil
	}
	return
}

func GetCommunityDetailByID(communityID int64) (community *models.CommunityDetail, err error) {
	community = new(models.CommunityDetail)
	sqlStr := `select community_id, community_name, introduction, create_time from community where id=?`
	err = db.Get(community, sqlStr, communityID)
	if err == sql.ErrNoRows {
		zap.L().Error("there is no this id in db")
		err = ErrorInvalidID
	}
	return
}
