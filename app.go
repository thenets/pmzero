package main

import (
	"fmt"
	"log"
	"os"

	lib "github.com/thenets/pmzero/lib"
	"github.com/urfave/cli"
)

func main() {
	// Special condition for "start" command
	// cause it can conflit with the command args.
	if len(os.Args) > 1 {
		if os.Args[1] == "start" {
			lib.CreateProcess(os.Args[2], os.Args[3:])
			os.Exit(0)
		}
	}

	// CLI
	app := cli.NewApp()
	app.Name = "pmzero"
	app.Usage = "Windows and Linux process manager."
	app.Version = "0.0.1-alpha"
	app.Commands = []cli.Command{
		{
			Name:    "load",
			Aliases: []string{"l"},
			Usage:   "load a config file.",
			Action: func(c *cli.Context) error {
				if len(c.Args().First()) == 0 {
					log.Fatalf("[ERROR] Required the config file path.\nExample: %v load <configFilePath>\n", app.Name)
				}

				var configData = lib.ReadConfigFile(c.Args().First())
				fmt.Printf("%#v\n", configData)

				return nil
			},
		},
		{
			Name:    "start",
			Aliases: []string{"c"},
			Usage:   "start a process.",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:    "print",
			Aliases: []string{"c"},
			Usage:   "DEBUG prints the args.",
			Action: func(c *cli.Context) error {
				fmt.Println(c.Args()[1:])
				return nil
			},
		},
		{
			Name:    "template",
			Aliases: []string{"t"},
			Usage:   "options for task templates",
			Subcommands: []cli.Command{
				{
					Name:  "add",
					Usage: "add a new template",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
				{
					Name:  "remove",
					Usage: "remove an existing template",
					Action: func(c *cli.Context) error {
						fmt.Println("removed task template: ", c.Args().First())
						return nil
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
