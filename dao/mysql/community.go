package mysql

import (
	"database/sql"
	"go.uber.org/zap"
	"web_app/models"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	sqlstr := `select community_id, community_name from community`
	if err = dbSqlx.Select(&communityList, sqlstr); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
	}
	return
}

func GetCommunityDetailByID(communityID int64) (communityDetail *models.CommunityDetails, err error) {
	communityDetail = new(models.CommunityDetails)
	sqlstr := `select community_id, community_name, introduction, create_time from community where community_id = ?`
	err = dbSqlx.Get(communityDetail, sqlstr, communityID)
	if err == sql.ErrNoRows {
		err = ErrorInvalidID
		return
	}
	if err != nil {
		zap.L().Error("query communityDetail failed!", zap.String("sql ", sqlstr), zap.Error(err))
		err = ErrorQueryFailed
		return
	}
	return
}
