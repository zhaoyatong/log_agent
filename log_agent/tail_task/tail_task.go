package tail_task

import (
	"github.com/hpcloud/tail"
	"log_agent/es"
	"time"
)

func Init(filename, index string) error {
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Whence: 2},
		MustExist: false,
		Poll:      true,
	}

	tailObj, err := tail.TailFile(filename, config)
	if err != nil {
		return err
	}

	// 开始执行任务
	go func() {
		for {
			select {
			case line := <-tailObj.Lines:
				info := &es.LogInfo{
					Message: line.Text,
					Time:    time.Now().Format("2006-01-02 15:04:05.000"),
				}

				logData := &es.LogData{
					Index: index,
					Data:  info,
				}

				es.SendToChan(logData)
			default:
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()

	return nil
}
