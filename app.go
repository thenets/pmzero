package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	lib "github.com/thenets/pmzero/lib"
	"github.com/urfave/cli"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Special condition for "run" command
	// cause it can conflict with the command args.
	// if len(os.Args) > 1 {
	// 	if os.Args[1] == "run" {
	// 		lib.createProcess(os.Args[2], os.Args[3:])
	// 		os.Exit(0)
	// 	}
	// }

	// CLI
	app := cli.NewApp()
	app.Name = "pmzero"
	app.Usage = "Easiest multi-platform process manager"
	app.Version = "0.0.3-alpha"
	app.Commands = []cli.Command{
		{
			Name:    "load",
			Aliases: []string{"l"},
			Usage:   "load a config file",
			Action: func(c *cli.Context) error {
				if len(c.Args().First()) == 0 {
					log.Fatalf("[ERROR] Required the config file path.\nExample: %v load <configFilePath>\n", app.Name)
				}

				// Load deployment file if it's a deployment type
				var data = lib.GetDeploymentByFilePath(c.Args().First())
				if data.Type == "deployment" {
					lib.AddDeploymentFile(c.Args().First())
				} else {
					log.Fatalf("[ERROR] config file type not supported: %v\n", data.Type)
				}

				return nil
			},
		},
		{
			Name:    "start",
			Aliases: []string{"c"},
			Usage:   "start a deployment",
			Action: func(c *cli.Context) error {
				if len(c.Args().First()) == 0 {
					log.Fatalf("[ERROR] require a deployment name.\nExample: %v stop <deploymentName>\n", app.Name)
				}
				if lib.HasDeployment(c.Args().First()) {
					lib.StartDeployment(c.Args().First())
				} else {
					log.Fatalf("[ERROR] deployment '%s' not found", c.Args().First())
				}
				return nil
			},
		},
		{
			Name:  "list",
			Usage: "list all deployments",
			Action: func(c *cli.Context) error {
				listDeployments()

				return nil
			},
		},
		{
			Name:  "delete",
			Usage: "delete a deployment",
			Action: func(c *cli.Context) error {
				if len(c.Args().First()) == 0 {
					log.Fatalf("[ERROR] require a deployment name.\nExample: %v delete <deploymentName>\n", app.Name)
				}
				if lib.HasDeployment(c.Args().First()) {
					lib.StopDeployment(c.Args().First())
				} else {
					log.Fatalf("[ERROR] deployment '%s' not found", c.Args().First())
				}

				lib.DeleteDeployment(c.Args().First())
				return nil
			},
		},
		{
			Name:  "stop",
			Usage: "stop a deployment",
			Action: func(c *cli.Context) error {
				if len(c.Args().First()) == 0 {
					log.Fatalf("[ERROR] require a deployment name.\nExample: %v stop <deploymentName>\n", app.Name)
				}
				if lib.HasDeployment(c.Args().First()) {
					lib.StopDeployment(c.Args().First())
				} else {
					log.Fatalf("[ERROR] deployment '%s' not found", c.Args().First())
				}

				return nil
			},
		},
		{
			Name:  "logs",
			Usage: "follow the logs files from a deployment",
			Action: func(c *cli.Context) error {
				if len(c.Args().First()) == 0 {
					log.Fatalf("[ERROR] require a deployment name.\nExample: %v logs <deploymentName>\n", app.Name)
				}
				if lib.HasDeployment(c.Args().First()) {
					lib.TailDeployment(c.Args().First())
				} else {
					log.Fatalf("[ERROR] deployment '%s' not found", c.Args().First())
				}

				return nil
			},
		},
		{
			Name:  "foreground",
			Usage: "keep running and respawn all deployments",
			Action: func(c *cli.Context) error {
				fmt.Printf("Press Ctrl+C to end\n")
				lib.ForegroundDeployments()

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	lib.UpdateState()

}

func listDeployments() {
	data := [][]string{}

	var deployments = lib.GetDeployments()

	for _, d := range deployments {
		data = append(data, []string{d.Name, d.Status, strconv.Itoa(d.PID), strings.Join(d.CMD[:], " ")})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Status", "PID", "CMD"})

	for _, v := range data {
		table.Append(v)
	}

	table.Render()
}
