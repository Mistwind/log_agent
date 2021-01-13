package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
)

var (
	client      sarama.SyncProducer
	logDataChan chan *LogData
)

// LogData 日志数据
type LogData struct {
	topic string
	data  string
}

// Init 初始化配置
func Init(addrs []string, chanMaxSize int) (err error) {
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

	// 初始化 logDataChan
	logDataChan = make(chan *LogData, chanMaxSize)
	// 开启后台的 goroutine，从通道中取数据发往 kafka
	go SendToKafka()
	return
}

// SendToChan 将日志数据发送到一个内部的 channel 中
func SendToChan(topic, data string) {
	msg := &LogData{
		topic: topic,
		data:  data,
	}
	logDataChan <- msg
}

// SendToKafka 发送消息到kafka
func SendToKafka() {
	for {
		select {
		case logData := <-logDataChan:
			msg := &sarama.ProducerMessage{
				Topic: logData.topic,
				Value: sarama.StringEncoder(logData.data),
			}
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				fmt.Printf("Send msg failed, err: %v\n", err)
			}
			fmt.Printf("Send msg success, pid: %v  offset: %v\n", pid, offset)
		}
	}
}
