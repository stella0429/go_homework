package producer

import (
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
)

type KafkaCfg struct {
	Host  string `json:"host"`
	Topic string `json:"topic"`
}

type Reporter struct {
	Producer sarama.SyncProducer
}

func NewReporter(cfg *KafkaCfg) *Reporter {
	reporter := &Reporter{}
	reporter.setProducer(cfg)
	return reporter
}

func (reporter *Reporter) setProducer(cfg *KafkaCfg) {
	var broker = []string{cfg.Host}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回
	producer, err := sarama.NewSyncProducer(broker, config)
	if err != nil {
		panic(err)
	}
	reporter.Producer = producer
}

func (reporter *Reporter) DoReport(topic string, msg []byte) (interface{}, error) {
	return reporter.do(topic, msg)
}

func (reporter *Reporter) do(topic string, msg []byte) (interface{}, error) {
	res := make(map[string]interface{})
	kafkaMsg := generateProducerMessage(topic, msg)
	pid, offset, err := reporter.Producer.SendMessage(kafkaMsg)
	if err != nil {
		return res, errors.Wrap(err, "Send msg failed!")
	}
	res["pid"] = pid
	res["offset"] = offset
	return res, nil
}

func generateProducerMessage(topic string, message []byte) *sarama.ProducerMessage {
	return &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
}
