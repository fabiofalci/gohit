package main

import (
	"fmt"
	"github.com/smallfish/simpleyaml"
	"io/ioutil"
	"strconv"
	"strings"
)

type Configuration struct {
	GlobalUrl       string
	GlobalHeaders   map[string]bool
	GlobalOptions   map[string]bool
	GlobalVariables map[string]interface{}
	Endpoints       map[string]*Endpoint
	Requests        map[string]*Request

	requestsConfiguration map[string]map[interface{}]interface{}
	reader                ConfReader
}

type ConfReader interface {
	Read()
	Directory() string
	Configuration() map[string][]byte
}

func NewConfiguration(confReader ConfReader) *Configuration {
	configuration := &Configuration{
		GlobalHeaders:         make(map[string]bool),
		GlobalOptions:         make(map[string]bool),
		GlobalVariables:       make(map[string]interface{}),
		Endpoints:             make(map[string]*Endpoint),
		Requests:              make(map[string]*Request),
		requestsConfiguration: make(map[string]map[interface{}]interface{}),
		reader:                confReader,
	}
	configuration.init()
	return configuration
}

func (conf *Configuration) init() {
	conf.reader.Read()
	for name, content := range conf.reader.Configuration() {
		conf.readConfiguration(name, content)
	}
	conf.loadEndpointGlobals()
	conf.loadRequests()
}

func (conf *Configuration) loadEndpointGlobals() {
	for _, endpoint := range conf.Endpoints {
		for globalHeader := range conf.GlobalHeaders {
			endpoint.Headers[globalHeader] = true
		}

		for globalOption := range conf.GlobalOptions {
			endpoint.Options[globalOption] = true
		}
	}
}

func (conf *Configuration) loadRequests() {
	for k := range conf.requestsConfiguration {
		conf.readRequests(conf.requestsConfiguration[k])
	}
}

func (conf *Configuration) readConfiguration(moduleDefinition string, source []byte) {
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
		} else if key != "endpoints" && key != "requests" {
			panic(fmt.Sprintf("Invalid yaml attribute '%v'", key))
		}
	}

	endpointsMap := asMap["endpoints"]
	if endpointsMap != nil {
		conf.readEndpoints(endpointsMap.(map[interface{}]interface{}), yaml)
	}
	requestsMap := asMap["requests"]
	if requestsMap != nil {
		conf.requestsConfiguration[moduleDefinition] = requestsMap.(map[interface{}]interface{})
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
	conf.Requests[name] = conf.createRequest(name, value)
}

func (conf *Configuration) createRequest(name string, value interface{}) *Request {
	request := &Request{
		Name:    name,
		Headers: make(map[string]bool),
		Options: make(map[string]bool),
	}

	request.Parameters = value.(map[interface{}]interface{})

	endpointName := request.Parameters["endpoint"].(string)
	endpoint := conf.Endpoints[endpointName]
	if endpoint == nil {
		panic(fmt.Sprintf("Cannot find endpoint %v", endpointName))
	}
	request.Method = endpoint.Method
	request.Url = endpoint.Url
	request.Path = endpoint.Path
	request.Query = endpoint.Query

	for k, v := range endpoint.Headers {
		request.Headers[k] = v
	}

	for k, v := range endpoint.Options {
		request.Options[k] = v
	}

	for k := range request.Parameters {
		toReplace := "{" + k.(string) + "}"
		conf.replaceAll(request, toReplace, request.Parameters[k])
	}

	for k := range conf.GlobalVariables {
		toReplace := "{" + k + "}"
		conf.replaceAll(request, toReplace, conf.GlobalVariables[k])
	}
	return request
}

func (conf *Configuration) replaceAll(request *Request, toReplace string, value interface{}) {
	replacement := conf.getReplacement(value)
	request.Url = strings.Replace(request.Url, toReplace, replacement, -1)
	request.Path = strings.Replace(request.Path, toReplace, replacement, -1)
	request.Query = strings.Replace(request.Query, toReplace, replacement, -1)

	for header := range request.Headers {
		replaced := strings.Replace(header, toReplace, replacement, -1)
		if header != replaced {
			request.Headers[replaced] = true
			delete(request.Headers, header)
		}
	}
	for option := range request.Options {
		replaced := strings.Replace(option, toReplace, replacement, -1)
		if option != replaced {
			request.Options[replaced] = true
			delete(request.Options, option)
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
		Name:    name,
		Headers: make(map[string]bool),
		Options: make(map[string]bool),
	}
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

	headers, err := yaml.GetPath("endpoints", name, "headers").Array()
	if err == nil {
		for i := range headers {
			endpoint.Headers[headers[i].(string)] = true
		}
	}

	options, err := yaml.GetPath("endpoints", name, "options").Array()
	if err == nil {
		for i := range options {
			endpoint.Options[options[i].(string)] = true
		}
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
	} else if name == "files" {
		files, _ := yaml.Get(name).Array()
		for i := range files {
			fileName := files[i].(string)
			source, err := ioutil.ReadFile(conf.reader.Directory() + "/" + fileName)
			if err != nil {
				panic(err)
			}
			conf.readConfiguration(fileName, source)
		}
	} else if name == "variables" {
		variables, _ := yaml.Get(name).Map()
		for i := range variables {
			conf.GlobalVariables[i.(string)] = variables[i]
		}
	}
}

func (conf *Configuration) isConfiguration(name string) bool {
	return name == "headers" || name == "url" || name == "options" || name == "files" || name == "variables"
}
