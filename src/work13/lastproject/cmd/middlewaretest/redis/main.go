package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"github.com/go-redis/redis/v8"
	"math/big"
	"strings"
	"time"
)

var (
	err           error
	clusterClient *redis.ClusterClient
	redisTTL      = 604800 * time.Second
	ctx           = context.Background()
)

type connOption struct {
	command      string
	redisKey     string
	redisVal     interface{}
	redisExpire  time.Duration
	startIndex   int64
	endIndex     int64
	lInsertType  string
	lInsertVal   string
	hashFields   string
	memberFields string
	zMin, zMax   string
}

// 连接redis集群
func init() {
	//todo ip修改
	clusterClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			"localhost:6379",
			"localhost:6380",
			"localhost:6381",
			"localhost:6382",
			"localhost:6383",
			"localhost:6384",
		},
		Password:     "123456ssw",
		DialTimeout:  50 * time.Second, // 设置连接超时
		ReadTimeout:  50 * time.Second, // 设置读取超时
		WriteTimeout: 50 * time.Second, // 设置写入超时
	})

	// 发送一个ping命令,测试是否通
	s := clusterClient.Do(ctx, "ping").String()
	fmt.Println(s)
	fmt.Println()
}

func main() {
	//basicOperation()
	//listOperation()
	//setOperation()
	//hashOperation()
	//sortedSetOperation()
	pipelineOperation()
	watchOperation()
}

//redis pipeline操作
func pipelineOperation() {
	zsetKey := "pipline_test_name"
	// 开启一个TxPipeline事务
	pipe := clusterClient.TxPipeline()

	// 执行事务操作，可以通过pipe读写redis
	incr := pipe.Incr(ctx, zsetKey)
	pipe.Expire(ctx, zsetKey, redisTTL)

	// 通过Exec函数提交redis事务
	_, err = pipe.Exec(ctx)
	PrintPanic(err)

	// 提交事务后，我们可以查询事务操作的结果
	// 前面执行Incr函数，在没有执行exec函数之前，实际上还没开始运行。
	fmt.Println(incr.Val(), err)
}

//redis watch操作
func watchOperation() {
	keyName := "watch_test_name"
	// 定义一个回调函数，用于处理事务逻辑
	fn := func(tx *redis.Tx) error {
		// 先查询下当前watch监听的key的值
		v, err := tx.Get(ctx, keyName).Result()
		if err != nil && err != redis.Nil {
			return err
		}

		fmt.Println(v)

		// 如果key的值没有改变的话，Pipelined函数才会调用成功
		_, err = tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			// 在这里给key设置最新值
			pipe.Set(ctx, keyName, "new value2", redisTTL)
			return nil
		})
		return err
	}

	// 使用Watch监听一些Key, 同时绑定一个回调函数fn, 监听Key后的逻辑写在fn这个回调函数里面
	// 如果想监听多个key，可以这么写：client.Watch(fn, "key1", "key2", "key3")
	clusterClient.Watch(ctx, fn, keyName)
}

//redis sorted set操作
func sortedSetOperation() {
	// ZAdd ZIncrBy zrange
	fmt.Println("---------------------redis  ZAdd ZIncrBy zrange 使用---------------------")
	HandlerRedisSortedSetCommand(connOption{command: "zadd", redisKey: "zadd_test_name", redisVal: &redis.Z{Score: 90.0, Member: "Golang"}})
	HandlerRedisSortedSetCommand(connOption{command: "zadd", redisKey: "zadd_test_name", redisVal: &redis.Z{Score: 98.0, Member: "Java"}})
	HandlerRedisSortedSetCommand(connOption{command: "zadd", redisKey: "zadd_test_name", redisVal: &redis.Z{Score: 95.0, Member: "Python"}})
	HandlerRedisSortedSetCommand(connOption{command: "zadd", redisKey: "zadd_test_name", redisVal: &redis.Z{Score: 97.0, Member: "JavaScript"}})
	HandlerRedisSortedSetCommand(connOption{command: "zadd", redisKey: "zadd_test_name", redisVal: &redis.Z{Score: 92.0, Member: "C/C++"}})
	HandlerRedisSortedSetCommand(connOption{command: "zincrby", redisKey: "zadd_test_name", memberFields: "Vue", redisVal: float64(8)})
	HandlerRedisSortedSetCommand(connOption{command: "zrange", redisKey: "zadd_test_name", startIndex: 0, endIndex: -1})
	fmt.Println()

	// ZCount ZCard ZRangeByScore
	fmt.Println("---------------------redis  ZCount  ZCard  ZRangeByScore 使用---------------------")
	HandlerRedisSortedSetCommand(connOption{command: "zcard", redisKey: "zadd_test_name"})
	HandlerRedisSortedSetCommand(connOption{command: "zcount", redisKey: "zadd_test_name", zMin: "95", zMax: "100"})
	HandlerRedisSortedSetCommand(connOption{command: "zrangebyscore", redisKey: "zadd_test_name", redisVal: &redis.ZRangeBy{
		Min:    "80",  // 最小分数
		Max:    "100", // 最大分数
		Offset: 0,     // 类似sql的limit, 表示开始偏移量
		Count:  5,     // 一次返回多少数据
	}})
	fmt.Println()

	// ZRank & ZScore
	fmt.Println("---------------------redis  ZRank & ZScore 使用---------------------")
	HandlerRedisSortedSetCommand(connOption{command: "zscore", redisKey: "zadd_test_name", memberFields: "Golang"})
	HandlerRedisSortedSetCommand(connOption{command: "zrank", redisKey: "zadd_test_name", memberFields: "Java"})
	fmt.Println()

	// ZRem & ZRemRangeByRank
	fmt.Println("---------------------redis  ZRem & ZRemRangeByRank 使用---------------------")
	HandlerRedisSortedSetCommand(connOption{command: "zrem", redisKey: "zadd_test_name", memberFields: "Java"})
	HandlerRedisSortedSetCommand(connOption{command: "zremrangebyrank", redisKey: "zadd_test_name", startIndex: 0, endIndex: 3})
	HandlerRedisSortedSetCommand(connOption{command: "zrange", redisKey: "zadd_test_name", startIndex: 0, endIndex: -1})
	fmt.Println()
}

//处理redis sorted set操作
func HandlerRedisSortedSetCommand(opt connOption) {
	switch opt.command {
	case "zadd":
		val, err := clusterClient.ZAdd(ctx, opt.redisKey, (opt.redisVal).(*redis.Z)).Result()
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "ZAdd添加成功元素：", val)
	case "zaddnx":
	case "zaddxx":
	case "zaddch":
	case "zaddnxch":
	case "zincr":
	case "zincrnx":
	case "zincrxx":
	case "zincrby":
		clusterClient.ZIncrBy(ctx, opt.redisKey, opt.redisVal.(float64), opt.memberFields)
	case "zinterstore":
	case "zcard":
		val, err := clusterClient.ZCard(ctx, opt.redisKey).Result()
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "查询出来的值为：", val)
	case "zcount":
		val, err := clusterClient.ZCount(ctx, opt.redisKey, opt.zMin, opt.zMax).Result()
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "查询出来的值为：", val)
	case "zrange":
		val, err := clusterClient.ZRange(ctx, opt.redisKey, opt.startIndex, opt.endIndex).Result()
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "查询出来的值为：", val)
	case "zrevrange":
	case "zrangebyscore":
		val, err := clusterClient.ZRangeByScore(ctx, opt.redisKey, (opt.redisVal).(*redis.ZRangeBy)).Result()
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "查询出来的值为：", val)
	case "zremrangebyscore":
	case "zrangewithscores":
	case "zrank":
		//根据元素名，查询集合元素在集合中的排名，从0开始算，集合元素按分数从小到大排序
		val, _ := clusterClient.ZRank(ctx, opt.redisKey, opt.memberFields).Result()
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "，字段名为：", opt.memberFields, "查询出来的值为：", val)
	case "zrevrank":
	case "zscore":
		// 查询集合元素Golang的分数
		val, err := clusterClient.ZScore(ctx, opt.redisKey, opt.memberFields).Result()
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "，字段名为：", opt.memberFields, "查询出来的值为：", val)
	case "zrem":
		clusterClient.ZRem(ctx, opt.redisKey, opt.memberFields)
	case "zremrangebyrank":
		clusterClient.ZRemRangeByRank(ctx, opt.redisKey, opt.startIndex, opt.endIndex)
	}
}

//redis set操作
func hashOperation() {
	// HGet HGetAll HMSet HMGet HSet HSetNX
	fmt.Println("---------------------redis  HGet HGetAll HMSet HMGet HSet HSetNX 使用---------------------")
	HandlerRedisHashCommand(connOption{command: "hset", redisKey: "hashset_test_name", hashFields: "username", redisVal: "admin"})
	HandlerRedisHashCommand(connOption{command: "hget", redisKey: "hashset_test_name", hashFields: "username"})
	HandlerRedisHashCommand(connOption{command: "hset", redisKey: "hashset_test_name", hashFields: "password", redisVal: "abc123"})
	HandlerRedisHashCommand(connOption{command: "hgetall", redisKey: "hashset_test_name"})
	batchData := make(map[string]interface{})
	batchData["username"] = "test"
	batchData["password"] = 123456
	HandlerRedisHashCommand(connOption{command: "hmset", redisKey: "hashmset_test_name", redisVal: batchData})
	HandlerRedisHashCommand(connOption{command: "hsetnx", redisKey: "hashmset_test_name", hashFields: "email", redisVal: "ourlang@foxmail.com"})
	HandlerRedisHashCommand(connOption{command: "hmget", redisKey: "hashmset_test_name", hashFields: "username,password,email"})
	fmt.Println()

	// HIncrBy & HIncrByFloat
	fmt.Println("---------------------redis  HIncrBy & HIncrByFloat 使用---------------------")
	HandlerRedisHashCommand(connOption{command: "hincrby", redisKey: "hashset_test_name", hashFields: "count", redisVal: int64(2)})
	HandlerRedisHashCommand(connOption{command: "hincrbyfloat", redisKey: "hashset_test_name", hashFields: "score", redisVal: float64(3.2)})
	HandlerRedisHashCommand(connOption{command: "hgetall", redisKey: "hashset_test_name"})
	fmt.Println()

	// HKeys & HLen
	fmt.Println("---------------------redis  HKeys & HLen 使用---------------------")
	HandlerRedisHashCommand(connOption{command: "hkeys", redisKey: "hashset_test_name"})
	HandlerRedisHashCommand(connOption{command: "hlen", redisKey: "hashset_test_name"})
	fmt.Println()

	// HDel & HExists
	fmt.Println("---------------------redis  HDel & HExists 使用---------------------")
	HandlerRedisHashCommand(connOption{command: "hdel", redisKey: "hashset_test_name", hashFields: "score"})
	HandlerRedisHashCommand(connOption{command: "hgetall", redisKey: "hashset_test_name"})
	HandlerRedisHashCommand(connOption{command: "hexists", redisKey: "hashset_test_name", hashFields: "id"})
	HandlerRedisHashCommand(connOption{command: "hexists", redisKey: "hashset_test_name", hashFields: "count"})
	fmt.Println()
}

//处理redis set操作
func HandlerRedisHashCommand(opt connOption) {
	switch opt.command {
	case "hdel":
		clusterClient.HDel(ctx, opt.redisKey, opt.hashFields)
	case "hexists":
		exists, err := clusterClient.HExists(ctx, opt.redisKey, opt.hashFields).Result()
		PrintPanic(err)
		if exists {
			fmt.Println("redis的key为：", opt.redisKey, "，字段名为：", opt.hashFields, "存在：", exists)
		} else {
			fmt.Println("redis的key为：", opt.redisKey, "，字段名为：", opt.hashFields, "不存在：", exists)
		}
	case "hget":
		val, err := clusterClient.HGet(ctx, opt.redisKey, opt.hashFields).Result()
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "，字段名为：", opt.hashFields, "查询出来的值为：", val)
	case "hgetall":
		val, err := clusterClient.HGetAll(ctx, opt.redisKey).Result()
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "查询出来的值为：", val)
	case "hincrby":
		val, err := clusterClient.HIncrBy(ctx, opt.redisKey, opt.hashFields, opt.redisVal.(int64)).Result()
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "，字段名为：", opt.hashFields, "查询出来的值为：", val)
	case "hincrbyfloat":
		val, err := clusterClient.HIncrByFloat(ctx, opt.redisKey, opt.hashFields, opt.redisVal.(float64)).Result()
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "，字段名为：", opt.hashFields, "查询出来的值为：", val)
	case "hkeys":
		val, err := clusterClient.HKeys(ctx, opt.redisKey).Result()
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "查询出来的值为：", val)
	case "hlen":
		val, err := clusterClient.HLen(ctx, opt.redisKey).Result()
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "查询出来的值为：", val)
	case "hmget":
		aKeys := strings.Split(opt.hashFields, ",")
		val, err := clusterClient.HMGet(ctx, opt.redisKey, aKeys[0], aKeys[1], aKeys[2]).Result()
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "，字段名为：", opt.hashFields, "查询出来的值为：", val)
	case "hmset":
		err = clusterClient.HMSet(ctx, opt.redisKey, opt.redisVal).Err()
		PrintPanic(err)
	case "hset":
		err = clusterClient.HSet(ctx, opt.redisKey, opt.hashFields, opt.redisVal).Err()
		PrintPanic(err)
	case "hsetnx":
		//如果字段不存在，则设置hash字段值
		clusterClient.HSetNX(ctx, opt.redisKey, opt.hashFields, opt.redisVal)
	}
}

//redis set操作
func setOperation() {
	// SAdd SCard SIsMember SMembers
	fmt.Println("---------------------redis  SAdd SCard SIsMember SMembers 使用---------------------")
	HandlerRedisSetCommand(connOption{command: "sadd", redisKey: "setadd_test_name", redisVal: 100})
	HandlerRedisSetCommand(connOption{command: "sadd", redisKey: "setadd_test_name", redisVal: 200})
	HandlerRedisSetCommand(connOption{command: "scard", redisKey: "setadd_test_name"})
	HandlerRedisSetCommand(connOption{command: "smembers", redisKey: "setadd_test_name"})
	HandlerRedisSetCommand(connOption{command: "sismember", redisKey: "setadd_test_name", redisVal: 100})
	HandlerRedisSetCommand(connOption{command: "sismember", redisKey: "setadd_test_name", redisVal: 300})
	fmt.Println()

	// SDiff SDiffStore SInter SInterStore SUnion SUnionStore
	fmt.Println("---------------------redis  SAdd SCard SIsMember SMembers 使用---------------------")
	HandlerRedisSetCommand(connOption{command: "sadd", redisKey: "{blacklist_test_name_}1", redisVal: "Obama"})
	HandlerRedisSetCommand(connOption{command: "sadd", redisKey: "{blacklist_test_name_}1", redisVal: "Hillary"})
	HandlerRedisSetCommand(connOption{command: "sadd", redisKey: "{blacklist_test_name_}1", redisVal: "the Elder"})
	HandlerRedisSetCommand(connOption{command: "sadd", redisKey: "{blacklist_test_name_}2", redisVal: "the Elder"})
	HandlerRedisSetCommand(connOption{command: "sinter", redisKey: "{blacklist_test_name_}1,{blacklist_test_name_}2"})
	HandlerRedisSetCommand(connOption{command: "sinterstore", redisKey: "{blacklist_test_name_}3,{blacklist_test_name_}1,{blacklist_test_name_}2"})
	HandlerRedisSetCommand(connOption{command: "sdiff", redisKey: "{blacklist_test_name_}1,{blacklist_test_name_}2"})
	HandlerRedisSetCommand(connOption{command: "sdiffstore", redisKey: "{blacklist_test_name_}3,{blacklist_test_name_}1,{blacklist_test_name_}2"})
	HandlerRedisSetCommand(connOption{command: "sunion", redisKey: "{blacklist_test_name_}1,{blacklist_test_name_}2"})
	HandlerRedisSetCommand(connOption{command: "sunionstore", redisKey: "{blacklist_test_name_}3,{blacklist_test_name_}1,{blacklist_test_name_}2"})
	fmt.Println()

	//SPop SPopN SRem
	fmt.Println("---------------------redis  SPop SPopN SRem 使用---------------------")
	HandlerRedisSetCommand(connOption{command: "sadd", redisKey: "setadd_test_name", redisVal: 300})
	HandlerRedisSetCommand(connOption{command: "sadd", redisKey: "setadd_test_name", redisVal: 400})
	HandlerRedisSetCommand(connOption{command: "sadd", redisKey: "setadd_test_name", redisVal: 500})
	HandlerRedisSetCommand(connOption{command: "sadd", redisKey: "setadd_test_name", redisVal: 600})
	HandlerRedisSetCommand(connOption{command: "srem", redisKey: "setadd_test_name", redisVal: 600})
	HandlerRedisSetCommand(connOption{command: "spop", redisKey: "setadd_test_name"})
	HandlerRedisSetCommand(connOption{command: "spopn", redisKey: "setadd_test_name", startIndex: 2})
	HandlerRedisSetCommand(connOption{command: "smembers", redisKey: "setadd_test_name"})
	fmt.Println()

	//SRandMember SRandMemberN
	fmt.Println("---------------------redis  SRandMember SRandMemberN 使用---------------------")
	HandlerRedisSetCommand(connOption{command: "sadd", redisKey: "setrandmember_test_name", redisVal: 300})
	HandlerRedisSetCommand(connOption{command: "sadd", redisKey: "setrandmember_test_name", redisVal: 400})
	HandlerRedisSetCommand(connOption{command: "sadd", redisKey: "setrandmember_test_name", redisVal: 500})
	HandlerRedisSetCommand(connOption{command: "sadd", redisKey: "setrandmember_test_name", redisVal: 600})
	HandlerRedisSetCommand(connOption{command: "sadd", redisKey: "setrandmember_test_name", redisVal: 200})
	HandlerRedisSetCommand(connOption{command: "smembers", redisKey: "setrandmember_test_name"})
	HandlerRedisSetCommand(connOption{command: "srandermenmber", redisKey: "setrandmember_test_name"})
	HandlerRedisSetCommand(connOption{command: "srandermenmbern", redisKey: "setrandmember_test_name", startIndex: 2})
	HandlerRedisSetCommand(connOption{command: "smembers", redisKey: "setrandmember_test_name"})
	fmt.Println()

	//SMembersMap & SMove
	fmt.Println("---------------------redis  SMembersMap & SMove 使用---------------------")
	HandlerRedisSetCommand(connOption{command: "sadd", redisKey: "{setmap_test_name_}1", redisVal: 300})
	HandlerRedisSetCommand(connOption{command: "sadd", redisKey: "{setmap_test_name_}1", redisVal: 400})
	HandlerRedisSetCommand(connOption{command: "sadd", redisKey: "{setmap_test_name_}1", redisVal: 500})
	HandlerRedisSetCommand(connOption{command: "sadd", redisKey: "{setmap_test_name_}1", redisVal: 600})
	HandlerRedisSetCommand(connOption{command: "sadd", redisKey: "{setmap_test_name_}2", redisVal: 200})
	HandlerRedisSetCommand(connOption{command: "smenmbermap", redisKey: "{setmap_test_name_}1"})
	HandlerRedisSetCommand(connOption{command: "smove", redisKey: "{setmap_test_name_}1,{setmap_test_name_}2", redisVal: 300})
	HandlerRedisSetCommand(connOption{command: "smembers", redisKey: "{setmap_test_name_}1"})
	HandlerRedisSetCommand(connOption{command: "smembers", redisKey: "{setmap_test_name_}2"})
	fmt.Println()
}

//处理redis set操作
func HandlerRedisSetCommand(opt connOption) {
	switch opt.command {
	case "sadd":
		err = clusterClient.SAdd(ctx, opt.redisKey, opt.redisVal).Err()
		PrintPanic(err)
	case "scard":
		val, err := clusterClient.SCard(ctx, opt.redisKey).Result()
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "，查询出来的值为：", val)
	case "sismember":
		flag, err := clusterClient.SIsMember(ctx, opt.redisKey, opt.redisVal).Result()
		PrintPanic(err)
		if flag {
			fmt.Println("集合:", opt.redisKey, "中包含指定元素:", opt.redisVal)
		} else {
			fmt.Println("集合:", opt.redisKey, "中不包含指定元素:", opt.redisVal)
		}
	case "smembers":
		val, err := clusterClient.SMembers(ctx, opt.redisKey).Result()
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "，查询出来的值为：", val)
	case "sdiff":
		aKeys := strings.Split(opt.redisKey, ",")
		val, err := clusterClient.SDiff(ctx, aKeys[0], aKeys[1]).Result()
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "，差集结果是：", val)
	case "sdiffstore":
		aKeys := strings.Split(opt.redisKey, ",")
		val, err := clusterClient.SDiffStore(ctx, aKeys[0], aKeys[1], aKeys[2]).Result()
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "，查询出来的值为：", val)
	case "sinter":
		aKeys := strings.Split(opt.redisKey, ",")
		val, err := clusterClient.SInter(ctx, aKeys[0], aKeys[1]).Result()
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "，交集结果是：", val)
	case "sinterstore":
		aKeys := strings.Split(opt.redisKey, ",")
		val, err := clusterClient.SInterStore(ctx, aKeys[0], aKeys[1], aKeys[2]).Result()
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "，查询出来的值为：", val)
	case "sunion":
		aKeys := strings.Split(opt.redisKey, ",")
		val, err := clusterClient.SUnion(ctx, aKeys[0], aKeys[1]).Result()
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "，并集结果是：", val)
	case "sunionstore":
		aKeys := strings.Split(opt.redisKey, ",")
		val, err := clusterClient.SUnionStore(ctx, aKeys[0], aKeys[1], aKeys[2]).Result()
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "，查询出来的值为：", val)
	case "spop":
		val, err := clusterClient.SPop(ctx, opt.redisKey).Result()
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "，查询出来的值为：", val)
	case "spopn":
		val, err := clusterClient.SPopN(ctx, opt.redisKey, opt.startIndex).Result()
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "，查询出来的值为：", val)
	case "srem":
		val, err := clusterClient.SRem(ctx, opt.redisKey, opt.redisVal).Result()
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "，查询出来的值为：", val)
	case "srandermenmber":
		val, err := clusterClient.SRandMember(ctx, opt.redisKey).Result()
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "，查询出来的值为：", val)
	case "srandermenmbern":
		val, err := clusterClient.SRandMemberN(ctx, opt.redisKey, opt.startIndex).Result()
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "，查询出来的值为：", val)
	case "smenmbermap":
		val, err := clusterClient.SMembersMap(ctx, opt.redisKey).Result()
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "，查询出来的值为：", val)
	case "smove":
		aKeys := strings.Split(opt.redisKey, ",")
		ok, err := clusterClient.SMove(ctx, aKeys[0], aKeys[1], opt.redisVal).Result()
		PrintPanic(err)
		if ok {
			fmt.Println("移动数据成功")
		}
	}
}

//redis list操作
func listOperation() {
	// LPush LPushX LRange LLen
	fmt.Println("---------------------redis  LPush LPushX LRange LLen 使用---------------------")
	HandlerRedisListCommand(connOption{command: "lpush", redisKey: "lpush_test_name"})
	HandlerRedisListCommand(connOption{command: "lpushx", redisKey: "lpush_test_name"})
	HandlerRedisListCommand(connOption{command: "lrange", redisKey: "lpush_test_name", startIndex: 0, endIndex: -1})
	HandlerRedisListCommand(connOption{command: "llen", redisKey: "lpush_test_name"})
	fmt.Println()

	// LTrim & LIndex
	fmt.Println("---------------------redis  LTrim & LIndex 使用---------------------")
	HandlerRedisListCommand(connOption{command: "lindex", redisKey: "lpush_test_name", startIndex: 3})
	HandlerRedisListCommand(connOption{command: "ltrim", redisKey: "lpush_test_name", startIndex: 0, endIndex: 3})
	HandlerRedisListCommand(connOption{command: "lrange", redisKey: "lpush_test_name", startIndex: 0, endIndex: -1})
	fmt.Println()

	// LSet & LInsert
	fmt.Println("---------------------redis  LSet & LInsert 使用---------------------")
	HandlerRedisListCommand(connOption{command: "lset", redisKey: "lpush_test_name", startIndex: 2, redisVal: "beer"})
	HandlerRedisListCommand(connOption{command: "lrange", redisKey: "lpush_test_name", startIndex: 0, endIndex: -1})
	HandlerRedisListCommand(connOption{command: "linsert", redisKey: "lpush_test_name", lInsertType: "before", redisVal: "lili", lInsertVal: "beer"})
	HandlerRedisListCommand(connOption{command: "lrange", redisKey: "lpush_test_name", startIndex: 0, endIndex: -1})
	HandlerRedisListCommand(connOption{command: "linsert", redisKey: "lpush_test_name", lInsertType: "after", redisVal: "jack", lInsertVal: "beer"})
	HandlerRedisListCommand(connOption{command: "lrange", redisKey: "lpush_test_name", startIndex: 0, endIndex: -1})
	fmt.Println()

	// LPop & LRem
	fmt.Println("---------------------redis  LPop & LRem 使用---------------------")
	HandlerRedisListCommand(connOption{command: "lpop", redisKey: "lpush_test_name"})
	HandlerRedisListCommand(connOption{command: "lrange", redisKey: "lpush_test_name", startIndex: 0, endIndex: -1})
	HandlerRedisListCommand(connOption{command: "lrem", redisKey: "lpush_test_name", startIndex: 10, redisVal: "lili"})
	HandlerRedisListCommand(connOption{command: "lrange", redisKey: "lpush_test_name", startIndex: 0, endIndex: -1})
	fmt.Println()
}

//处理redis list操作
func HandlerRedisListCommand(opt connOption) {
	switch opt.command {
	case "lpush":
		clusterClient.LPush(ctx, opt.redisKey, CreateRandomString(20))
	case "lpushx":
		clusterClient.LPushX(ctx, opt.redisKey, CreateRandomString(20))
	case "lrange":
		val, err := clusterClient.LRange(ctx, opt.redisKey, opt.startIndex, opt.endIndex).Result()
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "，查询出来的值为：", val)
	case "llen":
		val, err := clusterClient.LLen(ctx, opt.redisKey).Result()
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "，查询出来的值为：", val)
	case "lindex":
		val, err := clusterClient.LIndex(ctx, opt.redisKey, opt.startIndex).Result()
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "，查询出来的值为：", val)
	case "ltrim":
		val := clusterClient.LTrim(ctx, opt.redisKey, opt.startIndex, opt.endIndex)
		fmt.Println("redis的key为：", opt.redisKey, "，查询出来的值为：", val)
	case "lset":
		clusterClient.LSet(ctx, opt.redisKey, opt.startIndex, opt.redisVal)
	case "linsert":
		clusterClient.LInsert(ctx, opt.redisKey, opt.lInsertType, opt.lInsertVal, opt.redisVal)
	case "lpop":
		clusterClient.LPop(ctx, opt.redisKey)
	case "lrem":
		clusterClient.LRem(ctx, opt.redisKey, opt.startIndex, opt.redisVal)
	}
}

//redis基本键值操作
func basicOperation() {
	//set & get
	fmt.Println("---------------------redis set & get 使用---------------------")
	HandlerRedisBasicCommand(connOption{command: "set", redisKey: "set_test_name"})
	HandlerRedisBasicCommand(connOption{command: "get", redisKey: "set_test_name"})
	fmt.Println()

	// GetSet & SetNX
	fmt.Println("---------------------redis  GetSet & SetNX 使用---------------------")
	HandlerRedisBasicCommand(connOption{command: "getset", redisKey: "set_test_name"})
	HandlerRedisBasicCommand(connOption{command: "setnx", redisKey: "set_test_name"})
	fmt.Println()

	// MGet & MSet
	fmt.Println("---------------------redis  MGet & MSet 使用---------------------")
	//mset在集群模式设置值时错误 CROSSSLOT Keys in request don’t hash to the same slot
	//解决方案HashTag,HashTag即是用{}包裹key的一个子串，如{user:}1, {user:}2。
	HandlerRedisBasicCommand(connOption{command: "mset", redisKey: "{mset_test_name_}1,{mset_test_name_}2,{mset_test_name_}3"})
	HandlerRedisBasicCommand(connOption{command: "mget", redisKey: "{mset_test_name_}1,{mset_test_name_}2,{mset_test_name_}3"})
	fmt.Println()

	// Incr IncrBy Decr DecrBy
	fmt.Println("---------------------redis  Incr IncrBy Decr DecrBy 使用---------------------")
	HandlerRedisBasicCommand(connOption{command: "set", redisKey: "set_incr_test_num", redisVal: 20})
	HandlerRedisBasicCommand(connOption{command: "get", redisKey: "set_incr_test_num"})
	HandlerRedisBasicCommand(connOption{command: "incr", redisKey: "set_incr_test_num"})
	HandlerRedisBasicCommand(connOption{command: "get", redisKey: "set_incr_test_num"})
	var num int64 = 5
	HandlerRedisBasicCommand(connOption{command: "incrby", redisKey: "set_incr_test_num", redisVal: num})
	HandlerRedisBasicCommand(connOption{command: "get", redisKey: "set_incr_test_num"})
	HandlerRedisBasicCommand(connOption{command: "decr", redisKey: "set_incr_test_num"})
	HandlerRedisBasicCommand(connOption{command: "get", redisKey: "set_incr_test_num"})
	num = 3
	HandlerRedisBasicCommand(connOption{command: "decrby", redisKey: "set_incr_test_num", redisVal: num})
	HandlerRedisBasicCommand(connOption{command: "get", redisKey: "set_incr_test_num"})
	fmt.Println()

	// Del & Expire & Append
	fmt.Println("---------------------redis  Del & Expire & Append 使用---------------------")
	HandlerRedisBasicCommand(connOption{command: "set", redisKey: "set_append_test_name"})
	HandlerRedisBasicCommand(connOption{command: "get", redisKey: "set_append_test_name"})
	HandlerRedisBasicCommand(connOption{command: "append", redisKey: "set_append_test_name", redisVal: "hahhaaa"})
	HandlerRedisBasicCommand(connOption{command: "get", redisKey: "set_append_test_name"})
	HandlerRedisBasicCommand(connOption{command: "exprie", redisKey: "set_append_test_name", redisExpire: 604800 * time.Second})
	HandlerRedisBasicCommand(connOption{command: "del", redisKey: "set_append_test_name"})
	fmt.Println()
}

//处理redis基本键值操作
func HandlerRedisBasicCommand(opt connOption) {
	switch opt.command {
	case "set":
		if opt.redisVal == nil {
			err = clusterClient.Set(ctx, opt.redisKey, CreateRandomString(20), redisTTL).Err()
		} else {
			err = clusterClient.Set(ctx, opt.redisKey, opt.redisVal, redisTTL).Err()
		}
		PrintPanic(err)
	case "get":
		val, err := clusterClient.Get(ctx, opt.redisKey).Result()
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "，查询出来的值为：", val)
	case "getset":
		var oldVal interface{}
		if opt.redisVal == nil {
			oldVal, err = clusterClient.GetSet(ctx, opt.redisKey, CreateRandomString(20)).Result()
		} else {
			oldVal, err = clusterClient.GetSet(ctx, opt.redisKey, opt.redisVal).Result()
		}
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "，查询出来的旧值为：", oldVal)
	case "setnx":
		//如果key不存在，则设置这个key的值,并设置key的失效时间。如果key存在，则设置不生效
		if opt.redisVal == nil {
			err = clusterClient.SetNX(ctx, opt.redisKey, CreateRandomString(20), redisTTL).Err()
		} else {
			err = clusterClient.SetNX(ctx, opt.redisKey, opt.redisVal, redisTTL).Err()
		}
		PrintPanic(err)
	case "mset":
		//批量设置key1对应的值为value1，key2对应的值为value2，key3对应的值为value3
		aKeys := strings.Split(opt.redisKey, ",")
		aKeyMap := make(map[string]interface{})
		for i := 0; i < len(aKeys); i++ {
			aKeyMap[aKeys[i]] = CreateRandomString(20)
		}
		err = clusterClient.MSet(ctx, aKeyMap).Err()
		PrintPanic(err)
	case "mget":
		// MGet函数可以传入任意个key，一次性返回多个值;这里Result返回两个值，第一个值是一个数组，第二个值是错误信息
		// 这里指定传参个数一定为3个，test用
		aKeys := strings.Split(opt.redisKey, ",")
		val, err := clusterClient.MGet(ctx, aKeys[0], aKeys[1], aKeys[2]).Result()
		PrintPanic(err)
		fmt.Println("redis的key为：", opt.redisKey, "，查询出来的值为：", val)
	case "incr":
		clusterClient.Incr(ctx, opt.redisKey)
	case "incrby":
		clusterClient.IncrBy(ctx, opt.redisKey, opt.redisVal.(int64))
	case "decr":
		clusterClient.Decr(ctx, opt.redisKey)
	case "decrby":
		clusterClient.DecrBy(ctx, opt.redisKey, opt.redisVal.(int64))
	case "append":
		clusterClient.Append(ctx, opt.redisKey, opt.redisVal.(string))
	case "exprie":
		clusterClient.Expire(ctx, opt.redisKey, opt.redisExpire)
	case "del":
		clusterClient.Del(ctx, opt.redisKey)
	}
}

//打印错误
func PrintPanic(err error) {
	if err != nil {
		panic(err)
	}
}

//随机指定长度字符串
func CreateRandomString(len int) string {
	var (
		res string
		str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	)
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		res += string(str[randomInt.Int64()])
	}
	return res
}
