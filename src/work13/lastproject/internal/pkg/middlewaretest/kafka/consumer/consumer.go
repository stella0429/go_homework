package consumer

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/my/repo/internal/pkg/middlewaretest/kafka/producer"
	"github.com/pkg/errors"
)

type Subscriber struct {
	Consumer sarama.Consumer
}

func NewSubscriber(cfg *producer.KafkaCfg) *Subscriber {
	subscriber := &Subscriber{}
	subscriber.setConsumer(cfg)
	return subscriber
}

func (subscriber *Subscriber) setConsumer(cfg *producer.KafkaCfg) {
	consumer, err := sarama.NewConsumer([]string{cfg.Host}, nil)
	if err != nil {
		panic(err)
	}
	subscriber.Consumer = consumer
}

func (subscriber *Subscriber) Consume(topic string, ch chan string) {
	defer func() {
		if err := subscriber.Consumer.Close(); err != nil {
			fmt.Println(errors.Wrap(err, "Fail to start consumer!"))
		}
	}()

	partitionList, err := subscriber.Consumer.Partitions(topic)
	if err != nil {
		fmt.Println(errors.Wrapf(err, "Fail to get list of partition!"))
		return
	}

	for _, partition := range partitionList {
		pc, err := subscriber.Consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			fmt.Println(errors.Wrapf(err, "Failed to start consumer for partition:%d", partition))
			return
		}
		go func(pc sarama.PartitionConsumer) {
			for message := range pc.Messages() {
				fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v", message.Partition, message.Offset, message.Key, message.Value)
				ch <- messageReceived(message)
			}
		}(pc)
	}
}

func messageReceived(message *sarama.ConsumerMessage) string {
	return string(message.Value)
}
