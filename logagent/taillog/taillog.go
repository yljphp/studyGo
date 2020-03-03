package taillog

import (
	"fmt"
	"github.com/hpcloud/tail"
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

	return
}

func (t *TailTask) ReadChan() <-chan *tail.Line {
	return t.instance.Lines
}
