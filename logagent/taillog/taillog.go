package taillog

import (
	"fmt"
	"github.com/hpcloud/tail"
	"logagent/kafka"
	"time"
)

// TailTask： 一个日志收集的任务
type TailTask struct {
	path     string
	topic    string
	instance *tail.Tail
}

func NewTailTask(path, topic string) (tailObj *TailTask, err error) {

	tailObj = &TailTask{
		path:  path,
		topic: topic,
	}

	err = tailObj.init()

	return
}

func (t *TailTask) init() (err error) {
	config := tail.Config{
		ReOpen:    true,                                 // 重新打开
		Follow:    true,                                 // 是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件的哪个地方开始读
		MustExist: false,                                // 文件不存在不报错
		Poll:      true,
	}

	t.instance, err = tail.TailFile(t.path, config)
	if err != nil {
		fmt.Println("tail file failed, err:", err)
	}

	go t.run()

	return
}

func (t *TailTask) run() {

	for {

		select {

		case line := <-t.instance.Lines:

			//把日志发送到通道中
			//kafka中有单独的goroutine去取日志，然后发到kafka中
			kafka.SendChan(t.topic, line.Text)

		default:
			time.Sleep(time.Millisecond * 50)

		}

	}
}

func (t *TailTask) ReadChan() <-chan *tail.Line {
	return t.instance.Lines
}
