package redis

//存放redis key

const (
	KeyPrefix          = "bluebell:"
	KeyPostTimeZSet    = "post:time"
	KeyPostScoreZSet   = "post:score"
	KeyPostVotedZSetPF = "post:voted:" // 这里更改了好像会有点麻烦
	KeyCommunitySetPF  = "community:"  // 保存每个community下面的post的集合
)

// 给key加上前缀
func getRedisKey(key string) string {
	return KeyPrefix + key
}
