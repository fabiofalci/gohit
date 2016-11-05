package main

import (
	"os"
	"io/ioutil"
	"text/template"
	"strings"
	"strconv"
	"fmt"

	"github.com/smallfish/simpleyaml"
	"github.com/urfave/cli"
	"bytes"
	"regexp"
	"bufio"
)

var (
	globalHeaders map[string]bool = make(map[string]bool)
	globalUrl string
	endpoints map[string]*Endpoint = make(map[string]*Endpoint)
	requests map[string]*Request = make(map[string]*Request)
)

const curlTemplate = `curl {{.Url}}{{.Path}}{{if .Query}}?{{.Query}}{{end}} \
{{- if .Headers}}
	{{- range $key, $value := .Headers }}
        -H '{{$key}}' \
	{{- end}}
{{- end}}
	-X{{.Method}}
`

type Endpoint struct {
	Name string
	Url string
	Path string
	Query string
	Method string
	Headers map[string]bool
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
}

func main() {
	app := cli.NewApp()
	app.Version = "0.1.0"
	var file string

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "file",
			Usage: "yaml file",
			Destination: &file,
		},
	}

	app.Commands = []cli.Command{
		{
			Name: "show-all",
			Action: func(c *cli.Context) error {
				loadAll(file)
				showAll()
				return nil
			},
		},
		{
			Name: "show",
			Action: func(c *cli.Context) error {
				loadAll(file)
				showRequest(c.Args().First())
				return nil
			},
		},
		{
			Name: "run",
			Action: func(c *cli.Context) error {
				loadAll(file)
				runRequest(c.Args().First())
				return nil
			},
		},
	}

	app.Run(os.Args)
}

func runRequest(requestName string) {
	request := requests[requestName]
	if request != nil {
		request.run()
	}

	endpoint := endpoints[requestName]
	if endpoint != nil {
		endpoint.run()
	}
}

func showRequest(requestName string) {
	request := requests[requestName]
	if request != nil {
		request.show()
		return
	}

	endpoint := endpoints[requestName]
	if endpoint != nil {
		endpoint.show()
	}
}

func loadAll(file string) {
	if !strings.HasSuffix(file, ".yaml") {
		file = file + ".yaml"
	}
	readConfiguration(file)
}

func readConfiguration(moduleDefinition string) {
	source, err := ioutil.ReadFile(moduleDefinition)
	if err != nil {
		panic(err)
	}

	yaml, err := simpleyaml.NewYaml(source)
	if err != nil {
		panic(err)
	}
	asMap, err := yaml.Map()
	if err != nil {
		panic(err)
	}

	for key := range asMap {
		if isConfiguration(key.(string)) {
			addConfiguration(key.(string), yaml)
		}
	}

	readEndpoints(asMap["endpoints"].(map[interface{}]interface{}), yaml)
	readRequests(asMap["requests"].(map[interface{}]interface{}))
}

func showAll() {
	for name := range requests {
		request := requests[name]
		request.show()
		fmt.Println("")
	}
}
func readRequests(requestsMap map[interface{}]interface{}) {
	for key := range requestsMap {
		keyAsString := key.(string)
		addRequest(keyAsString, requestsMap[key])
	}
}

func readEndpoints(endpointMap map[interface{}]interface{}, yaml *simpleyaml.Yaml) {
	for key := range endpointMap {
		keyAsString := key.(string)
		addEndpoint(keyAsString, yaml)
	}
}

func addRequest(name string, value interface{}) {
	request := &Request{Name: name, Headers: make(map[string]bool)}
	requests[name] = request

	request.Parameters = value.(map[interface{}]interface{})

	endpointName := request.Parameters["endpoint"]
	endpoint := endpoints[endpointName.(string)]
	request.Method = endpoint.Method
	request.Url = endpoint.Url
	request.Path = endpoint.Path
	request.Query = endpoint.Query

	for k,v := range endpoint.Headers {
		request.Headers[k] = v
	}

	for k := range request.Parameters {
		toReplace := "{"+k.(string)+"}"
		replacement := getReplacement(request.Parameters[k])
		request.Url = strings.Replace(request.Url, toReplace, replacement, -1)
		request.Path = strings.Replace(request.Path, toReplace, replacement, -1)
		request.Query = strings.Replace(request.Query, toReplace, replacement, -1)

		for header := range request.Headers {
			replaced := strings.Replace(header, toReplace, replacement, -1)
			if (header != replaced) {
				request.Headers[replaced] = true
				delete(request.Headers, header)
			}
		}
	}
}

func getReplacement(value interface{}) string {
	switch value.(type) {
	default:
		return "<<Error: invalid tpe>>"
	case bool:
		return strconv.FormatBool(value.(bool))
	case int:
		return strconv.Itoa(value.(int))
	case string:
		return value.(string)
	}
}


func addEndpoint(name string, yaml *simpleyaml.Yaml) {
	endpoint := &Endpoint{Name: name, Headers: make(map[string]bool)};
	endpoints[name] = endpoint

	path, err := yaml.GetPath("endpoints", name, "path").String()
	if err != nil {
		panic(err)
	}
	endpoint.Path = path

	query, err := yaml.GetPath("endpoints", name, "query").String()
	if err == nil {
		endpoint.Query = query
	}

	url, err := yaml.GetPath("endpoints", name, "url").String()
	if err == nil {
		endpoint.Url = url
	} else {
		endpoint.Url = globalUrl
	}

	method, err := yaml.GetPath("endpoints", name, "method").String()
	if err == nil {
		endpoint.Method = method
	} else {
		endpoint.Method = "GET"
	}

	data, err := yaml.GetPath("endpoints", name, "data").Bool()
	if err == nil {
		endpoint.HasData = data
	}

	headers, err := yaml.GetPath("endpoints", name, "headers").Array()
	if err == nil {
		for i := range headers {
			endpoint.Headers[headers[i].(string)] = true
		}
	}

	for globalHeader := range globalHeaders {
		endpoint.Headers[globalHeader] = true
	}
}

func (endpoint Endpoint) show() {
	t := template.Must(template.New("curlTemplate").Parse(curlTemplate))
	t.Execute(os.Stdout, endpoint)
}

func (request Request) show() {
	t := template.Must(template.New("curlTemplate").Parse(curlTemplate))
	fmt.Printf("Request %v:\n", request.Name)
	t.Execute(os.Stdout, request)
}

func (request Request) run() {
	t := template.Must(template.New("curlTemplate").Parse(curlTemplate))
	fmt.Printf("Request %v:\n", request.Name)
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

	println(requestAsString)
}

func (endpoint Endpoint) run() {
	t := template.Must(template.New("curlTemplate").Parse(curlTemplate))
	fmt.Printf("Request %v:\n", endpoint.Name)
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

	println(endpointAsString)
}

func hasResolvedAllVariables(request string) bool {
	return strings.Index(request, "{") == -1
}

func addConfiguration(name string, yaml *simpleyaml.Yaml) {
	if name == "headers" {
		headers, _ := yaml.Get(name).Array()
		for i := range headers {
			globalHeaders[headers[i].(string)] = true
		}
	} else if name == "url" {
		globalUrl, _ = yaml.Get(name).String()
	}
}

func isConfiguration(name string) bool {
	return name == "headers" || name == "url"
}
