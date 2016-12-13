package main

import (
	"os"
	"github.com/urfave/cli"
)

type Endpoint struct {
	Name    string
	Url     string
	Path    string
	Query   string
	Method  string
	Headers map[string]bool
	Options map[string]bool
}

type Request struct {
	Name       string
	Url        string
	Path       string
	Query      string
	Method     string
	Headers    map[string]bool
	Options    map[string]bool
	Parameters map[interface{}]interface{}
}

type Executable interface {
	GetName() string
	GetOptions() map[string]bool
}

func main() {
	conf := NewConfiguration()
	executor := &Executor{conf: conf}
	app := cli.NewApp()
	app.Version = "0.1.0"
	printer := NewPrinter(conf)

	var loadAllFiles bool
	var file string

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "f",
			Usage:       "Load one yaml file",
			Destination: &file,
		},
		cli.BoolFlag{
			Name:        "r",
			Usage:       "Recursively load all yaml files",
			Destination: &loadAllFiles,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:      "requests",
			ShortName: "r",
			Action: func(c *cli.Context) error {
				conf.Init(loadAllFiles, file)
				printer.ShowRequests()
				return nil
			},
		},
		{
			Name:      "endpoints",
			ShortName: "e",
			Action: func(c *cli.Context) error {
				conf.Init(loadAllFiles, file)
				printer.ShowEndpoints()
				return nil
			},
		},
		{
			Name: "show",
			Action: func(c *cli.Context) error {
				conf.Init(loadAllFiles, file)
				printer.ShowRequestOrEndpoint(c.Args().First())
				return nil
			},
		},
		{
			Name: "run",
			Action: func(c *cli.Context) error {
				conf.Init(loadAllFiles, file)
				executor.RunRequest(c.Args().First())
				return nil
			},
		},
	}

	app.Run(os.Args)
}

func (endpoint *Endpoint) GetName() string {
	return endpoint.Name
}

func (request *Request) GetName() string {
	return request.Name
}

func (endpoint *Endpoint) GetOptions() map[string]bool {
	return endpoint.Options
}

func (request *Request) GetOptions() map[string]bool {
	return request.Options
}
