package tail_task

import (
	"fmt"
	"github.com/hpcloud/tail"
	"log_agent/kafka"
	"log_agent/logger"
	"time"
)

func (t *tailTask) init() (err error) {
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Whence: 2},
		MustExist: false,
		Poll:      true,
	}

	t.instance, err = tail.TailFile(t.filename, config)
	if err != nil{
		return
	}

	// 开始执行任务
	go t.run()

	return
}

func (t tailTask) run() {
	for {
		select {
		case <-t.ctx.Done():
			logger.Info(fmt.Sprintf("Task filename:%s,topic:%s exited.", t.filename, t.topic))
			return
		case line := <-t.instance.Lines:
			kafka.SendToChan(t.topic, line.Text)
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}
