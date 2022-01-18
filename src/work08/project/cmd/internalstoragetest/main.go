package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"
	rdb "project/project/pkg/redisdb"
)

const (
	address = "127.0.0.1:6379"
	pwd     = ""
	keyLen  = 50
	reps    = 10
)

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
	//分析1w数据大小内存占用
	calculateAverageMemory(10000)

	//分析2w数据大小内存占用
	calculateAverageMemory(20000)

	//分析5w数据大小内存占用
	calculateAverageMemory(50000)

	//分析10w数据大小内存占用
	calculateAverageMemory(100000)

	//分析20w数据大小内存占用
	calculateAverageMemory(200000)

	//分析30w数据大小内存占用
	calculateAverageMemory(300000)

	//分析50w数据大小内存占用
	calculateAverageMemory(500000)
}

//计算平均每个key占用内存空间
func calculateAverageMemory(len int) {
	fmt.Println("******************分析开始******************")

	var sum int64 = 0
	for i := 0; i < reps; i++ {
		sum += analyseRedisMemory(len)
	}

	fmt.Println("******************分析结果******************")
	fmt.Println("测试value大小为：", len, "字节，平均每个 key 的占用内存空间", sum/int64(reps), "\n\n")
}

//比较redis set值前后的内存差异，进行分析
func analyseRedisMemory(len int) int64 {
	//随机指定长度字符串以及key
	string := CreateRandomString(len)
	key := CreateRandomString(keyLen)

	//统计设置值前后内存大小
	startTotal, _ := rdb.DoInfo()
	//fmt.Println("startTotal:", startTotal)
	rdb.Set(key, string, 120)
	endTotal, _ := rdb.DoInfo()
	//fmt.Println("endTotal:", endTotal)

	//计算前后内存差
	diff := endTotal - startTotal
	fmt.Println("测试的key值为：", key, "value大小为：", len, "字节，占用内存大小为：", diff)
	return diff
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
