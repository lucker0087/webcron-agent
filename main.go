package main

import (
	"os"
	"webcron-agent/commands"

	"github.com/urfave/cli"
)

const VERSION = "1.0.0"

func main() {
	app := cli.NewApp()
	app.Name = "WebCronAgent"
	app.Version = VERSION
	app.Commands = append(app.Commands, commands.AgentCommands()...)
	app.Run(os.Args)
}
