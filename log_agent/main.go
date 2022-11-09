package main

import (
	"gopkg.in/ini.v1"
	"log_agent/config"
	"log_agent/etcd"
	"log_agent/kafka"
	"log_agent/logger"
	"log_agent/tail_task"
	"sync"
	"time"
)

var cfg = new(config.Config)

func main() {
	// 读取配置文件,初始化配置
	err := ini.MapTo(cfg, "./conf.ini")
	if err != nil {
		logger.Error("Init config failed, err:" + err.Error())
		return
	}

	// 初始化etcd
	err = etcd.Init(cfg.Etcd.Addr, time.Duration(cfg.Etcd.Timeout) * time.Second)
	if err != nil {
		logger.Error("Init Etcd failed, err:" + err.Error())
		return
	}
	logger.Info("Etcd init complete.")

	// 获取etcd的日志配置
	logConf, err := etcd.GetConf(cfg.Etcd.Key)
	if err != nil {
		logger.Error("Get Etcd config failed, err:" + err.Error())
		return
	}
	logger.Info("Etcd get config complete.")

	// 初始化Kafka
	err = kafka.Init(cfg.KafKa.Addr, cfg.KafKa.ChannelMaxSize, cfg.KafKa.ProcessNum)
	if err != nil {
		logger.Error("Init kafka failed, err:" + err.Error())
		return
	}
	logger.Info("Kafka init complete.")

	// 初始化所有任务
	tail_task.Init(logConf)

	var wg sync.WaitGroup
	wg.Add(1)

	// 监控etcd配置，进行任务热更新
	go etcd.WatchConf(cfg.Etcd.Key)

	wg.Wait()
}
