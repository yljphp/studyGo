package main

import (
	"fmt"
	"github.com/hpcloud/tail"
	"io"
	"time"
)

func main() {

	InitTail("/tmp/taillog/test.log")
}

func InitTail(path string) {

	tail, err := tail.TailFile(path, tail.Config{
		ReOpen:    true,                                          // 重新打开
		Follow:    true,                                          // 是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd}, // 从文件的哪个地方开始读
		MustExist: false,                                         // 文件不存在不报错
		Poll:      true,
	})

	if err != nil {
		fmt.Printf("init tail failed,err:%v\n", err)
		return
	}

	for {
		select {
		case line := <-tail.Lines:
			fmt.Println(line.Text)
		default:
			time.Sleep(time.Second)
		}
	}

}
