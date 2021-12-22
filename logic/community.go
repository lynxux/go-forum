package logic

import (
	"web_app/dao/mysql"
	"web_app/models"
)

func GetCommunityList() ([]*models.Community, error) {
	return mysql.GetCommunityList()
}

func GetCommunityDetail(communityID int64) (*models.CommunityDetails, error) {
	return mysql.GetCommunityDetailByID(communityID)
}
