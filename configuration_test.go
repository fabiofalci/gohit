package main

import (
	"testing"
	"errors"
)

func TestLoadAllConfigurationGlobal(t *testing.T) {
	conf, err := NewConfiguration(NewSilentConfigurationReader("_resources/valid", "api-requests.yaml"))

	if err != nil {
		t.Errorf("Should not throw an error '%v'", err)
	}

	if len(conf.GlobalHeaders) != 3 {
		t.Error("Should be 3 global headers")
	}

	if len(conf.GlobalOptions) != 4 {
		t.Error("Should be 4 global options")
	}

	if conf.GlobalUrl != "https://localhost" {
		t.Error("Should be httos://localhost")
	}

	if len(conf.GlobalVariables) != 1 {
		t.Error("Should be 1 global variables")
	}
}

func TestLoadAllConfigurationEndpoints(t *testing.T) {
	conf, _ := NewConfiguration(NewSilentConfigurationReader("_resources/valid", "api-requests.yaml"))

	if len(conf.Endpoints) != 5 {
		t.Error("Should be 5 endpoints")
	}

	endpoint1 := conf.Endpoints["endpoint1"]
	if endpoint1.Name != "endpoint1" ||
		endpoint1.Path != "/path1" ||
		endpoint1.Url != "https://localhost" ||
		endpoint1.Query != "" ||
		len(endpoint1.Headers) != 3 ||
		len(endpoint1.Options) != 4 ||
		endpoint1.Method != "GET" {
		t.Errorf("Endpoint1 configuration problem %v", endpoint1)
	}

	endpoint2 := conf.Endpoints["endpoint2"]
	if endpoint2.Name != "endpoint2" ||
		endpoint2.Path != "/path2/{variable}/something" ||
		endpoint2.Url != "https://localhost" ||
		endpoint2.Query != "" ||
		len(endpoint2.Headers) != 3 ||
		len(endpoint2.Options) != 4 ||
		endpoint2.Method != "GET" {
		t.Errorf("Endpoint2 configuration problem %v", endpoint2)
	}

	endpoint3 := conf.Endpoints["endpoint3"]
	if endpoint3.Name != "endpoint3" ||
		endpoint3.Path != "/path3" ||
		endpoint3.Url != "https://localhost" ||
		endpoint3.Query != "" ||
		len(endpoint3.Headers) != 4 ||
		len(endpoint3.Options) != 4 ||
		endpoint3.Method != "PUT" {
		t.Errorf("Endpoint3 configuration problem %v", endpoint3)
	}

	endpoint4 := conf.Endpoints["endpoint4"]
	if endpoint4.Name != "endpoint4" ||
		endpoint4.Path != "/path4/{variable}" ||
		endpoint4.Url != "https://localhost" ||
		endpoint4.Query != "name={name}&date={date}" ||
		len(endpoint4.Headers) != 3 ||
		len(endpoint4.Options) != 4 ||
		endpoint4.Method != "POST" {
		t.Errorf("Endpoint4 configuration problem %v", endpoint4)
	}

	endpoint5 := conf.Endpoints["endpoint5"]
	if endpoint5.Name != "endpoint5" ||
		endpoint5.Path != "/" ||
		endpoint5.Url != "https://localhost" ||
		endpoint5.Query != "" ||
		len(endpoint5.Headers) != 3 ||
		len(endpoint5.Options) != 4 ||
		endpoint5.Method != "DELETE" {
		t.Errorf("Endpoint5 configuration problem %v", endpoint5)
	}
}

func TestLoadAllConfigurationRequests(t *testing.T) {
	conf, _ := NewConfiguration(NewSilentConfigurationReader("_resources/valid", "api-requests.yaml"))

	if len(conf.Requests) != 6 {
		t.Error("Should be 6 endpoints")
	}

	request1 := conf.Requests["request1"]
	if request1.Name != "request1" ||
		request1.Path != "/path1" ||
		request1.Url != "https://localhost" ||
		request1.Query != "" ||
		len(request1.Headers) != 3 ||
		len(request1.Options) != 4 ||
		len(request1.Parameters) != 1 ||
		request1.Method != "GET" {
		t.Errorf("Request1 configuration problem %v", request1)
	}

	request2 := conf.Requests["request2"]
	if request2.Name != "request2" ||
		request2.Path != "/path2/value/something" ||
		request2.Url != "https://localhost" ||
		request2.Query != "" ||
		len(request2.Headers) != 3 ||
		len(request2.Options) != 4 ||
		len(request2.Parameters) != 1 ||
		request2.Method != "GET" {
		t.Errorf("Endpoint2 configuration problem %v", request2)
	}

	request3 := conf.Requests["request3"]
	if request3.Name != "request3" ||
		request3.Path != "/path3" ||
		request3.Url != "https://localhost" ||
		request3.Query != "" ||
		len(request3.Headers) != 4 ||
		len(request3.Options) != 4 ||
		len(request3.Parameters) != 1 ||
		request3.Method != "PUT" {
		t.Errorf("Request4 configuration problem %v", request3)
	}

	request4 := conf.Requests["request4"]
	if request4.Name != "request4" ||
		request4.Path != "/path4/value" ||
		request4.Url != "https://localhost" ||
		request4.Query != "name=gohit&date=today" ||
		len(request4.Headers) != 3 ||
		len(request4.Options) != 4 ||
		len(request4.Parameters) != 3 ||
		request4.Method != "POST" {
		t.Errorf("Request4 configuration problem %v", request4)
	}

	request41 := conf.Requests["request4_1"]
	if request41.Name != "request4_1" ||
		request41.Path != "/path4/value" ||
		request41.Url != "https://localhost" ||
		request41.Query != "name=gohit1&date=today1" ||
		len(request41.Headers) != 3 ||
		len(request41.Options) != 4 ||
		len(request41.Parameters) != 3 ||
		request4.Method != "POST" {
		t.Errorf("Request4_1 configuration problem %v", request4)
	}

	request5 := conf.Requests["request5"]
	if request5.Name != "request5" ||
		request5.Path != "/" ||
		request5.Url != "https://localhost" ||
		request5.Query != "" ||
		len(request5.Headers) != 3 ||
		len(request5.Options) != 4 ||
		len(request5.Parameters) != 1 ||
		request5.Method != "DELETE" {
		t.Errorf("Request5 configuration problem %v", request5)
	}
}

func TestOverrideUrl(t *testing.T) {
	conf, _ := NewConfiguration(NewSilentConfigurationReader("_resources/valid", "api-requests2.yaml"))

	if len(conf.Requests) != 6 {
		t.Error("Should be 6 endpoints")
	}

	request1 := conf.Requests["request1"]
	if request1.Name != "request1" ||
		request1.Path != "/path1" ||
		request1.Url != "https://127.0.0.1" ||
		request1.Query != "" ||
		len(request1.Headers) != 3 ||
		len(request1.Options) != 4 ||
		len(request1.Parameters) != 1 ||
		request1.Method != "GET" {
		t.Errorf("Request1 configuration problem %v", request1)
	}
}

func TestLoadMissingEndpoints(t *testing.T) {
	reader := &MockReader{configurations: make(map[string][]byte)}
	reader.configurations["test"] = []byte(
`
url: local
`)


	if _, err := NewConfiguration(reader); err == nil || err.Error() != "Missing endpoints" {
		t.Error("Should have thrown a missing endpoints error")
	}
}

func TestLoadMissingUrl(t *testing.T) {
	reader := &MockReader{configurations: make(map[string][]byte)}
	reader.configurations["test"] = []byte(
`
endpoints:
  test:
    path: /test
`)


	if _, err := NewConfiguration(reader); err == nil || err.Error() != "Missing URL" {
		t.Error("Should have thrown a missing URL error")
	}
}

func TestMissingPath(t *testing.T) {
	reader := &MockReader{configurations: make(map[string][]byte)}
	reader.configurations["test"] = []byte(
`
url: local

endpoints:
  test:
  method: GET
`)

	if _, err := NewConfiguration(reader); err == nil || err.Error() != "Endpoint test missing path" {
		t.Error("Should have thrown a missing path error but got ", err)
	}
}

func TestRequestUsingMissingEndpoint(t *testing.T) {
	reader := &MockReader{configurations: make(map[string][]byte)}
	reader.configurations["test"] = []byte(
`
url: local

endpoints:
  test:
    path: /test

requests:
  my_request:
    endpoint: test_1
`)


	if _, err := NewConfiguration(reader); err == nil || err.Error() != "Request my_request couldn't find endpoint test_1" {
		t.Error("Should have thrown a missing endpoint error")
	}
}

func TestInvalidAttribute(t *testing.T) {
	reader := &MockReader{configurations: make(map[string][]byte)}
	reader.configurations["test"] = []byte(
`
url: local

invalid: true

endpoints:
  test:
    path: /test

`)

	if _, err := NewConfiguration(reader); err == nil || err.Error() != "Invalid yaml attribute 'invalid'" {
		t.Error("Should have thrown a missing endpoint error")
	}
}

func TestFileNotFound(t *testing.T) {
	reader := &MockReader{configurations: make(map[string][]byte), errorWhenReading: errors.New("Test error")}

	if _, err := NewConfiguration(reader); err == nil {
		t.Error("Should have thrown a file not found error")
	}
}

func (reader *MockReader) Read() error {
	return reader.errorWhenReading
}

func (reader *MockReader) Configuration() map[string][]byte {
	return reader.configurations
}

func (reader *MockReader) Directory() string {
	return "test"
}

type MockReader struct {
	configurations   map[string][]byte
	errorWhenReading error
}