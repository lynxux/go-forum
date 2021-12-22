package redis

const (
	KeyPrefix          = "bluebell:"
	KeyPostTimeZSet    = "post:time"
	KeyPostScoreZSet   = "post:score"
	keyPostVotedZSetPF = "post:voted:" //PF->prefix, 参数是post_id

	KeyCommunitySetPF = "community:"
)

func getRedisKey(key string) string {
	return KeyPrefix + key
}
