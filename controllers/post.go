package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"strconv"
	"web_app/logic"
	"web_app/models"
)

func CreatePostHandler(c *gin.Context) {
	//1.参数处理
	//完整性校验
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		zap.L().Error("create post with invalid param !", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	//校验能否获取author_id，即当前的user_id
	userId, err := GetCurrentUserID(c)
	if err != nil {
		zap.L().Error("logic.GetCurrentUserID failed !", zap.Error(err))
		ResponseError(c, CodeNeedLogin)
		return
	}
	post.AuthorId = userId

	//2.创建帖子---插入数据到数据库
	if err := logic.CreatePost(&post); err != nil {
		zap.L().Error("logic.CreatePost() failed !", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	//3.返回响应
	ResponseSuccess(c, nil)
}

func PostDetailHandler(c *gin.Context) {
	//1. param --Post_id
	postIdStr := c.Param("id")
	postId, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail-id failed! ", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	//2 query post from post-table
	post, err := logic.GetPostById(postId)
	if err != nil {
		zap.L().Error("logic.GetPost() failed !", zap.String("post_id", postIdStr), zap.Error(err))
		ResponseError(c, CodeServerBusy)
	}

	//3
	ResponseSuccess(c, post)
}

func PostListHandler(c *gin.Context) {
	//1. param -> page,size
	page, size := GetPageInfo(c)
	//2.
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("get post list failed !", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3
	ResponseSuccess(c, data)
}

//根据前端传来的参数 动态的获取帖子的列表
//可以按照帖子的 创建时间time 以及 获得的投票score 进行排序
func PostListHandler2(c *gin.Context) {
	//1 param => api/v1/posts2?page=1?size=2&order=time
	// 除了使用c.Query 还可以使用 c.shouldBindQuery
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("Get PostListHandler2 failed with invalid param ", zap.Error(err))
		return
	}

	//2.data =>  redis -> mysql
	data, err := logic.GetPostListNew(p)

	if err != nil {
		zap.L().Error("get post list in order failed !", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3 response
	ResponseSuccess(c, data)
}

//func GetCommunityPostListHandler(c *gin.Context) {
//	//1 param => api/v1/posts2?page=1?size=2&order=time
//	// 除了使用c.Query 还可以使用 c.shouldBindQuery
//	p := &models.ParamPostList{
//		CommunityID: 1,
//		Page:        1,
//		Size:        10,
//		Order:       models.OrderTime,
//	}
//	if err := c.ShouldBindQuery(p); err != nil {
//		zap.L().Error("Get GetCommunityPostListHandler failed with invalid param ", zap.Error(err))
//		return
//	}
//
//	//2.data =>  redis -> mysql
//	data, err := logic.GetCommunityPostList(p)
//	if err != nil {
//		zap.L().Error("get community post list in order failed !", zap.Error(err))
//		ResponseError(c, CodeServerBusy)
//		return
//	}
//	//3 response
//	ResponseSuccess(c, data)
//
//}

func VotePostHandler(c *gin.Context) {
	//1. param ->
	p := new(models.ParamVoteData) //return the point of object
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors) //类型断言
		if !ok {                                     //返回的不是validator类型的错误，说明不是参数bing的问题？
			zap.L().Error("vote post with invalid param !", zap.Error(err))
			ResponseError(c, CodeInvalidParam)
		}
		//否则是参数banding的问题？ 还是说json也算作validator里？
		errData := errs.Translate(trans)
		ResponseErrorWithMsg(c, CodeInvalidParam, errData)
		return
	}

	//2
	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	if err := logic.VoteForPost(userID, p); err != nil {
		zap.L().Error("logic.VoteForPost() failed! ", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	//3
	ResponseSuccess(c, nil)
}
