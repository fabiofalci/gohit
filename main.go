package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
	"runtime"
	"strconv"
	"time"
)

type Endpoint struct {
	Name       string
	Url        string
	Path       string
	QueryRaw   string
	QueryList  map[string]string
	Method     string
	Headers    map[string]bool
	Options    map[string]bool
	Parameters map[interface{}]interface{}
}

type Request struct {
	Name       string
	Url        string
	Path       string
	QueryRaw   string
	QueryList  map[string]string
	Method     string
	Headers    map[string]bool
	Options    map[string]bool
	Parameters map[interface{}]interface{}
}

type Executable interface {
	GetName() string
	GetOptions() map[string]bool
}

var version string
var commit string
var buildDate string

func main() {
	app := cli.NewApp()
	app.Version = version
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println("Version: " + c.App.Version)
		fmt.Println("Git commit: " + commit)
		if i, err := strconv.ParseInt(buildDate, 10, 64); err == nil {
			fmt.Println("Build date: " + time.Unix(i, 0).UTC().String())
		}
		fmt.Println("Go version: " + runtime.Version())
	}

	var file string
	var directory string
	var oneLine bool

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "f",
			Usage:       "Load one yaml file",
			Destination: &file,
		},
		cli.StringFlag{
			Name:        "d",
			Usage:       "Specify directory to load the yaml files",
			Destination: &directory,
			Value:       ".",
		},
		cli.BoolFlag{
			Name:        "oneline, ol",
			Usage:       "Print commands in one line",
			Destination: &oneLine,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:      "requests",
			ShortName: "r",
			Action: func(c *cli.Context) error {
				conf, err := NewConfiguration(NewDefaultConfigurationReader(directory, file))
				if err != nil {
					return err
				}
				printer := &Printer{conf: conf, writer: os.Stdout, oneLine: oneLine}
				printer.ShowRequests()
				return nil
			},
		},
		{
			Name:      "endpoints",
			ShortName: "e",
			Action: func(c *cli.Context) error {
				conf, err := NewConfiguration(NewDefaultConfigurationReader(directory, file))
				if err != nil {
					return err
				}
				printer := &Printer{conf: conf, writer: os.Stdout, oneLine: oneLine}
				printer.ShowEndpoints()
				return nil
			},
		},
		{
			Name: "show",
			Action: func(c *cli.Context) error {
				conf, err := NewConfiguration(NewDefaultConfigurationReader(directory, file))
				if err != nil {
					return err
				}
				printer := &Printer{conf: conf, writer: os.Stdout, oneLine: oneLine}
				printer.ShowRequestOrEndpoint(c.Args().First())
				return nil
			},
		},
		{
			Name: "run",
			Action: func(c *cli.Context) error {
				conf, err := NewConfiguration(NewDefaultConfigurationReader(directory, file))
				if err != nil {
					return err
				}
				executor := NewDefaultExecutor(conf)
				requestName := c.Args().First()
				return executor.RunRequest(requestName, c.Args().Tail())
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
	return fmt.Sprintf("%v %v %v %v QueryRaw=%v QueryList=%v Headers=%v Options=%v Param=%v", request.Name,
		request.Method,
		request.Url,
		request.Path,
		request.QueryRaw,
		request.QueryList,
		len(request.Headers),
		len(request.Options),
		len(request.Parameters),
	)
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%v %v %v %v QueryRaw=%v QueryList=%v Headers=%v Options=%v", endpoint.Name,
		endpoint.Method,
		endpoint.Url,
		endpoint.Path,
		endpoint.QueryRaw,
		endpoint.QueryList,
		len(endpoint.Headers),
		len(endpoint.Options),
	)
}
