package main

import (
	"fmt"
	"github.com/smallfish/simpleyaml"
	"io/ioutil"
	"strconv"
	"strings"
	"errors"
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
	Read() error
	Directory() string
	Configuration() map[string][]byte
}

func NewConfiguration(confReader ConfReader) (*Configuration, error) {
	configuration := &Configuration{
		GlobalHeaders:         make(map[string]bool),
		GlobalOptions:         make(map[string]bool),
		GlobalVariables:       make(map[string]interface{}),
		Endpoints:             make(map[string]*Endpoint),
		Requests:              make(map[string]*Request),
		requestsConfiguration: make(map[string]map[interface{}]interface{}),
		reader:                confReader,
	}
	if err := configuration.init(); err != nil {
		return nil, err
	}

	return configuration, nil
}

func (conf *Configuration) init() error {
	if err := conf.reader.Read(); err != nil {
		return err
	}
	for name, content := range conf.reader.Configuration() {
		if err := conf.readConfiguration(name, content); err != nil {
			return err
		}
	}
	conf.loadEndpointGlobals()
	if err := conf.loadRequests(); err != nil {
		return err
	}

	if err := conf.validate(); err != nil {
		return err
	}

	return nil
}

func (conf *Configuration) validate() error {
	if len(conf.Endpoints) == 0 {
		return errors.New("Missing endpoints")
	}
	for _, endpoint := range conf.Endpoints {
		if endpoint.Url == "" {
			return errors.New("Missing URL")
		}
		// need to validate just the first one
		break
	}

	for _, endpoint := range conf.Endpoints {
		if endpoint.Path == "" {
			return errors.New(fmt.Sprintf("Endpoint %v missing path", endpoint.Name))
		}
	}

	return nil
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

func (conf *Configuration) loadRequests() error {
	for k := range conf.requestsConfiguration {
		if err := conf.readRequests(conf.requestsConfiguration[k]); err != nil {
			return err
		}
	}
	return nil
}

func (conf *Configuration) readConfiguration(moduleDefinition string, source []byte) error {
	yaml, err := simpleyaml.NewYaml(source)
	if err != nil {
		return err
	}
	asMap, err := yaml.Map()
	if err != nil {
		return err
	}

	for key := range asMap {
		if conf.isConfiguration(key.(string)) {
			if err := conf.addConfiguration(key.(string), yaml); err != nil {
				return err
			}
		} else if key != "endpoints" && key != "requests" {
			return errors.New(fmt.Sprintf("Invalid yaml attribute '%v'", key))
		}
	}

	endpointsMap := asMap["endpoints"]
	if endpointsMap != nil {
		if err := conf.readEndpoints(endpointsMap.(map[interface{}]interface{}), yaml); err != nil {
			return err
		}
	}
	requestsMap := asMap["requests"]
	if requestsMap != nil {
		conf.requestsConfiguration[moduleDefinition] = requestsMap.(map[interface{}]interface{})
	}

	return nil
}

func (conf *Configuration) readRequests(requestsMap map[interface{}]interface{}) error {
	var err error
	for key := range requestsMap {
		keyAsString := key.(string)
		conf.Requests[keyAsString], err = conf.createRequest(keyAsString, requestsMap[key])
		if err != nil {
			return err
		}
	}
	return nil
}

func (conf *Configuration) readEndpoints(endpointMap map[interface{}]interface{}, yaml *simpleyaml.Yaml) error {
	for key := range endpointMap {
		keyAsString := key.(string)
		if err := conf.addEndpoint(keyAsString, yaml); err != nil {
			return err
		}
	}
	return nil
}

func (conf *Configuration) createRequest(name string, value interface{}) (*Request, error) {
	request := &Request{
		Name:    name,
		Headers: make(map[string]bool),
		Options: make(map[string]bool),
	}

	request.Parameters = value.(map[interface{}]interface{})

	endpointName := request.Parameters["endpoint"].(string)
	endpoint := conf.Endpoints[endpointName]
	if endpoint == nil {
		return nil, errors.New(fmt.Sprintf("Request %v couldn't find endpoint %v", name, endpointName))
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
	return request, nil
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

func (conf *Configuration) addEndpoint(name string, yaml *simpleyaml.Yaml) error {
	endpoint := &Endpoint{
		Name:    name,
		Headers: make(map[string]bool),
		Options: make(map[string]bool),
	}
	conf.Endpoints[name] = endpoint

	if path, err := yaml.GetPath("endpoints", name, "path").String(); err == nil {
		endpoint.Path = path
	}

	if query, err := yaml.GetPath("endpoints", name, "query").String(); err == nil {
		endpoint.Query = query
	}

	if url, err := yaml.GetPath("endpoints", name, "url").String(); err == nil {
		endpoint.Url = url
	} else {
		endpoint.Url = conf.GlobalUrl
	}

	if method, err := yaml.GetPath("endpoints", name, "method").String(); err == nil {
		endpoint.Method = method
	} else {
		endpoint.Method = "GET"
	}

	if headers, err := yaml.GetPath("endpoints", name, "headers").Array(); err == nil {
		for i := range headers {
			endpoint.Headers[headers[i].(string)] = true
		}
	}

	if options, err := yaml.GetPath("endpoints", name, "options").Array(); err == nil {
		for i := range options {
			endpoint.Options[options[i].(string)] = true
		}
	}

	return nil
}

func (conf *Configuration) addConfiguration(name string, yaml *simpleyaml.Yaml) error {
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
			if err := conf.readConfiguration(fileName, source); err != nil {
				return err
			}
		}
	} else if name == "variables" {
		variables, _ := yaml.Get(name).Map()
		for i := range variables {
			conf.GlobalVariables[i.(string)] = variables[i]
		}
	}

	return nil
}

func (conf *Configuration) isConfiguration(name string) bool {
	return name == "headers" || name == "url" || name == "options" || name == "files" || name == "variables"
}
