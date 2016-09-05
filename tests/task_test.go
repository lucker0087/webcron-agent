package tests

import (
	"encoding/json"
	"fmt"
	"testing"
	"webcron-agent/jobs"

	"github.com/davecgh/go-spew/spew"
)

func init() {
}

func TestDataDecode(t *testing.T) {
	data := "[{\"id\":1,\"user_id\":1,\"server_id\":1,\"task_name\":\"测试脚本\",\"task_type\":1,\"description\":\"这是我的测试脚本描述\",\"cron_spec\":\"* * * * * *\",\"concurrent\":0,\"command\":\"ls /tmp\",\"status\":1,\"notify\":1,\"notify_email\":\"zhanghuitao@lianjia.com\",\"timeout\":3,\"execute_times\":3,\"prev_time\":1472720249,\"create_time\":1472720309}]"

	var task_list []*jobs.Task
	err := json.Unmarshal([]byte(data), &task_list)
	if err != nil {
		fmt.Println(err.Error())
	}
	spew.Dump(task_list)
}
