package main

import (
	"testing"
	"bytes"
	"strings"
)

func TestShowEndpointsOneLine(t *testing.T) {
	conf, _ := NewConfiguration(NewSilentConfigurationReader("_resources/valid", "api-requests.yaml"))

	var b bytes.Buffer
	printer := &Printer{conf: conf, writer: &b, oneLine: true}

	printer.ShowEndpoints()

	if allEndpointsOutputOneLine != strings.Trim(b.String(), " \n\t") {
		t.Error("All endpoints output doesn't look correct")
	}
}

func TestShowRequestsOneLine(t *testing.T) {
	conf, _ := NewConfiguration(NewSilentConfigurationReader("_resources/valid", "api-requests.yaml"))

	var b bytes.Buffer
	printer := &Printer{conf: conf, writer: &b, oneLine: true}

	printer.ShowRequests()

	if allRequestsOutputOneLine != strings.Trim(b.String(), " \n\t") {
		t.Error("All requests output doesn't look correct")
	}
}

func TestShowEndpointOneLine(t *testing.T) {
	conf, _ := NewConfiguration(NewSilentConfigurationReader("_resources/valid", "api-requests.yaml"))

	var b bytes.Buffer
	printer := &Printer{conf: conf, writer: &b, oneLine: true}

	printer.ShowRequestOrEndpoint("endpoint1")

	if endpoint1OutputOneLine != strings.Trim(b.String(), " \n\t") {
		t.Error("endpoint1 output doesn't look correct")
	}
}

func TestShowRequestOneLine(t *testing.T) {
	conf, _ := NewConfiguration(NewSilentConfigurationReader("_resources/valid", "api-requests.yaml"))

	var b bytes.Buffer
	printer := &Printer{conf: conf, writer: &b, oneLine: true}

	printer.ShowRequestOrEndpoint("request1")

	if request1OutputOneLine != strings.Trim(b.String(), " \n\t") {
		t.Error("endpoint1 output doesn't look correct")
	}
}

var endpoint1OutputOneLine = `Endpoint endpoint1:
curl 'https://localhost/path1' -H 'Accept: application/vnd.github.v3+json' -H 'Authorization: bearer a12b3c' -H 'Custom: value' -G --data-urlencode 'format=json' --data-urlencode 'version=v2' --compress --silent -s -vvv -XGET`

var request1OutputOneLine =`Endpoint request1:
curl 'https://localhost/path1' -H 'Accept: application/vnd.github.v3+json' -H 'Authorization: bearer a12b3c' -H 'Custom: value' -G --data-urlencode 'format=json' --data-urlencode 'version=v2' --compress --silent -s -vvv -XGET`

var allEndpointsOutputOneLine =`Endpoint endpoint1:
curl 'https://localhost/path1' -H 'Accept: application/vnd.github.v3+json' -H 'Authorization: bearer a12b3c' -H 'Custom: value' -G --data-urlencode 'format=json' --data-urlencode 'version=v2' --compress --silent -s -vvv -XGET

Endpoint endpoint2:
curl 'https://localhost/path2/{variable}/something' -H 'Accept: application/vnd.github.v3+json' -H 'Authorization: bearer a12b3c' -H 'Custom: value' --compress --silent -s -vvv -XGET

Endpoint endpoint3:
curl 'https://localhost/path3' -H 'Accept: application/vnd.github.v3+json' -H 'Authorization: bearer a12b3c' -H 'Content-length: 0' -H 'Custom: value' --compress --silent -s -vvv -XPUT

Endpoint endpoint4:
curl 'https://localhost/path4/{variable}?name={name}&date={date}' -H 'Accept: application/vnd.github.v3+json' -H 'Authorization: bearer a12b3c' -H 'Custom: value' --compress --silent -s -vvv -XPOST

Endpoint endpoint5:
curl 'https://localhost/' -H 'Accept: application/vnd.github.v3+json' -H 'Authorization: bearer a12b3c' -H 'Custom: value' --compress --silent -s -vvv -XDELETE`

var allRequestsOutputOneLine = `Endpoint request1:
curl 'https://localhost/path1' -H 'Accept: application/vnd.github.v3+json' -H 'Authorization: bearer a12b3c' -H 'Custom: value' -G --data-urlencode 'format=json' --data-urlencode 'version=v2' --compress --silent -s -vvv -XGET

Endpoint request2:
curl 'https://localhost/path2/value/something' -H 'Accept: application/vnd.github.v3+json' -H 'Authorization: bearer a12b3c' -H 'Custom: value' --compress --silent -s -vvv -XGET

Endpoint request3:
curl 'https://localhost/path3' -H 'Accept: application/vnd.github.v3+json' -H 'Authorization: bearer a12b3c' -H 'Content-length: 0' -H 'Custom: value' --compress --silent -s -vvv -XPUT

Endpoint request4:
curl 'https://localhost/path4/value?name=gohit&date=today' -H 'Accept: application/vnd.github.v3+json' -H 'Authorization: bearer a12b3c' -H 'Custom: value' --compress --silent -s -vvv -XPOST

Endpoint request4_1:
curl 'https://localhost/path4/value?name=gohit1&date=today1' -H 'Accept: application/vnd.github.v3+json' -H 'Authorization: bearer a12b3c' -H 'Custom: value' --compress --silent -s -vvv -XPOST

Endpoint request5:
curl 'https://localhost/' -H 'Accept: application/vnd.github.v3+json' -H 'Authorization: bearer a12b3c' -H 'Custom: value' --compress --silent -s -vvv -XDELETE`
