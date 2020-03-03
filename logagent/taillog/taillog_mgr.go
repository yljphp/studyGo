package taillog

import (
	"fmt"
	"logagent/etcd"
	"logagent/kafka"
	"time"
)

var (
	tskMgr *tailLogMgr
)

// tailTask 管理者
type tailLogMgr struct {
	logEntry []*etcd.LogEntry
}

func Init(logEntryConf []*etcd.LogEntry) {
	tskMgr = &tailLogMgr{
		logEntry: logEntryConf, // 把当前的日志收集项配置信息保存起来
	}



	go tskMgr.run()
}

func (t *tailLogMgr) run() {

	for _, logEntry := range t.logEntry {

		// 初始化的时候起了多少个tailtask 都要记下来，为了后续判断方便
		taskObj, err := NewTailTask(logEntry.Path, logEntry.Topic)
		if err != nil {
			fmt.Printf("taillog new tail task failed,err:%v\n", err)
			return
		}

		for {

			select {

			case line := <-taskObj.instance.Lines:
				kafka.SendMsg(taskObj.topic, line.Text)
			default:
				time.Sleep(time.Second)

			}

		}

	}

}
