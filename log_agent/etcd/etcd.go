package etcd

import (
	"context"
	"encoding/json"
	"go.etcd.io/etcd/client/v3"
	"log_agent/config"
	"log_agent/logger"
	"log_agent/tail_task"
	"time"
)

var client *clientv3.Client

func Init(addr []string, timeout time.Duration) (err error) {
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   addr,
		DialTimeout: timeout,
	})

	return
}

func GetConf(key string) (logConf []*config.LogConf, err error) {
	response, err := client.Get(context.Background(), key)
	if err != nil {
		return nil, err
	}
	for _, kv := range response.Kvs {
		err = json.Unmarshal(kv.Value, &logConf)
		if err != nil {
			return nil, err
		}
	}
	return
}

func WatchConf(key string){
	watchChan := client.Watch(context.Background(), key)
	for watchResp := range watchChan {
		for _, ev := range watchResp.Events {
			var logConf []*config.LogConf
			if ev.Type == clientv3.EventTypePut{
				err := json.Unmarshal(ev.Kv.Value, &logConf)
				if err != nil {
					logger.Error("Watch etcd config error, err:" + err.Error())
				}
			}
			tail_task.UpdateTask(logConf)
		}
	}
}
