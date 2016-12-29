package main

import (
	"os"
	"github.com/urfave/cli"
	"fmt"
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
	conf := NewDefaultConfiguration();
	executor := &Executor{conf: conf}
	app := cli.NewApp()
	app.Version = "0.1.0"
	printer := &Printer{conf: conf, writer: os.Stdout}

	var loadAllFiles bool
	var file string
	var directory string

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
		cli.StringFlag{
			Name:        "d",
			Usage:       "Specify directory to load the yaml files",
			Destination: &directory,
			Value:       ".",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:      "requests",
			ShortName: "r",
			Action: func(c *cli.Context) error {
				conf.Init(loadAllFiles, directory, file)
				printer.ShowRequests()
				return nil
			},
		},
		{
			Name:      "endpoints",
			ShortName: "e",
			Action: func(c *cli.Context) error {
				conf.Init(loadAllFiles, directory, file)
				printer.ShowEndpoints()
				return nil
			},
		},
		{
			Name: "show",
			Action: func(c *cli.Context) error {
				conf.Init(loadAllFiles, directory, file)
				printer.ShowRequestOrEndpoint(c.Args().First())
				return nil
			},
		},
		{
			Name: "run",
			Action: func(c *cli.Context) error {
				conf.Init(loadAllFiles, directory, file)
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

func (request *Request) String() string {
	return fmt.Sprintf("%v %v %v %v %v Headers=%v Options=%v Param=%v", request.Name,
		request.Method,
		request.Url,
		request.Path,
		request.Query,
		len(request.Headers),
		len(request.Options),
		len(request.Parameters),
	)
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%v %v %v %v %v Headers=%v Options=%v", endpoint.Name,
		endpoint.Method,
		endpoint.Url,
		endpoint.Path,
		endpoint.Query,
		len(endpoint.Headers),
		len(endpoint.Options),
	)
}
