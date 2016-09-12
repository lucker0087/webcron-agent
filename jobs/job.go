package jobs

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"time"
	"webcron-agent/config"
	"webcron-agent/libs"
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

type JobResult struct {
	Id            int    `json:job_id`
	SuccessOutput string `json:response_success`
	ErrorOutput   string `json:response_error`
	Error         string `json:error`
	Consume       int    `json:consume`
	Timeout       bool   `json:timeout`
}

func (j *Job) Run() {

	var logger = libs.NewTaskLogger()
	if !j.Concurrent && j.status > 0 {
		logger.Warning(fmt.Sprintf("任务[%d]上一次执行尚未结束，本次被忽略。", j.id))
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

	j.status++

	logger.Info(fmt.Sprintf("开始执行任务: %s, 当前任务数: %d\n", j.name, j.status))
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

	var result JobResult

	result.Id = j.id
	result.SuccessOutput = cmdOut
	result.ErrorOutput = cmdErr
	result.Timeout = isTimeout
	result.Error = ""
	result.Consume = int(consume)
	if err != nil {
		result.Error = err.Error()
	}

	data, _ := encodeResult(result)

	config, _ := config.GetConfig()
	client, err := NewSocketClient(config.Master.Server, config.Master.Port)
	if err != nil {
		result.Error = "与服务器失去连接, 执行结果未同步"
		data, _ := encodeResult(result)
		//log
		logger.Warning(string(data))
		return
	}
	defer client.Close()

	err = client.Send(data)
	if err != nil {
		result.Error = err.Error()
		data, _ := encodeResult(result)
		//log
		fmt.Println(string(data))
		return
	}
	fmt.Println(string(data))
	return
}

func encodeResult(result JobResult) ([]byte, error) {
	var data Data
	data.Type = "job_response"
	data.Data = result
	data.Time = time.Now().String()
	return json.Marshal(data)
}
