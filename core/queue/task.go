package queue

import (
	"daosuan/utils/log"
	"fmt"
	"time"
)

var Task chan func()

// 初始化任务队列
func InitTaskQueue() {

	Task = make(chan func(), 100)
	go func() {
		for {
			function := <-Task
			startTime := time.Now()
			logUtils.Println("启动任务")
			logUtils.Println(fmt.Sprintf("当前任务队列等待数%d", len(Task)))
			function()
			logUtils.Println(fmt.Sprintf("任务结束，耗时: %s", time.Now().Sub(startTime).String()))
		}
	}()
}