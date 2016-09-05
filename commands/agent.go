package commands

import (
	"fmt"
	"webcron-agent/service"

	"github.com/urfave/cli"
)

func init() {
}

func startCron(c *cli.Context) error {
	agent := service.NewAgentService()
	return agent.Start()
}

func listCommands(c *cli.Context) error {
	//args := c.Args()
	fmt.Println("1.	* * * * * ls /tmp")
	fmt.Println("2.	* * * * * ls /tmp")
	fmt.Println("3.	* * * * * ls /tmp")
	return nil
}

func runCommand(c *cli.Context) error {
	args := c.Args()
	fmt.Println(args[0])
	return nil
}

func editCommand(c *cli.Context) error {
	args := c.Args()
	fmt.Println(args[0])
	return nil
}
