package controllers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"web_app/logic"
)

func CommunityHandler(c *gin.Context) {
	//1. 无（这里不需要参数校验等）
	//2. 查询到社区所有的 community_id community_name 以列表的形式返回 （业务处理-查询数据库等）
	CommunityList, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed!", zap.Error(err))
		//当发生错误的时候，不返回过于详细的信息给前端，而是把错误记录到日志中
		ResponseError(c, CodeServerBusy)
		return
	}
	//fmt.Println(CommunityList)
	//3.返回响应
	ResponseSuccess(c, CommunityList)
}

func CommunityDetailHandler(c *gin.Context) {
	//1. 参数处理-获取社区id
	communityIDStr := c.Param("id")
	communityID, err := strconv.ParseInt(communityIDStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	//2. 查询到指定的community_id的详细信息返回 （业务处理-查询数据库等）
	CommunityDetail, err := logic.GetCommunityDetail(communityID)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail(communityID) failed!", zap.Error(err))
		//当发生错误的时候，不返回过于详细的信息给前端，而是把错误记录到日志中
		ResponseError(c, CodeServerBusy)
		return
	}
	//fmt.Println(CommunityDetail)
	//3.返回响应
	ResponseSuccess(c, CommunityDetail)
}
