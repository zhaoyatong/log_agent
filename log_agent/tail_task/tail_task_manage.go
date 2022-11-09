package tail_task

import (
	"context"
	"github.com/hpcloud/tail"
	"log_agent/config"
	"log_agent/logger"
)

// 日志收集任务
type tailTask struct {
	filename string             // 文件名
	topic    string             // kafka的主题topic
	instance *tail.Tail         //tail客户端实例
	retain   bool               //是否保留任务，当etcd配置任务更新时，false的任务将被cancel
	ctx      context.Context    // context上下文
	cancel   context.CancelFunc // cancel函数
}

var (
	taskList = make(map[string]*tailTask, 16) // 任务列表
)

func Init(logConf []*config.LogConf) {
	for _, conf := range logConf {
		taskObj, err := newTailTask(conf.Filename, conf.Topic)
		if err != nil {
			logger.Error("Create task failed, err:" + err.Error())
			continue
		}
		taskList[conf.Filename+"_"+conf.Topic] = taskObj
	}
}

func newTailTask(filename, topic string) (taskObj *tailTask, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	taskObj = &tailTask{
		filename: filename,
		topic:    topic,
		retain:   true,
		ctx:      ctx,
		cancel:   cancel,
	}

	err = taskObj.init()
	return
}

func UpdateTask(newConf []*config.LogConf) {
	// 先将保留状态置为否
	for _, task := range taskList {
		task.retain = false
	}

	for _, conf := range newConf {
		_, ok := taskList[conf.Filename+"_"+conf.Topic]
		// 存在设置为保留任务，不存在则新增任务
		if ok {
			taskList[conf.Filename+"_"+conf.Topic].retain = true
		} else {
			taskObj, err := newTailTask(conf.Filename, conf.Topic)
			if err != nil {
				logger.Error("Create task failed, err:" + err.Error())
				continue
			}
			taskList[conf.Filename+"_"+conf.Topic] = taskObj
		}
	}
	// 移除不保留的任务
	for _, task := range taskList {
		if !task.retain {
			task.cancel()
		}
	}
}
