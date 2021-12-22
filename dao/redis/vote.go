package redis

import (
	"errors"
	"github.com/go-redis/redis/v8"
	"math"
	"time"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 // 每一票值多少分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过！")
	ErrRepeatVote     = errors.New("重复投票！")
)

//这里使用float64， 是由于redis Val()返回float64的值，这里为了方便处理
func VoteforPost(userID, postID string, directionValue float64) (err error) {
	// 1 判断投票权限
	// 获取发帖时间
	postTime := rdb.ZScore(ctx, getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}

	// 2 更新分数
	// 2.1 获取之前的投票状态 => 1: 之前投agree  -1：之前投反对  0：之前未投
	oldDirecValue := rdb.ZScore(ctx, getRedisKey(keyPostVotedZSetPF+postID), userID).Val()
	// 2.2 update
	if oldDirecValue == directionValue {
		return ErrRepeatVote
	}
	var op float64
	if directionValue > oldDirecValue {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(oldDirecValue - directionValue) // 计算两次投票差值
	pipeline := rdb.TxPipeline()
	pipeline.ZIncrBy(ctx, getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID)

	// 3 记录用户为帖子投票的数据
	if directionValue == 0 { // 表示恢复到未投票状态，直接删除bluebell：vote：userid 的数据
		pipeline.ZRem(ctx, getRedisKey(keyPostVotedZSetPF+postID), userID)
	} else { //新增已投票状态
		pipeline.ZAdd(ctx, getRedisKey(keyPostVotedZSetPF+postID), &redis.Z{
			Score:  directionValue,
			Member: userID,
		})
	}
	_, err = pipeline.Exec(ctx)
	return err
}
