package logger

import (
	"fmt"
	"time"
)

func Info(data string){
	fmt.Printf("[%s] Info:%s\n", time.Now().Format("2006-01-02 15:04:05.000"), data)
}

func Error(data string){
	fmt.Printf("[%s] Error:%s\n", time.Now().Format("2006-01-02 15:04:05.000"), data)
}
