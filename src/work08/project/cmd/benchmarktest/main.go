package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"
	rdb "project/project/pkg/redisdb"
	"strconv"
	"time"
)

//自己用来随便测试玩的
const (
	address = "127.0.0.1:6379"
	pwd     = ""
)

var index = 0

//基础版本的redis初始化
func init() {
	redisConf := &rdb.RedisConf{
		Addr:     address,
		Password: pwd,
		DB:       0,
	}
	if err := rdb.RedisDial(redisConf); err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func main() {
	RedisBenchMark(10)
	RedisBenchMark(20)
	RedisBenchMark(50)
	RedisBenchMark(100)
	RedisBenchMark(200)
	RedisBenchMark(1000)
	RedisBenchMark(5000)
}

func RedisBenchMark(len int) {
	//随机指定长度字符串
	string := CreateRandomString(len)
	key := "test_redis1_" + strconv.Itoa(index)

	//统计set使用时间
	startSet := time.Now()
	rdb.Set(key, string, 120)
	elapsedSet := time.Since(startSet)
	fmt.Println("测试value大小为：", len, "字节，set执行时间为：", elapsedSet)

	//统计get使用时间
	startGet := time.Now()
	_ = rdb.Get(key)
	elapsedGet := time.Since(startGet)
	fmt.Println("测试value大小为：", len, "字节，get执行时间为：", elapsedGet)

	index++
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
