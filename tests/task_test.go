package tests

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
	"webcron-agent/jobs"
	"webcron-agent/libs"

	"github.com/davecgh/go-spew/spew"
)

func init() {
}

func TestDataDecode(t *testing.T) {
	data := "[{\"id\":1,\"user_id\":1,\"server_id\":1,\"task_name\":\"测试脚本\",\"task_type\":1,\"description\":\"这是我的测试脚本描述\",\"cron_spec\":\"* * * * * *\",\"concurrent\":0,\"command\":\"ls /tmp\",\"status\":1,\"notify\":1,\"notify_email\":\"zhanghuitao@lianjia.com\",\"timeout\":3,\"execute_times\":3,\"prev_time\":1472720249,\"create_time\":1472720309}]"
	task_list, _ := jobs.DecodeTask([]byte(data))
	for _, task := range task_list {
		fmt.Println(task.TaskName)
	}
}

func TestStopTask(t *testing.T) {
	ids := [...]int{1, 2, 3}
	d, _ := json.Marshal(ids)
	jobs.StopTasks(d)
}

func TestLocalTask(t *testing.T) {
	jobs.RunLocalTask()
	//time.Sleep(time.Second * 10)
}

func TestRemeberCron(t *testing.T) {
	data := "[{\"id\":1,\"user_id\":1,\"server_id\":1,\"task_name\":\"测试脚本\",\"task_type\":1,\"description\":\"这是我的测试脚本描述\",\"cron_spec\":\"* * * * * *\",\"concurrent\":0,\"command\":\"ls /tmp\",\"status\":1,\"notify\":1,\"notify_email\":\"zhanghuitao@lianjia.com\",\"timeout\":3,\"execute_times\":3,\"prev_time\":1472720249,\"create_time\":1472720309}]"
	task_list, _ := jobs.DecodeTask([]byte(data))
	jobs.RememberCron(task_list, true)
}

func TestCreateLog(t *testing.T) {
	libs.CreateLogFile("logs/task.log")
}

func TestTimeformat(t *testing.T) {
	now := time.Now()
	spew.Dump(now.Format("20060102"))
}
