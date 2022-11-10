package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/ini.v1"
	"log_agent/config"
	"log_agent/es"
	"log_agent/logger"
	"log_agent/tail_task"
)

var cfg = new(config.Config)

func main() {
	// 读取配置文件,初始化配置
	err := ini.MapTo(cfg, "./conf.ini")
	if err != nil {
		logger.Error("Init config failed, err:" + err.Error())
		return
	}

	// 初始化ES
	err = es.Init(cfg.ESConfig.Addr, cfg.ESConfig.Process)
	if err != nil {
		logger.Error("Init ES failed, err:" + err.Error())
		return
	}
	logger.Info("ES init complete.")

	// 读取任务配置的json
	taskConf := make([]config.TaskConf, 0, 2)
	err = json.Unmarshal([]byte(cfg.Task.Conf), &taskConf)
	if err != nil {
		logger.Error("Analysis task config json failed, err:" + err.Error())
		return
	}

	// 初始化tail任务
	allFail := true
	for _, conf := range taskConf {
		err = tail_task.Init(conf.Path, conf.Index)
		if err != nil {
			logger.Error(
				fmt.Sprintf("Init tail task failed, task:%s %s, err:%s", conf.Path, conf.Index, err.Error()),
				)
			continue
		}
		allFail = false
	}
	if allFail{
		logger.Info("All task start failed.")
		return
	}
	logger.Info("Tail task init complete.")
	logger.Info("Task running.")

	select {}
}
