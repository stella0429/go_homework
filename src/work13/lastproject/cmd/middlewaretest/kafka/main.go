package main

import (
	"encoding/json"
	"fmt"
	"github.com/my/repo/internal/pkg/middlewaretest/kafka/consumer"
	"github.com/my/repo/internal/pkg/middlewaretest/kafka/producer"
)

var (
	//todo ip修改
	config = `{
	"host": "localhost:9092",  
    "topic": "kafka_test"
  }`
)

func main() {
	var cfg producer.KafkaCfg
	json.Unmarshal([]byte(config), &cfg)

	//reporter
	myproducer := producer.NewReporter(&cfg)
	//subscriber
	myconsumer := consumer.NewSubscriber(&cfg)
	message := "Hello Kafka World."

	ch := make(chan string)
	myconsumer.Consume(cfg.Topic, ch)
	res, err := myproducer.DoReport(cfg.Topic, []byte(message))
	fmt.Println("生产结果：", res, err)

	select {
	case msg := <-ch:
		fmt.Println("Got msg: ", msg)
		break
	}
}
