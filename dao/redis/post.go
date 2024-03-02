package redis

import (
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/xiaorui/web_app/models"
)

func CreatePost(pid, communityID int64) error {

	//使用事务操作， redis事务
	pipeline := RDB.Client.TxPipeline()
	//在redis中加入一个创建的post的记录
	pipeline.ZAdd(RDB.Context, getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: pid, //需要确认这里是否需要是string类型的？先暂时使用int
	})
	// 在创建post的时候， 也加入了对于分数的初始化
	pipeline.ZAdd(RDB.Context, getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: pid, //需要确认这里是否需要是string类型的？先暂时使用int
	})

	communityKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityID)))
	pipeline.SAdd(RDB.Context, communityKey, pid) // 加入member， 但是不需要score

	_, err := pipeline.Exec(RDB.Context)
	return err
}

func getIDSFromKey(key string, page, size int64) ([]string, error) {
	//得到按照某种方式进行排序
	start := (page - 1) * size
	end := page * size
	return RDB.Client.ZRevRange(RDB.Context, key, start, end).Result()
}

func GetPostIDListByOrder(p *models.ParamPostList) ([]string, error) {
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}

	return getIDSFromKey(key, p.Page, p.Size)
}

func GetVotesByPostIDS(pidList []string) ([]int64, error) {
	pipeline := RDB.Client.Pipeline()
	for _, id := range pidList {
		key := getRedisKey(KeyPostVotedZSetPF + id)
		pipeline.ZCount(RDB.Context, key, "1", "1")
	}
	cmders, err := pipeline.Exec(RDB.Context)
	if err != nil {
		return nil, err
	}
	data := make([]int64, 0, len(pidList))
	for _, cmder := range cmders {
		value := cmder.(*redis.IntCmd).Val()
		data = append(data, value)
	}
	return data, nil
}

func GetCommunityPostIDListByOrder(p *models.ParamPostList) ([]string, error) {

	// 再多加上一个key， 如果一段时间内重复查询会更快， 也就是加上一个对之前查询结果的缓存
	communityKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityID)))
	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	// 如果不存在， 也就是如果缓存里面没有， 那么就需要查询了
	if RDB.Client.Exists(RDB.Context, key).Val() < 1 { //
		// 需要计算
		pipeline := RDB.Client.Pipeline()
		pipeline.ZInterStore(RDB.Context, key, &redis.ZStore{
			Aggregate: "MAX", // 这里的意思是相同元素的聚合方式
			Keys:      []string{communityKey, orderKey},
		}) // 注意， 值最终是保存到一个zset中的
		pipeline.Expire(RDB.Context, key, time.Second*60)
		_, err := pipeline.Exec(RDB.Context)
		if err != nil {
			return nil, err
		}
	}

	// 上面结束之后就得到key中的值就是获取了key的对应的值， 就是说所有的id
	return getIDSFromKey(key, p.Page, p.Size)
}
