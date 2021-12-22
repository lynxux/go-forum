package redis

import (
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
	"web_app/models"
)

func CreatePost(p *models.Post) (err error) {
	// 这里要求两个操作都成功，即要求事务
	pipeline := rdb.TxPipeline()

	//time
	pipeline.ZAdd(ctx, getRedisKey(KeyPostTimeZSet), &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: p.PostID,
	})

	//score
	pipeline.ZAdd(ctx, getRedisKey(KeyPostScoreZSet), &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: p.PostID,
	})

	//community
	pipeline.SAdd(ctx, getRedisKey(KeyCommunitySetPF+strconv.Itoa(int(p.CommunityID))), p.PostID)

	_, err = pipeline.Exec(ctx)
	return
}

func getIDsFormKey(key string, page, size int64) ([]string, error) {
	start := (page - 1) * size
	end := start + size - 1
	return rdb.ZRevRange(ctx, key, start, end).Result()
}

func GetPostListIDInOrder(p *models.ParamPostList) ([]string, error) {
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	return getIDsFormKey(key, p.Page, p.Size)
}

// GetPostVoteData 根据ids查询每篇帖子的投票数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	//data = make([]int64, 0, len(ids))
	//for _, id := range ids {
	//	key := getRedisKey(keyPostVotedZSetPF + id)
	//	//这里查询了分数是1的元素的数据量 -》 即赞成票的数量
	//	v1 := rdb.ZCount(ctx, key, "1", "1").Val()
	//	data = append(data, v1)
	//}

	//
	//减少请求redis的次数
	pipeline := rdb.Pipeline()
	for _, id := range ids {
		key := getRedisKey(keyPostVotedZSetPF + id)
		pipeline.ZCount(ctx, key, "1", "1")
	}
	cmders, err := pipeline.Exec(ctx)
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

// GetCommunityPostListIDInOrder 按社区查询ids
func GetCommunityPostListIDInOrder(p *models.ParamPostList) ([]string, error) {
	//使用ZInterStore 把分区信息的set与分数or时间的zset生成一个新zset

	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}

	// 社区的key
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityID)))
	// 利用缓存key减少zinterstore执行的次数
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	if rdb.Exists(ctx, key).Val() < 1 { //这里的key为合并后新建的->bluebell:post:time/score+{communityID}, 不是ckey->bluebell:community:{communityID}
		// 不存在，需要计算
		pipeline := rdb.Pipeline()
		pipeline.ZInterStore(ctx, key, &redis.ZStore{
			Keys:      []string{cKey, orderKey},
			Aggregate: "MAX",
		}) // zinterstore 计算
		pipeline.Expire(ctx, key, 60*time.Second) // 设置超时时间
		_, err := pipeline.Exec(ctx)
		if err != nil {
			return nil, err
		}
	}

	return getIDsFormKey(key, p.Page, p.Size)
}
