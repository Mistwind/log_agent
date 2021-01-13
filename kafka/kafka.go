package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
)

var client sarama.SyncProducer

// Init 初始化配置
func Init(addrs []string) (err error) {
	config := sarama.NewConfig()

	// 副本数据同步的等级
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	client, err = sarama.NewSyncProducer(addrs, config)
	if err != nil {
		fmt.Printf("producer closed, err: %v\n", err)
		return
	}
	return
}

// SendToKafka 发送消息到kafka
func SendToKafka(topic, data string) {
	// 构建发送到 kafka 的消息
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(data),
	}
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Printf("Send msg failed, err: %v\n", err)
	}
	fmt.Printf("pid: %v  offset: %v\n", pid, offset)
	fmt.Println("发送成功")
}
