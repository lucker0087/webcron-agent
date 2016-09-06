package jobs

import (
	"fmt"
	"runtime/debug"
	"time"
)

type Job struct {
	id         int                                               // 任务ID
	logId      int64                                             // 日志记录ID
	name       string                                            // 任务名称
	task       *Task                                             // 任务对象
	runFunc    func(time.Duration) (string, string, error, bool) // 执行函数
	status     int                                               // 任务状态，大于0表示正在执行中
	Concurrent bool                                              // 同一个任务是否允许并行执行
}

func (j *Job) Run() {
	if !j.Concurrent && j.status > 0 {
		fmt.Printf("任务[%d]上一次执行尚未结束，本次被忽略。", j.id)
		return
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf(string(debug.Stack()))
		}
	}()

	if workPool != nil {
		workPool <- true
		defer func() {
			<-workPool
		}()
	}

	fmt.Printf("开始执行任务: %s\n", j.name)

	j.status++
	defer func() {
		j.status--
	}()

	t := time.Now()
	timeout := time.Duration(time.Hour * 24)
	if j.task.Timeout > 0 {
		timeout = time.Second * time.Duration(j.task.Timeout)
	}

	cmdOut, cmdErr, err, isTimeout := j.runFunc(timeout)
	consume := time.Now().Sub(t) / time.Millisecond

	//client, err := NewSocketClient("127.0.0.1", 9999)
	//if err != nil {
	//	//todo
	//}

	if !isTimeout {
		fmt.Printf("任务:%s, 正常输出:%s, 异常输出:%s, 错误:%s, 执行时间:%d", j.name, cmdOut, cmdErr, err, int(consume))
	} else {
		fmt.Printf("任务:%s, 执行超过了%d秒, 异常输出:%s", j.name, timeout/time.Second, cmdErr)
	}
}
