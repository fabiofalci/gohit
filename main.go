package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"bufio"
	"bytes"
	"github.com/urfave/cli"
	"os/exec"
	"regexp"
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

type Executor struct {
	Conf *Configuration
}

func main() {
	conf := NewConfiguration()
	executor := &Executor{Conf: conf}
	app := cli.NewApp()
	app.Version = "0.1.0"
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
				if loadAllFiles {
					conf.LoadAll()
				} else {
					conf.LoadConfigurationAndEndpoints(file)
				}
				conf.LoadRequests()
				conf.ShowRequests()
				return nil
			},
		},
		{
			Name:      "endpoints",
			ShortName: "e",
			Action: func(c *cli.Context) error {
				if loadAllFiles {
					conf.LoadAll()
				} else {
					conf.LoadConfigurationAndEndpoints(file)
				}
				conf.LoadRequests()
				conf.ShowEndpoints()
				return nil
			},
		},
		{
			Name: "show",
			Action: func(c *cli.Context) error {
				if loadAllFiles {
					conf.LoadAll()
				} else {
					conf.LoadConfigurationAndEndpoints(file)
				}
				conf.LoadRequests()
				conf.ShowRequestOrEndpoint(c.Args().First())
				return nil
			},
		},
		{
			Name: "run",
			Action: func(c *cli.Context) error {
				if loadAllFiles {
					conf.LoadAll()
				} else {
					conf.LoadConfigurationAndEndpoints(file)
				}
				conf.LoadRequests()
				executor.RunRequest(c.Args().First())
				return nil
			},
		},
	}

	app.Run(os.Args)
}

func (executor *Executor) RunRequest(requestName string) {
	request := executor.Conf.Requests[requestName]
	if request != nil {
		runExecutable(request)
		return
	}

	endpoint := executor.Conf.Endpoints[requestName]
	if endpoint != nil {
		m := make(map[interface{}]interface{})
		m["endpoint"] = requestName
		request := executor.Conf.createRequest(requestName, m)
		runExecutable(request)
	}
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

func showExecutable(executable Executable) {
	t := template.Must(template.New("curlTemplate").Parse(showCurlTemplate))
	fmt.Printf("Endpoint %v:\n", executable.GetName())
	t.Execute(os.Stdout, executable)
}

func runExecutable(executable Executable) {
	t := template.Must(template.New("curlTemplate").Parse(runCurlTemplate))
	fmt.Printf("%v:\n", executable.GetName())
	buf := new(bytes.Buffer)
	t.Execute(buf, executable)

	requestAsString := buf.String()

	if !hasResolvedAllVariables(requestAsString) {
		re := regexp.MustCompile("{(.+?)}")
		for _, v := range re.FindAllString(requestAsString, -1) {
			reader := bufio.NewReader(os.Stdin)
			fmt.Printf("Enter %v: ", v)
			value, _ := reader.ReadString('\n')
			value = strings.TrimSpace(value)
			requestAsString = strings.Replace(requestAsString, v, value, -1)
		}
	}

	executeCurlCommand(requestAsString)
}

func executeCurlCommand(command string) {
	cmd := exec.Command("curl", strings.Split(command, "\n")...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	fmt.Println(out.String())
}

func hasResolvedAllVariables(request string) bool {
	return strings.Index(request, "{") == -1
}
