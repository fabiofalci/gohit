package main

import (
	"fmt"
	"bytes"
	"regexp"
	"bufio"
	"os"
	"strings"
	"os/exec"
	"text/template"
	"errors"
)

type Executor struct {
	conf   *Configuration
	runner CommandRunner
}

type CommandRunner interface {
	Run(command []string) error
}

type DefaultRunner struct {
}

func NewDefaultExecutor(conf *Configuration) *Executor {
	runner := &DefaultRunner{}
	return NewExecutor(conf, runner)
}

func NewExecutor(conf *Configuration, runner CommandRunner) *Executor {
	executor := &Executor{
		conf: conf,
		runner: runner,
	}
	return executor
}

func (executor *Executor) RunRequest(requestName string) error {
	request := executor.conf.Requests[requestName]
	if request != nil {
		return executor.runExecutable(request)
	}

	endpoint := executor.conf.Endpoints[requestName]
	if endpoint != nil {
		if r, err := executor.createTemporaryRequest(requestName); err == nil {
			return executor.runExecutable(r)
		} else {
			return err
		}
	}
	return errors.New(fmt.Sprint("Could not find request/endpoint ", requestName))
}

func (executor *Executor) createTemporaryRequest(requestName string) (*Request, error) {
	m := make(map[interface{}]interface{}, 1)
	m["endpoint"] = requestName
	return executor.conf.createRequest(requestName, m)
}

func (executor *Executor) runExecutable(executable Executable) error {
	t := template.Must(template.New("curlTemplate").Parse(runCurlTemplate))
	buf := new(bytes.Buffer)
	t.Execute(buf, executable)

	requestAsString := buf.String()

	if !executor.hasResolvedAllVariables(requestAsString) {
		re := regexp.MustCompile("{(.+?)}")
		for _, v := range re.FindAllString(requestAsString, -1) {
			reader := bufio.NewReader(os.Stdin)
			fmt.Printf("Enter %v: ", v)
			value, _ := reader.ReadString('\n')
			value = strings.TrimSpace(value)
			requestAsString = strings.Replace(requestAsString, v, value, -1)
		}
	}

	asArray := strings.Split(requestAsString, "\n")
	return executor.runner.Run(asArray)
}

func (runner *DefaultRunner) Run(command []string) error {
	cmd := exec.Command("curl", command...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	fmt.Println(out.String())
	return nil
}

func (executor *Executor) hasResolvedAllVariables(request string) bool {
	return strings.Index(request, "{") == -1
}
