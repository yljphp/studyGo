package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

var (
	cli *clientv3.Client
)

func InitEtcd() (err error) {
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		fmt.Printf("init etcd failed,err:%v\n", err)
		return
	}

	return
}

func main() {

	err := InitEtcd()
	if err != nil {
		fmt.Printf("init etcd failed,err:%v\n", err)
		return
	}

	fmt.Println("connect to etcd success")


	//ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)


	info := `[{"path":"/tmp/web.log","topic":"web_log"},{"path":"/tmp/nginx.log","topic":"nginx_log"}]`

	resp, err := cli.Put(context.TODO(), "/logagent/config", info)
	if err != nil {
		fmt.Printf("etcd.putInfo failed,err:%v\n", err)
		return
	}

	fmt.Println(resp)
}
