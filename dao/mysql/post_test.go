package mysql

import (
	"testing"
	"web_app/models"
	"web_app/settings"
)

func init() {
	dbCfg := settings.MySqlConfig{
		Host:         "127.0.0.1",
		User:         "root",
		Password:     "123456",
		DbName:       "bluebell",
		Port:         3306,
		MaxOpenConns: 10,
		MaxIdleConns: 10,
	}
	err := Init(&dbCfg)
	if err != nil {
		panic(err)
	}
}

func TestCreatePost(t *testing.T) {
	post := models.Post{
		PostID:      10,
		Title:       "test",
		Content:     "just for test",
		AuthorId:    123,
		CommunityID: 1,
	}

	err := CreatePost(&post)
	if err != nil {
		t.Fatalf("CreatePost (insert record to mysql) failed, err:%v\n", err)
	}
	t.Logf("CreatePost (insert record to mysql) success !")
}
