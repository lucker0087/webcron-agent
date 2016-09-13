package jobs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
	"webcron-agent/config"
	"webcron-agent/libs"
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

func CronTask(data []byte, remeber bool) error {
	task_list, err := DecodeTask(data)

	if err != nil {
		return err
	}

	if remeber {
		go RememberCron(task_list, true)
	}

	logger := libs.NewTaskLogger()

	for _, task := range task_list {
		job, err := NewJob(task)
		if err != nil {
			logger.Warning(fmt.Sprintf("InitJobs error :%s", err.Error()))
			continue
		}
		AddJob(task.CronSpec, job)
	}
	return nil
}

func RunTask(data []byte) error {
	task_list, err := DecodeTask(data)

	if err != nil {
		return err
	}

	logger := libs.NewTaskLogger()
	for _, task := range task_list {
		job, err := NewJob(task)
		if err != nil {
			logger.Warning(fmt.Sprintf("InitJobs error :%s", err.Error()))
			continue
		}
		job.Run()
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

	logger := libs.NewTaskLogger()
	task_list, _ := GetLocalCron()
	for _, id := range ids {
		if len(task_list) > 0 {
			for _, task := range task_list {
				if task.Id == id {
					RemoveJob(id)
					delete(task_list, id)
					logger.Info(fmt.Sprintf("停止任务:%d\n", id))
				} else {
					logger.Warning(fmt.Sprintf("任务不存在%d\n", id))
				}
			}
		}
	}

	return RememberCron(task_list, false)
}

func NewJob(task *Task) (*Job, error) {
	if task.Id < 1 {
		return nil, fmt.Errorf("ToJob: 缺少id")
	}
	job := &Job{
		id:         task.Id,
		name:       task.TaskName,
		Concurrent: task.Concurrent == 1,
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

func RememberCron(task_list map[int]*Task, append bool) error {

	var local_task_list = make(map[int]*Task, 0)

	if append {
		local_task_list, _ = GetLocalCron()
	}

	l.Lock()
	defer l.Unlock()

	if len(local_task_list) > 0 {
		for id, task := range local_task_list {
			task_list[id] = task
		}
	}

	data, _ := EncodeTask(task_list)

	return ioutil.WriteFile(GetCronDataPath(), data, 0644)
}

func GetLocalCron() (map[int]*Task, error) {
	l.RLock()
	defer l.RUnlock()

	data, err := ioutil.ReadFile(GetCronDataPath())
	if err != nil {
		return nil, err
	}

	if len(data) > 0 {
		return DecodeTask(data)
	}

	return nil, nil
}

func GetCronDataPath() string {
	config, _ := config.GetConfig()
	return filepath.Join(config.App.Path, config.Cron.DataPath)
}

func RunLocalTask() error {
	task_list, err := GetLocalCron()
	if err != nil {
		return err
	}
	data, err := EncodeTask(task_list)
	if err != nil || len(data) == 0 {
		return err
	}
	return CronTask(data, false)
}
