package mysql

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"strings"
	"web_app/models"
)

func CreatePost(p *models.Post) (err error) {
	sqlstr := `insert into post(post_id, title, content, author_id, community_id) values (?, ?, ?, ?, ?)`
	_, err = dbSqlx.Exec(sqlstr, p.PostID, p.Title, p.Content, p.AuthorId, p.CommunityID)
	if err != nil {
		zap.L().Error("insert post failed !", zap.Error(err))
		err = ErrorInsertFailed
		return
	}
	return
}

func GetPostById(postId int64) (p *models.Post, err error) {
	p = new(models.Post)
	sqlstr := `select post_id, title, content, author_id, community_id, status, create_time from post where post_id = ?`
	err = dbSqlx.Get(p, sqlstr, postId)
	if err != nil {
		zap.L().Error("Get Post by id failed !", zap.Int64("post_id", postId), zap.Error(err))
		return
	}
	return
}

func GetPostList(page, size int64) (postList []*models.Post, err error) {
	sqlstr := `select post_id, title, content, author_id, community_id, status, create_time from post
				ORDER BY create_time
				DESC 
				limit ?,?`

	postList = make([]*models.Post, 0, 2)

	err = dbSqlx.Select(&postList, sqlstr, (page-1)*size, size)
	if err != nil {
		zap.L().Error("Get post list failed !", zap.Error(err))
		return
	}
	return
}

func GetPostLIstByIDs(ids []string) (postList []*models.Post, err error) {
	sqlstr := `select post_id, title, content, author_id, community_id, create_time from post
			where post_id in (?)
			order by FIND_IN_SET(post_id, ?)`
	query, args, err := sqlx.In(sqlstr, ids, strings.Join(ids, ",")) //组合参数
	if err != nil {
		return nil, err
	}
	query = dbSqlx.Rebind(query) //

	err = dbSqlx.Select(&postList, query, args...) // 这里的args后面一定要加上...!!!!
	return
}
