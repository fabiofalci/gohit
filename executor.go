package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/template"
)

type Executor struct {
	conf      *Configuration
	runner    CommandRunner
	varReader VariableReader
	jq        bool
}

type CommandRunner interface {
	Run(command []string, jq bool) error
}

type DefaultRunner struct {
}

type VariableReader interface {
	Read(variableName string) string
}

type DefaultVariableReader struct {
}

func NewDefaultExecutor(conf *Configuration, jq bool) *Executor {
	runner := &DefaultRunner{}
	varReader := &DefaultVariableReader{}
	return NewExecutor(conf, runner, varReader, jq)
}

func NewExecutor(conf *Configuration, runner CommandRunner, varReader VariableReader, jq bool) *Executor {
	executor := &Executor{
		conf:      conf,
		runner:    runner,
		varReader: varReader,
		jq:        jq,
	}
	return executor
}

func (executor *Executor) RunRequest(requestName string, args []string) error {
	request := executor.conf.Requests[requestName]
	if request != nil {
		return executor.runExecutable(request, args)
	}

	endpoint := executor.conf.Endpoints[requestName]
	if endpoint != nil {
		if r, err := executor.createTemporaryRequest(requestName); err == nil {
			return executor.runExecutable(r, args)
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

func (executor *Executor) runExecutable(executable Executable, args []string) error {
	t := template.Must(template.New("curlTemplate").Parse(runCurlTemplate))
	buf := new(bytes.Buffer)
	t.Execute(buf, executable)

	requestAsString := buf.String()
	if !executor.hasResolvedAllVariables(requestAsString) {
		re := regexp.MustCompile("{(.+?)}")
		for i, v := range re.FindAllString(requestAsString, -1) {
			value := executor.getValue(v, i, args)
			requestAsString = strings.Replace(requestAsString, v, value, -1)
		}
	}

	asArray := strings.Split(requestAsString, "\n")
	return executor.runner.Run(asArray, executor.jq)
}

func (executor *Executor) getValue(variableName string, position int, args []string) string {
	if len(args) > position {
		return args[position]
	}
	return executor.varReader.Read(variableName)
}

func (parameterReader *DefaultVariableReader) Read(variableName string) string {
	fmt.Printf("Enter %v: ", variableName)
	reader := bufio.NewReader(os.Stdin)
	value, _ := reader.ReadString('\n')
	return strings.TrimSpace(value)
}

func (runner *DefaultRunner) Run(command []string, jq bool) error {
	cmd := exec.Command("curl", command...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	var pipedCommand *exec.Cmd
	cmd.Stderr = &stderr

	if jq {
		pipedCommand = exec.Command("jq")
		pipedCommand.Stdin, _ = cmd.StdoutPipe()
		pipedCommand.Stdout = &out
		if err := pipedCommand.Start(); err != nil {
			return err
		}
	} else {
		cmd.Stdout = &out
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	if pipedCommand != nil {
		if err := pipedCommand.Wait(); err != nil {
			return err
		}
	}

	fmt.Println(out.String())
	if stderr.Len() > 0 {
		fmt.Println("#### Stderr ####")
		fmt.Println(stderr.String())
	}
	return nil
}

func (executor *Executor) hasResolvedAllVariables(request string) bool {
	return strings.Index(request, "{") == -1
}
