package main

import (
	"strings"
	"io/ioutil"
	"github.com/smallfish/simpleyaml"
	"fmt"
	"strconv"
	"path/filepath"
	"os"
)

type Configuration struct {
	GlobalUrl string
	GlobalHeaders map[string]bool
	GlobalOptions map[string]bool
	Endpoints map[string]*Endpoint
	Requests map[string]*Request
}

func NewConfiguration() *Configuration {
	configuration := &Configuration{
		GlobalHeaders: make(map[string]bool),
		GlobalOptions: make(map[string]bool),
		Endpoints: make(map[string]*Endpoint),
		Requests: make(map[string]*Request),
	}
	return configuration
}

func (conf *Configuration) LoadFile(file string) {
	if !strings.HasSuffix(file, ".yaml") {
		file = file + ".yaml"
	}
	conf.readConfiguration(file)
}

func (conf *Configuration) visit(path string, f os.FileInfo, err error) error {
	if !f.IsDir() {
		if strings.HasSuffix(path, ".yaml") {
			fmt.Printf("Loading: %s\n", path)
			conf.readConfiguration(path)
		}
	}
	return nil
}

func (conf *Configuration) LoadAll() {
	filepath.Walk(".", conf.visit)
}

func (conf *Configuration) readConfiguration(moduleDefinition string) {
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
		if conf.isConfiguration(key.(string)) {
			conf.addConfiguration(key.(string), yaml)
		}
	}

	conf.readEndpoints(asMap["endpoints"].(map[interface{}]interface{}), yaml)
	conf.readRequests(asMap["requests"].(map[interface{}]interface{}))
}

func (conf *Configuration) ShowAll() {
	for name := range conf.Requests {
		request := conf.Requests[name]
		request.show()
		fmt.Println("")
	}
}

func (conf *Configuration) ShowRequest(requestName string) {
	request := conf.Requests[requestName]
	if request != nil {
		request.show()
		return
	}

	endpoint := conf.Endpoints[requestName]
	if endpoint != nil {
		endpoint.show()
	}
}

func (conf *Configuration) readRequests(requestsMap map[interface{}]interface{}) {
	for key := range requestsMap {
		keyAsString := key.(string)
		conf.addRequest(keyAsString, requestsMap[key])
	}
}

func (conf *Configuration) readEndpoints(endpointMap map[interface{}]interface{}, yaml *simpleyaml.Yaml) {
	for key := range endpointMap {
		keyAsString := key.(string)
		conf.addEndpoint(keyAsString, yaml)
	}
}

func (conf *Configuration) addRequest(name string, value interface{}) {
	request := &Request{
		Name: name,
		Headers: make(map[string]bool),
		Options: make(map[string]bool),
	}
	conf.Requests[name] = request

	request.Parameters = value.(map[interface{}]interface{})

	endpointName := request.Parameters["endpoint"]
	endpoint := conf.Endpoints[endpointName.(string)]
	request.Method = endpoint.Method
	request.Url = endpoint.Url
	request.Path = endpoint.Path
	request.Query = endpoint.Query

	for k,v := range endpoint.Headers {
		request.Headers[k] = v
	}

	for k,v := range endpoint.Options {
		request.Options[k] = v
	}

	for k := range request.Parameters {
		toReplace := "{"+k.(string)+"}"
		replacement := conf.getReplacement(request.Parameters[k])
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
		for option := range request.Options {
			replaced := strings.Replace(option, toReplace, replacement, -1)
			if (option != replaced) {
				request.Options[replaced] = true
				delete(request.Options, option)
			}
		}
	}
}


func (conf *Configuration) getReplacement(value interface{}) string {
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


func (conf *Configuration) addEndpoint(name string, yaml *simpleyaml.Yaml) {
	endpoint := &Endpoint{
		Name: name,
		Headers: make(map[string]bool),
		Options: make(map[string]bool),
	};
	conf.Endpoints[name] = endpoint

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
		endpoint.Url = conf.GlobalUrl
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

	for globalHeader := range conf.GlobalHeaders {
		endpoint.Headers[globalHeader] = true
	}

	options, err := yaml.GetPath("endpoints", name, "options").Array()
	if err == nil {
		for i := range options {
			endpoint.Options[options[i].(string)] = true
		}
	}

	for globalOption := range conf.GlobalOptions {
		endpoint.Options[globalOption] = true
	}
}

func (conf *Configuration) addConfiguration(name string, yaml *simpleyaml.Yaml) {
	if name == "headers" {
		headers, _ := yaml.Get(name).Array()
		for i := range headers {
			conf.GlobalHeaders[headers[i].(string)] = true
		}
	} else if name == "url" {
		conf.GlobalUrl, _ = yaml.Get(name).String()
	} else if name == "options" {
		options, _ := yaml.Get(name).Array()
		for i := range options {
			conf.GlobalOptions[options[i].(string)] = true
		}
	}
}

func (conf *Configuration) isConfiguration(name string) bool {
	return name == "headers" || name == "url" || name == "options"
}

