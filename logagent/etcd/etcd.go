package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

var (
	cli *clientv3.Client
)

type LogEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

func Init(addr []string, timeout time.Duration) (err error) {

	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   addr,
		DialTimeout: timeout * time.Second,
	})

	if err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	return

}

func GetConf(key string) (logEntryConf []*LogEntry,err error) {

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)

	resp, err := cli.Get(ctx, key)
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return
	}

	cancelFunc()

	for _, ev := range resp.Kvs {
		//fmt.Printf("%s:%s\n", ev.Key, ev.Value)
		err = json.Unmarshal(ev.Value, &logEntryConf)
		if err != nil {
			fmt.Printf("unmarshal etcd value failed,err:%v\n", err)
			return
		}
	}
	return
}
