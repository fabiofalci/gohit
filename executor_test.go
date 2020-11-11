package main

import (
	"errors"
	"reflect"
	"testing"
)

func TestRequestEndpointNotFound(t *testing.T) {
	conf, _ := NewConfiguration(NewSilentConfigurationReader("_resources/valid", "api-requests.yaml"))
	executor := NewDefaultExecutor(conf)

	if err := executor.RunRequest("not-found", nil); err == nil || err.Error() != "Could not find request/endpoint not-found" {
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
	executor := NewExecutor(conf, &MockCommandRunner{command: command}, &MockVariableReader{})

	if err := executor.RunRequest("my_request", nil); err != nil {
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
	executor := NewExecutor(conf, &MockCommandRunner{command: command}, &MockVariableReader{})

	if err := executor.RunRequest("test", nil); err != nil {
		t.Error("Should not throw an error ", err)
	}
}

func TestExecuteEndpointWithParams(t *testing.T) {
	reader := &MockReader{configurations: make(map[string][]byte)}
	reader.configurations["test"] = []byte(
		`
url: local

endpoints:
  test:
    path: /test/{param}
    query: test=1
`)

	conf, err := NewConfiguration(reader)
	if err != nil {
		t.Error(err)
		return
	}
	command := []string{"local/test/value?test=1", "-XGET"}
	executor := NewExecutor(conf, &MockCommandRunner{command: command}, &MockVariableReader{})

	if err := executor.RunRequest("test", nil); err != nil {
		t.Error("Should not throw an error ", err)
	}
}

func TestExecuteEndpointWithParamsResolvedByArgs(t *testing.T) {
	reader := &MockReader{configurations: make(map[string][]byte)}
	reader.configurations["test"] = []byte(
		`
url: local

endpoints:
  test:
    path: /test/{param}
    query: test=1
`)

	conf, err := NewConfiguration(reader)
	if err != nil {
		t.Error(err)
		return
	}
	command := []string{"local/test/argParam?test=1", "-XGET"}
	executor := NewExecutor(conf, &MockCommandRunner{command: command}, &MockVariableReader{})

	if err := executor.RunRequest("test", []string{"argParam"}); err != nil {
		t.Error("Should not throw an error ", err)
	}
}

func TestExecuteRequestWithParams(t *testing.T) {
	reader := &MockReader{configurations: make(map[string][]byte)}
	reader.configurations["test"] = []byte(
		`
url: local

endpoints:
  test:
    path: /test/{param}

requests:
  my_request:
    endpoint: test
`)

	conf, err := NewConfiguration(reader)
	if err != nil {
		t.Error(err)
		return
	}
	command := []string{"local/test/value", "-XGET"}
	executor := NewExecutor(conf, &MockCommandRunner{command: command}, &MockVariableReader{})

	if err := executor.RunRequest("my_request", nil); err != nil {
		t.Error("Should not throw an error ", err)
	}
}

func TestExecuteRequestWithParamsResolvedByArgs(t *testing.T) {
	reader := &MockReader{configurations: make(map[string][]byte)}
	reader.configurations["test"] = []byte(
		`
url: local

endpoints:
  test:
    path: /test/{param}

requests:
  my_request:
    endpoint: test
`)

	conf, err := NewConfiguration(reader)
	if err != nil {
		t.Error(err)
		return
	}
	command := []string{"local/test/argParam", "-XGET"}
	executor := NewExecutor(conf, &MockCommandRunner{command: command}, &MockVariableReader{})

	if err := executor.RunRequest("my_request", []string{"argParam"}); err != nil {
		t.Error("Should not throw an error ", err)
	}
}

func (runner *MockCommandRunner) Run(command []string) error {
	if !reflect.DeepEqual(command, runner.command) {
		return errors.New("CommandRunner array is not correct")
	}
	return nil
}

func (parameterReader *MockVariableReader) Read(variableName string) string {
	return "value"
}

type MockCommandRunner struct {
	command []string
}

type MockVariableReader struct {
}
