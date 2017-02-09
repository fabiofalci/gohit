package main

import (
	"testing"
	"errors"
	"reflect"
)

func TestRequestEndpointNotFound(t *testing.T) {
	conf, _:= NewConfiguration(NewSilentConfigurationReader("_resources/valid", "api-requests.yaml"))
	executor := NewDefaultExecutor(conf)

	if err := executor.RunRequest("not-found"); err == nil || err.Error() != "Could not find request/endpoint not-found" {
		t.Error("Should not throw a not found error")
	}
}

func TestExecuteRequest(t *testing.T) {
	reader := &MockReader{configurations: make(map[string][]byte)}
	reader.configurations["test"] = []byte(
`
url: local

endpoints:
  test:
    path: /test

requests:
  my_request:
    endpoint: test
`)

	conf, err := NewConfiguration(reader)
	if err != nil {
		t.Error(err)
		return
	}
	command := []string{"local/test", "-XGET"}
	executor := NewExecutor(conf, &MockCommandRunner{command: command})

	if err := executor.RunRequest("my_request"); err != nil {
		t.Error("Should not throw an error ", err)
	}
}

func TestExecuteEndpoint(t *testing.T) {
	reader := &MockReader{configurations: make(map[string][]byte)}
	reader.configurations["test"] = []byte(
`
url: local

endpoints:
  test:
    path: /test
    query: test=1
`)

	conf, err := NewConfiguration(reader)
	if err != nil {
		t.Error(err)
		return
	}
	command := []string{"local/test?test=1", "-XGET"}
	executor := NewExecutor(conf, &MockCommandRunner{command: command})

	if err := executor.RunRequest("test"); err != nil {
		t.Error("Should not throw an error ", err)
	}
}

func (runner *MockCommandRunner) Run(command []string) error {
	if !reflect.DeepEqual(command, runner.command) {
		return errors.New("CommandRunner array is not correct")
	}
	return nil
}

type MockCommandRunner struct {
	command []string
}
