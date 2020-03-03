package taillog

import (
	"fmt"
	"logagent/etcd"
)

var (
	tskMgr *tailLogMgr
)

// tailTask 管理者
type tailLogMgr struct {
	logEntry []*etcd.LogEntry
}

func Init(logEntryConf []*etcd.LogEntry) (err error) {
	tskMgr = &tailLogMgr{
		logEntry: logEntryConf, // 把当前的日志收集项配置信息保存起来
	}

	for _, logEntry := range logEntryConf {

		_, err = NewTailTask(logEntry.Path, logEntry.Topic)

		if err != nil {
			fmt.Printf("taillog new tail task failed,err:%v\n", err)
			return
		}

	}

	return
}
