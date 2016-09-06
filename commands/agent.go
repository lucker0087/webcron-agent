package commands

import (
	"fmt"
	"webcron-agent/jobs"
	"webcron-agent/service"

	"github.com/urfave/cli"
)

func init() {
}

func startCron(c *cli.Context) error {
	agent := service.NewAgentService()
	return agent.Start()
}

func listCron(c *cli.Context) error {
	task_list, err := jobs.GetLocalCron()
	if err != nil {
		fmt.Printf("Get cron job error:%s\n", err.Error())
		return err
	}

	if len(task_list) == 0 {
		fmt.Println("No crontab running")
	}

	for _, task := range task_list {
		fmt.Printf("%d. %s %s %s\t%s", task.Id, task.TaskName, task.CronSpec, task.Command, task.Description)
	}
	return nil
}

func StopCron(c *cli.Context) error {
	//args := c.Args()
	//fmt.Println(args[0])
	return nil
}
