package jobs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
	"time"
)

const (
	TASK_SUCCESS = 0  // 任务执行成功
	TASK_ERROR   = -1 // 任务执行出错
	TASK_TIMEOUT = -2 // 任务执行超时
)

type Task struct {
	Id           int    `json:"id"`
	UserId       int    `json:"user_id"`
	ServerId     int    `json:"server_id"`
	TaskName     string `json:"task_name"`
	TaskType     int    `json:"task_type"`
	Description  string `json:"description"`
	CronSpec     string `json:"cron_spec"`
	Concurrent   int    `json:"concurrent"`
	Command      string `json:"command"`
	Status       int    `json:"status"`
	Notify       int    `json:"notify"`
	NotifyEmail  string `json:"notify_email"`
	Timeout      int    `json:"timeout"`
	ExecuteTimes int    `json:"execute_times"`
	PrevTime     int64  `json:"prev_time"`
	CreateTime   int64  `json:"create_time"`
}

type Job struct {
	id         int                                               // 任务ID
	logId      int64                                             // 日志记录ID
	name       string                                            // 任务名称
	task       *Task                                             // 任务对象
	runFunc    func(time.Duration) (string, string, error, bool) // 执行函数
	status     int                                               // 任务状态，大于0表示正在执行中
	Concurrent bool                                              // 同一个任务是否允许并行执行
}

func RunTask(data []byte, remeber bool) error {
	task_list, err := DecodeTask(data)
	if err != nil {
		return err
	}

	if remeber {
		go RememberCron(data)
	}

	for _, task := range task_list {
		job, err := NewJob(task)
		if err != nil {
			log.Printf("InitJobs error :%s", err.Error())
			continue
		}
		AddJob(task.CronSpec, job)
	}
	return nil
}

func StopTasks(data []byte) error {
	var ids = make([]int, 0)
	err := json.Unmarshal(data, &ids)
	if err != nil {
		return err
	}

	if len(ids) == 0 {
		return nil
	}

	for _, id := range ids {
		RemoveJob(id)
		fmt.Printf("停止任务:%d\n", id)
	}

	return nil
}

func NewJob(task *Task) (*Job, error) {
	if task.Id < 1 {
		return nil, fmt.Errorf("ToJob: 缺少id")
	}
	job := &Job{
		id:   task.Id,
		name: task.TaskName,
	}
	job.task = task
	job.runFunc = func(timeout time.Duration) (string, string, error, bool) {
		bufOut := new(bytes.Buffer)
		bufErr := new(bytes.Buffer)
		cmd := exec.Command("/bin/bash", "-c", task.Command)
		cmd.Stdout = bufOut
		cmd.Stderr = bufErr
		cmd.Start()
		err, isTimeout := runCmdWithTimeout(cmd, timeout)

		return bufOut.String(), bufErr.String(), err, isTimeout
	}
	return job, nil
}

func DecodeTask(d []byte) ([]*Task, error) {
	var task_list []*Task
	err := json.Unmarshal(d, &task_list)
	if err != nil {
		return nil, err
	}
	return task_list, nil
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

	fmt.Printf("开始执行任务: %s", j.name)

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

	if !isTimeout {
		fmt.Printf("任务:%s, 正常输出:%s, 异常输出:%s, 错误:%s, 执行时间:%d", j.name, cmdOut, cmdErr, err, int(consume))
	} else {
		fmt.Printf("任务:%s, 执行超过了%d秒, 异常输出:%s", j.name, timeout/time.Second, cmdErr)
	}
}

func RememberCron(data []byte) error {
	path, _ := os.Getwd()
	return ioutil.WriteFile(path+"/data/cron.data", data, 0644)
}

func GetLocalCron() ([]byte, error) {
	path, _ := os.Getwd()
	return ioutil.ReadFile(path + "/data/cron.data")
}

func RunLocalTask() error {
	data, err := GetLocalCron()
	if err != nil {
		return err
	}
	return RunTask(data, false)
}
