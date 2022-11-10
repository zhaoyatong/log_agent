package es

import (
	"context"
	"github.com/olivere/elastic/v7"
	"log_agent/logger"
	"time"
)

type LogData struct {
	Index string
	Data  *LogInfo
}

type LogInfo struct {
	Message string `json:"message"`
	Time    string `json:"time"`
}

var (
	client  *elastic.Client
	logChan chan *LogData
)

func Init(addr string, process int) (err error) {
	client, err = elastic.NewClient(elastic.SetURL(addr))
	logChan = make(chan *LogData, 100000)

	for i := 0; i < process; i++ {
		go run()
	}

	return
}

func SendToChan(logData *LogData) {
	logChan <- logData
}

func run() {
	for {
		select {
		case log := <-logChan:
			_, err := client.Index().
				Index(log.Index).
				BodyJson(*log.Data).
				Do(context.Background())
			if err != nil {
				logger.Error("Send data to ES failed, err:" + err.Error())
			}
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}
