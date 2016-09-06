package commands

import (
	"github.com/urfave/cli"
)

func AgentCommands() []cli.Command {
	return []cli.Command{
		{
			Name:  "agent",
			Usage: "web cron agent commands",
			Subcommands: []cli.Command{
				{
					Name:   "start",
					Usage:  "start agent",
					Action: startCron,
				},
				{
					Name:   "list",
					Usage:  "List cron commands",
					Action: listCron,
				},
				{
					Name:   "stop",
					Usage:  "stop the agent server",
					Action: stopCron,
				},
			},
		},
	}
}
