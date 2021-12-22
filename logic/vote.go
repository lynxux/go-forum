package logic

import (
	"go.uber.org/zap"
	"strconv"
	"web_app/dao/redis"
	"web_app/models"
)

// 推荐阅读
// 基于用户投票的相关算法：http://www.ruanyifeng.com/blog/algorithm/

// 本项目使用简化版的投票分数
// 投一票就加432分   86400/200  --> 200张赞成票可以给你的帖子续一天
// 分数越高则排名越靠前，初始值为帖子创建的时间戳

/* 投票的几种情况：
direction=1时，有两种情况：
	1. 之前没有投过票0，现在投赞成票1    --> 更新分数和投票记录   diff: 1  +432
	2. 之前投反对票-1，现在改投赞成票1   --> 更新分数和投票记录   diff: 2  +432*2
direction=0时，有两种情况：
	1. 之前投过赞成票1，现在要取消投票0  --> 更新分数和投票记录   diff:  -432
	2. 之前投过反对票-1，现在要取消投票0  --> 更新分数和投票记录   diff:  +432
direction=-1时，有两种情况：
	1. 之前没有投过票0，现在投反对票-1    --> 更新分数和投票记录  diff: 1  -432
	2. 之前投赞成票1，现在改投反对票-1    --> 更新分数和投票记录  diff: 2  -432*2

投票的限制：
每个贴子自发表之日起一个星期之内允许用户投票，超过一个星期就不允许再投票了。   =>post:time
	1. 到期之后将redis中保存的赞成票数及反对票数存储到mysql表中
	2. 到期之后删除那个 KeyPostVotedZSetPF
*/

func VoteForPost(userid int64, p *models.ParamVoteData) (err error) {
	//实现业务-----这里都是redis的操作
	zap.L().Debug("VoteForPost",
		zap.Int64("userID", userid),
		zap.String("postID", p.PostID),
		zap.Int8("direction", p.Direction))
	return redis.VoteforPost(strconv.Itoa(int(userid)), p.PostID, float64(p.Direction))

}
