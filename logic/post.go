package logic

import (
	"go.uber.org/zap"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/models"
	"web_app/pkg/snowflake"
)

func CreatePost(p *models.Post) (err error) {
	//1 generate post_id
	p.PostID = snowflake.GenID()

	//2 to mysql
	err = mysql.CreatePost(p)
	if err != nil {
		return
	}

	redis.CreatePost(p)

	return
}

func GetPostById(postId int64) (p *models.ApiPostDetail, err error) {
	//query and merge the data to the struct
	//p = new(models.ApiPostDetail)
	//query post
	post, err := mysql.GetPostById(postId)
	if err != nil {
		zap.L().Error("mysql.GetPostById failed ! ", zap.Int64("postID", postId), zap.Error(err))
		return
	}
	//query user
	user, err := mysql.GetUserById(post.AuthorId)
	if err != nil {
		zap.L().Error("mysql.GetUserById ! ", zap.Int64("author_id", post.AuthorId), zap.Error(err))
		return
	}
	//query community
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID failed ! ", zap.Int64("community_id", post.CommunityID), zap.Error(err))
		return
	}

	p = &models.ApiPostDetail{
		AuthorName:       user.Username,
		Post:             post,
		CommunityDetails: community,
	}
	return
}

func GetPostList(page, size int64) (postList []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}

	postList = make([]*models.ApiPostDetail, 0, len(posts))

	for _, post := range posts {
		zap.L().Debug("post_time", zap.Int64("post_id", post.PostID), zap.Float64("post_time", float64(post.CreateTime.Unix())))
		//query user
		user, err_1 := mysql.GetUserById(post.AuthorId)
		if err_1 != nil {
			zap.L().Error("mysql.GetUserById ! ", zap.Int64("author_id", post.AuthorId), zap.Error(err))
			err = err_1
			return
		}
		//query community
		community, err_2 := mysql.GetCommunityDetailByID(post.CommunityID)
		if err_2 != nil {
			zap.L().Error("mysql.GetCommunityDetailByID failed ! ", zap.Int64("community_id", post.CommunityID), zap.Error(err))
			err = err_2
			return
		}

		postDetail := &models.ApiPostDetail{
			AuthorName:       user.Username,
			Post:             post,
			CommunityDetails: community,
		}
		postList = append(postList, postDetail)
	}
	return
}

func GetPostList2(p *models.ParamPostList) (postList []*models.ApiPostDetail, err error) {
	//logic 层主要写业务逻辑，但是由于这里业务逻辑简单，主要是查询数据库
	//这里要实现按照投票进行排序，只能先从redis中取出，否则无法获得该顺序  （好像redis的数据没有写回数据库？）

	//redis -> ids(idList)
	ids, err := redis.GetPostListIDInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostListIDInOrder success but return 0 data")
		return
	}

	//获取投票相关的信息并返回
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {

	}

	//mysql -> posts
	posts, err := mysql.GetPostLIstByIDs(ids)
	if err != nil {
		return
	}
	postList = make([]*models.ApiPostDetail, 0, len(posts))
	for idx, post := range posts {
		zap.L().Debug("post_time", zap.Int64("post_id", post.PostID), zap.Float64("post_time", float64(post.CreateTime.Unix())))
		//query user
		user, err_1 := mysql.GetUserById(post.AuthorId)
		if err_1 != nil {
			zap.L().Error("mysql.GetUserById ! ", zap.Int64("author_id", post.AuthorId), zap.Error(err))
			err = err_1
			return
		}
		//query community
		community, err_2 := mysql.GetCommunityDetailByID(post.CommunityID)
		if err_2 != nil {
			zap.L().Error("mysql.GetCommunityDetailByID failed ! ", zap.Int64("community_id", post.CommunityID), zap.Error(err))
			err = err_2
			return
		}

		postDetail := &models.ApiPostDetail{
			AuthorName:       user.Username,
			Post:             post,
			CommunityDetails: community,
			VoteNumber:       voteData[idx],
		}
		postList = append(postList, postDetail)
	}
	return
}

func GetCommunityPostList(p *models.ParamPostList) (postList []*models.ApiPostDetail, err error) {
	//redis -> ids(idList)
	ids, err := redis.GetCommunityPostListIDInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostListIDInOrder success but return 0 data")
		return
	}

	//获取投票相关的信息并返回
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {

	}

	//mysql -> posts
	posts, err := mysql.GetPostLIstByIDs(ids)
	if err != nil {
		return
	}
	postList = make([]*models.ApiPostDetail, 0, len(posts))
	for idx, post := range posts {
		zap.L().Debug("post_time", zap.Int64("post_id", post.PostID), zap.Float64("post_time", float64(post.CreateTime.Unix())))
		//query user
		user, err_1 := mysql.GetUserById(post.AuthorId)
		if err_1 != nil {
			zap.L().Error("mysql.GetUserById ! ", zap.Int64("author_id", post.AuthorId), zap.Error(err))
			err = err_1
			return
		}
		//query community
		community, err_2 := mysql.GetCommunityDetailByID(post.CommunityID)
		if err_2 != nil {
			zap.L().Error("mysql.GetCommunityDetailByID failed ! ", zap.Int64("community_id", post.CommunityID), zap.Error(err))
			err = err_2
			return
		}

		postDetail := &models.ApiPostDetail{
			AuthorName:       user.Username,
			Post:             post,
			CommunityDetails: community,
			VoteNumber:       voteData[idx],
		}
		postList = append(postList, postDetail)
	}
	return
}

func GetPostListNew(p *models.ParamPostList) (postList []*models.ApiPostDetail, err error) {
	if p.CommunityID == 0 {
		postList, err = GetPostList2(p)
	} else {
		postList, err = GetCommunityPostList(p)
	}
	if err != nil {
		zap.L().Error("get PostListNew failed", zap.Error(err))
		return
	}
	return
}
