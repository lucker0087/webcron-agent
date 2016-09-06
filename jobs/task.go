package jobs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

var l sync.RWMutex

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

func RunTask(data []byte, remeber bool) error {
	task_list, err := DecodeTask(data)

	if err != nil {
		return err
	}

	if remeber {
		go RememberCron(task_list)
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

	task_list, _ := GetLocalCron()

	for _, id := range ids {
		if len(task_list) > 0 {
			for _, task := range task_list {
				if task.Id == id {
					RemoveJob(id)
					delete(task_list, id)
					fmt.Printf("停止任务:%d\n", id)
				} else {
					fmt.Printf("任务不存在%d\n", id)
				}
			}
		}
	}

	return RememberCron(task_list)
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

func DecodeTask(d []byte) (map[int]*Task, error) {
	var task_list []*Task

	new_task_list := map[int]*Task{}

	err := json.Unmarshal(d, &task_list)
	if err != nil {
		return nil, err
	}

	if len(task_list) > 0 {
		for _, task := range task_list {
			new_task_list[task.Id] = task
		}
	}
	return new_task_list, nil
}

func EncodeTask(task_list map[int]*Task) ([]byte, error) {
	var list []*Task
	if len(task_list) == 0 {
		return nil, errors.New("Task list empty")
	}

	for _, task := range task_list {
		list = append(list, task)
	}
	return json.Marshal(list)
}

func RememberCron(task_list map[int]*Task) error {
	path, _ := os.Getwd()
	local_task_list, _ := GetLocalCron()

	l.Lock()
	defer l.Unlock()

	if len(local_task_list) > 0 {
		for id, task := range local_task_list {
			task_list[id] = task
		}
	}

	data, _ := EncodeTask(task_list)

	//return ioutil.WriteFile(filepath.Dir(path)+"/data/cron.data", data, 0644)
	return ioutil.WriteFile(path+"/data/cron.data", data, 0644)
}

func GetLocalCron() (map[int]*Task, error) {
	l.RLock()
	defer l.RUnlock()
	path, _ := os.Getwd()
	//data, err := ioutil.ReadFile(filepath.Dir(path) + "/data/cron.data")
	data, err := ioutil.ReadFile(path + "/data/cron.data")
	if err != nil {
		return nil, err
	}
	return DecodeTask(data)
}

func RunLocalTask() error {
	task_list, err := GetLocalCron()
	if err != nil {
		return err
	}
	data, err := EncodeTask(task_list)
	if err != nil {
		return err
	}
	return RunTask(data, false)
}
