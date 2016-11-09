package main

import (
	"os"
	"text/template"
	"strings"
	"fmt"

	"github.com/urfave/cli"
	"bytes"
	"regexp"
	"bufio"
	"os/exec"
)

const showCurlTemplate = `curl {{.Url}}{{.Path}}{{if .Query}}?{{.Query}}{{end}} \
{{- if .Headers}}
	{{- range $key, $value := .Headers }}
        -H '{{$key}}' \
	{{- end}}
{{- end}}
{{- if .Options}}
	{{- range $key, $value := .Options }}
        {{$key}} \
	{{- end}}
{{- end}}
    -X{{.Method}}
`

// A separated template for running as it needs to transform the command to an array fo string.
// It splits on newline.
const runCurlTemplate = `{{.Url}}{{.Path}}{{if .Query}}?{{.Query}}{{end}}
{{- if .Headers}}
	{{- range $key, $value := .Headers }}
-H
{{$key}}
	{{- end}}
{{- end}}
{{- if .Options}}
{{.OptionsAsToken}}
{{- end}}
-X{{.Method}}`

type Endpoint struct {
	Name string
	Url string
	Path string
	Query string
	Method string
	Headers map[string]bool
	Options map[string]bool
	HasData bool
}

type Request struct {
	Name string
	Parameters map[interface{}]interface{}

	Url string
	Path string
	Query string
	Method string
	Headers map[string]bool
	Options map[string]bool
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

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "f",
			Usage: "Load one yaml file",
			Destination: &file,
		},
		cli.BoolFlag{
			Name: "r",
			Usage: "Recursively load all yaml files",
			Destination: &loadAllFiles,
		},
	}

	app.Commands = []cli.Command{
		{
			Name: "show-all",
			Action: func(c *cli.Context) error {
				if loadAllFiles {
					conf.LoadAll()
				} else {
					conf.LoadFile(file)
				}
				conf.ShowAll()
				return nil
			},
		},
		{
			Name: "show",
			Action: func(c *cli.Context) error {
				if loadAllFiles {
					conf.LoadAll()
				} else {
					conf.LoadFile(file)
				}
				conf.ShowRequest(c.Args().First())
				return nil
			},
		},
		{
			Name: "run",
			Action: func(c *cli.Context) error {
				if loadAllFiles {
					conf.LoadAll()
				} else {
					conf.LoadFile(file)
				}
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
		request.run()
	}

	endpoint := executor.Conf.Endpoints[requestName]
	if endpoint != nil {
		endpoint.run()
	}
}



func (endpoint Endpoint) show() {
	t := template.Must(template.New("curlTemplate").Parse(showCurlTemplate))
	t.Execute(os.Stdout, endpoint)
}

func (request Request) show() {
	t := template.Must(template.New("curlTemplate").Parse(showCurlTemplate))
	fmt.Printf("Request %v:\n", request.Name)
	t.Execute(os.Stdout, request)
}

func (request Request) run() {
	t := template.Must(template.New("curlTemplate").Parse(runCurlTemplate))
	fmt.Printf("%v:\n", request.Name)
	buf := new(bytes.Buffer)
	t.Execute(buf, request)

	requestAsString := buf.String()

	if (!hasResolvedAllVariables(requestAsString)) {
		re := regexp.MustCompile("{(.+?)}")
		for _, v  := range re.FindAllString(requestAsString, -1) {
			reader := bufio.NewReader(os.Stdin)
			fmt.Printf("Enter %v: ", v)
			value, _ := reader.ReadString('\n')
			value = strings.TrimSpace(value)
			requestAsString = strings.Replace(requestAsString, v, value, -1)
		}
	}

	executeCurlCommand(requestAsString)
}

func (endpoint Endpoint) run() {
	t := template.Must(template.New("curlTemplate").Parse(runCurlTemplate))
	fmt.Printf("%v:\n", endpoint.Name)
	buf := new(bytes.Buffer)
	t.Execute(buf, endpoint)

	endpointAsString := buf.String()

	if (!hasResolvedAllVariables(endpointAsString)) {
		re := regexp.MustCompile("{(.+?)}")
		for _, v  := range re.FindAllString(endpointAsString, -1) {
			reader := bufio.NewReader(os.Stdin)
			fmt.Printf("Enter %v: ", v)
			value, _ := reader.ReadString('\n')
			value = strings.TrimSpace(value)
			endpointAsString = strings.Replace(endpointAsString, v, value, -1)
		}
	}

	executeCurlCommand(endpointAsString)
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

func (request Request) OptionsAsToken() string {
	oneLineOptions := ""
	for option := range request.Options {
		re := regexp.MustCompile("[^\\s\"']+|\"([^\"]*)\"|'([^']*)'")
		for _, v  := range re.FindAllString(option, -1) {
			oneLineOptions = oneLineOptions + "\n" + v
		}
	}

	return strings.TrimPrefix(oneLineOptions, "\n")
}

