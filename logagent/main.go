package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"logagent/conf"
	"logagent/etcd"
	"logagent/kafka"
	"logagent/taillog"
	"sync"
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
	err = kafka.Init([]string{cfg.KafkaConf.Address}, cfg.KafkaConf.ChanMaxSize)
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
	//4.1 循环每个日志收集项，创建TailObj
	err = taillog.Init(logEntryConf)
	if err != nil {
		fmt.Printf("init taillog failed,err:%v\n", err)
		return
	}

	fmt.Println("init taillog success.")

	var wg sync.WaitGroup
	wg.Add(1)

	wg.Wait()


}
