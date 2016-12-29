package main

import (
	"testing"
	"bytes"
	"strings"
)

func TestShowEndpoints(t *testing.T) {
	conf := NewSilentConfiguration()
	conf.Init(true, "_resources/", "")

	var b bytes.Buffer
	printer := &Printer{conf: conf, writer: &b}

	printer.ShowEndpoints()

	if allEndpoints != strings.Trim(b.String(), " \n\t") {
		t.Error("All endpoints output doesn't look correct")
	}
}

func TestShowRequests(t *testing.T) {
	conf := NewSilentConfiguration()
	conf.Init(true, "_resources/", "")

	var b bytes.Buffer
	printer := &Printer{conf: conf, writer: &b}

	printer.ShowRequests()

	if allRequests != strings.Trim(b.String(), " \n\t") {
		t.Error("All requests output doesn't look correct")
	}
}

var allEndpoints =`Endpoint endpoint1:
curl 'https://localhost/path1' \
        -H 'Accept: application/vnd.github.v3+json' \
        -H 'Authorization: bearer a12b3c' \
        -H 'Custom: value' \
        --compress \
        --silent \
        -s \
        -vvv \
        -XGET

Endpoint endpoint2:
curl 'https://localhost/path2/{variable}/something' \
        -H 'Accept: application/vnd.github.v3+json' \
        -H 'Authorization: bearer a12b3c' \
        -H 'Custom: value' \
        --compress \
        --silent \
        -s \
        -vvv \
        -XGET

Endpoint endpoint3:
curl 'https://localhost/path3' \
        -H 'Accept: application/vnd.github.v3+json' \
        -H 'Authorization: bearer a12b3c' \
        -H 'Content-length: 0' \
        -H 'Custom: value' \
        --compress \
        --silent \
        -s \
        -vvv \
        -XPUT

Endpoint endpoint4:
curl 'https://localhost/path4/{variable}?name={name}&date={date}' \
        -H 'Accept: application/vnd.github.v3+json' \
        -H 'Authorization: bearer a12b3c' \
        -H 'Custom: value' \
        --compress \
        --silent \
        -s \
        -vvv \
        -XPOST

Endpoint endpoint5:
curl 'https://localhost/' \
        -H 'Accept: application/vnd.github.v3+json' \
        -H 'Authorization: bearer a12b3c' \
        -H 'Custom: value' \
        --compress \
        --silent \
        -s \
        -vvv \
        -XDELETE`

var allRequests = `Endpoint request1:
curl 'https://localhost/path1' \
        -H 'Accept: application/vnd.github.v3+json' \
        -H 'Authorization: bearer a12b3c' \
        -H 'Custom: value' \
        --compress \
        --silent \
        -s \
        -vvv \
        -XGET

Endpoint request2:
curl 'https://localhost/path2/value/something' \
        -H 'Accept: application/vnd.github.v3+json' \
        -H 'Authorization: bearer a12b3c' \
        -H 'Custom: value' \
        --compress \
        --silent \
        -s \
        -vvv \
        -XGET

Endpoint request3:
curl 'https://localhost/path3' \
        -H 'Accept: application/vnd.github.v3+json' \
        -H 'Authorization: bearer a12b3c' \
        -H 'Content-length: 0' \
        -H 'Custom: value' \
        --compress \
        --silent \
        -s \
        -vvv \
        -XPUT

Endpoint request4:
curl 'https://localhost/path4/value?name=gohit&date=today' \
        -H 'Accept: application/vnd.github.v3+json' \
        -H 'Authorization: bearer a12b3c' \
        -H 'Custom: value' \
        --compress \
        --silent \
        -s \
        -vvv \
        -XPOST

Endpoint request4_1:
curl 'https://localhost/path4/value?name=gohit1&date=today1' \
        -H 'Accept: application/vnd.github.v3+json' \
        -H 'Authorization: bearer a12b3c' \
        -H 'Custom: value' \
        --compress \
        --silent \
        -s \
        -vvv \
        -XPOST

Endpoint request5:
curl 'https://localhost/' \
        -H 'Accept: application/vnd.github.v3+json' \
        -H 'Authorization: bearer a12b3c' \
        -H 'Custom: value' \
        --compress \
        --silent \
        -s \
        -vvv \
        -XDELETE`
