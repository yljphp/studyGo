package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"logagent/conf"
	"logagent/etcd"
	"logagent/kafka"
	"logagent/taillog"
	"time"
)

var (
	cfg = new(conf.AppConf)
)

func main() {

	//1 加载配置文件
	err := ini.MapTo(cfg, "./conf/config.ini")
	if err != nil {
		fmt.Printf("load ini failed, err:%v\n", err)
		return
	}

	fmt.Println("load ini success.")

	//2 初始化kafka连接
	err = kafka.Init([]string{cfg.KafkaConf.Address})
	if err != nil {
		fmt.Printf("init Kafka failed,err:%v\n", err)
		return
	}

	fmt.Println("init kafka success.")

	//3 初始化etcd
	err = etcd.Init([]string{cfg.EtcdConf.Address}, time.Duration(cfg.EtcdConf.Timeout))
	if err != nil {
		fmt.Printf("init etcd failed,err:%v\n", err)
		return
	}
	fmt.Println("init etcd success.")

	//3.1从etcd中获取日志收集项的配置信息
	logEntryConf, err := etcd.GetConf(cfg.EtcdConf.LogKey)
	if err != nil {
		fmt.Printf("etcd.GetConf failed,err:%v\n", err)
		return
	}

	fmt.Printf("get conf from etcd success, %v\n", logEntryConf)

	//3.2 派一个哨兵去监视日志收集项的变化（有变化及时通知我的logAgent实现热加载配置）

	//4 收集日志发往kafka
	for index, value := range logEntryConf {

		fmt.Printf("index:%v value:%v\n", index, value)

		task, err := taillog.NewTailTask(value.Path, value.Topic)
		if err != nil {
			fmt.Printf("taillog new tail task failed,err:%v\n", err)
			return
		}

		for {
			select {
			case line := <-task.ReadChan():

				fmt.Println("get tail msg text: ", line.Text)

				err := kafka.SendMsg(value.Topic, line.Text)
				if err != nil {
					fmt.Printf("send msg to kafka send ,err:%v\n", err)
					continue
				}
			default:
				time.Sleep(time.Second)
			}

		}

	}

}
