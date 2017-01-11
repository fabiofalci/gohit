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
)

type Executor struct {
	conf *Configuration
}

func (executor *Executor) RunRequest(requestName string) error {
	request := executor.conf.Requests[requestName]
	if request != nil {
		executor.runExecutable(request)
		return nil
	}

	endpoint := executor.conf.Endpoints[requestName]
	if endpoint != nil {
		m := make(map[interface{}]interface{})
		m["endpoint"] = requestName
		request, err := executor.conf.createRequest(requestName, m)
		if err != nil {
			return err
		}
		executor.runExecutable(request)
	}
	return nil
}

func (executor *Executor) runExecutable(executable Executable) {
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

	executor.executeCurlCommand(requestAsString)
}

func (executor *Executor) executeCurlCommand(command string) {
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

func (executor *Executor) hasResolvedAllVariables(request string) bool {
	return strings.Index(request, "{") == -1
}
