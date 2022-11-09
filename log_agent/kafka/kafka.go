package kafka

import (
	"github.com/Shopify/sarama"
	"log_agent/logger"
)

type logData struct {
	topic string
	data  string
}

var (
	producer    sarama.SyncProducer
	consumer    sarama.Consumer
	logDataChan chan *logData
)

func Init(addr []string, maxSize, processNum int) (err error) {
	// 生产者配置
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送数据需leader和follow确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新随机一个partition
	config.Producer.Return.Successes = true                   // 回复成功消息

	producer, err = sarama.NewSyncProducer(addr, config)
	if err != nil{
		return
	}

	logDataChan = make(chan *logData, maxSize)

	for i := 0; i < processNum; i++ {
		go sendToKafka()
	}

	return
}

func SendToChan(topic, data string) {
	logDataChan <- &logData{
		topic: topic,
		data:  data,
	}
}

func sendToKafka() {
	for logData := range logDataChan {
		msg := &sarama.ProducerMessage{
			Topic: logData.topic,
			Value: sarama.StringEncoder(logData.data),
		}
		_, _, err := producer.SendMessage(msg)
		if err != nil {
			logger.Error("Send data to kafka failed, err:" + err.Error())
		}
	}
}
