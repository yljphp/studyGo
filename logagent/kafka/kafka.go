package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

var (
	// 声明一个全局的连接kafka的生产者client
	client      sarama.SyncProducer
	logDataChan chan *logData
)

type logData struct {
	topic string
	data  string
}

func Init(addrs []string, chanMaxSize int) (err error) {

	config := sarama.NewConfig()

	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出⼀个 partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

	//连接kafka
	client, err = sarama.NewSyncProducer(addrs, config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
		return
	}

	//初始化chan
	logDataChan = make(chan *logData, chanMaxSize)

	//后台启动消费任务
	go sendMsg()

	return
}

func SendChan(topic, data string) {

	logDataChan <- &logData{
		topic,
		data,
	}
}

/**
真正的往kafka中扔数据

*/
func sendMsg() {

	for {

		select {
		case tmpLog := <-logDataChan:

			message := &sarama.ProducerMessage{
				Topic: tmpLog.topic,
				Value: sarama.StringEncoder(tmpLog.data),
			}

			pid, offset, err := client.SendMessage(message)
			if err != nil {
				fmt.Println("send msg failed, err:", err)
				return
			}
			fmt.Printf("pid:%v offset:%v\n", pid, offset)

		default:
			time.Sleep(time.Millisecond * 50)

		}

	}

}
