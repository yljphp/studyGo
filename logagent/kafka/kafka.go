package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
)

var(
	// 声明一个全局的连接kafka的生产者client
	client sarama.SyncProducer
)

func Init(addrs []string)(err error) {

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
	return
}
