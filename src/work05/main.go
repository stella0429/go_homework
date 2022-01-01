package main

import (
	"fmt"
	"math/rand"
	"time"
)

type reqStatistics struct {
	totalReq   int
	successReq int
}

type sumStatistics struct {
	sumReq     reqStatistics
	indexStamp int
}

var (
	bucketNum  = 10
	bucketMap  = make(map[int]reqStatistics)
	bucketInMs = 100
	sumBucket  = &sumStatistics{}
	baseTime   = makeTimeStamp()
)

//获取当前毫秒时间戳
func makeTimeStamp() int {
	return int(time.Now().UnixNano() / 1e6)
}

//简易版滑动窗口计数器
func bucketedCounter() {
	var (
		nowTime   = makeTimeStamp()
		timeDiff  = nowTime - baseTime
		tempIndex = baseTime + 100*(timeDiff/bucketInMs)
	)
	if _, ok := bucketMap[tempIndex]; !ok {
		bucketMap[tempIndex] = reqStatistics{}
	}

	//记录总请求数和成功请求数
	tempValue := bucketMap[tempIndex]
	tempValue.totalReq += 1
	//随机种子模拟请求成功or失败
	rand.Seed(time.Now().UnixNano())
	seedTemp := rand.Intn(100)
	//假设成功概率为70%
	if seedTemp <= 70 {
		tempValue.successReq += 1
	}
	bucketMap[tempIndex] = tempValue

	//如果超过了整个桶，对历史10个桶做统计
	if timeDiff >= bucketInMs*bucketNum {
		//更新最新基准时间（考虑有可能有跳过的情况）
		baseTime = baseTime + ((timeDiff/bucketInMs)-bucketNum)*bucketInMs
		//判断是否已经计算过了，计算过的就不需要更新了,否则更新统计数据
		if sumBucket.indexStamp < baseTime {
			sumBucket.indexStamp = baseTime
			sumBucket.sumReq.totalReq = 0
			sumBucket.sumReq.successReq = 0
			for k, v := range bucketMap {
				if k >= baseTime && k < baseTime+bucketNum*bucketInMs {
					sumBucket.sumReq.totalReq += v.totalReq
					sumBucket.sumReq.successReq += v.successReq
				}
				//删除过期的桶
				if k < baseTime {
					delete(bucketMap, k)
				}
			}
			fmt.Println("bucketMap value is :", bucketMap)
			fmt.Println("sumBucket value is : ", sumBucket)
			fmt.Println()
		}
	}
}

func main() {
	//模拟请求调用，每10ms请求一次
	for {
		time.Sleep(time.Duration(10 * time.Millisecond))
		bucketedCounter()
	}
}
