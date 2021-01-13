package taillog

import (
	"fmt"

	"github.com/hpcloud/tail"
)

var (
	tailObj *tail.Tail
	logChan chan string
)

// Init 初始化tail的配置
func Init(filename string) (err error) {
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		MustExist: false,
		Poll:      true,
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: 2,
		},
	}
	tailObj, err = tail.TailFile(filename, config)
	if err != nil {
		fmt.Printf("Tail file failed, err: %v\n", err)
		return
	}
	return
}

// ReadChan 读取数据写入channel
func ReadChan() <-chan *tail.Line {
	return tailObj.Lines
}
