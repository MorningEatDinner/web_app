package redis

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 // 每个投票增加的分数
)

var (
	ErrVoteTimeExpire = errors.New("超出时间限制")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

func VoteForPost(userID int64, postID int64, direction int8) (err error) {
	//1. 判断投票限制
	//判断发帖时间
	postTime := RDB.Client.ZScore(RDB.Context, getRedisKey(KeyPostTimeZSet), fmt.Sprintf("%d", postID)).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	//2 更新帖子分数
	//先查看之前投票的分数
	preDirection := RDB.Client.ZScore(RDB.Context, getRedisKey(KeyPostVotedZSetPF+fmt.Sprintf("%d", postID)),
		fmt.Sprintf("%d", userID)).Val() // 就是之前这个用户给这个post的投票记录
	//如果投票记录相同， 则不需要发起请求了
	if int8(preDirection) == direction {
		return ErrVoteRepeated
	}
	diff := math.Abs(preDirection - float64(direction))
	var dir float64
	if float64(direction) > preDirection {
		dir = 1
	} else {
		dir = -1
	}
	pipeline := RDB.Client.TxPipeline()
	pipeline.ZIncrBy(RDB.Context, getRedisKey(KeyPostScoreZSet), dir*diff*scorePerVote, fmt.Sprintf("%d", postID))
	//3. 更新用户为该帖子投票的数据
	if direction == 0 {
		// 是取消投票， 那么就要删除投票记录哦
		pipeline.ZRem(RDB.Context, getRedisKey(KeyPostVotedZSetPF+fmt.Sprintf("%d", postID)), fmt.Sprintf("%d", userID))
	} else {
		pipeline.ZAdd(RDB.Context, getRedisKey(KeyPostVotedZSetPF+fmt.Sprintf("%d", postID)), redis.Z{
			Score:  float64(direction),
			Member: fmt.Sprintf("%d", userID),
		})
	}
	_, err = pipeline.Exec(RDB.Context)

	return
}
