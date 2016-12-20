package main

import (
	"fmt"
	"os"
	"text/template"
	"sort"
)

type Printer struct {
	conf *Configuration
}

func (printer *Printer) ShowRequests() {
	keys := make([]string, len(printer.conf.Requests))
	i := 0
	for k := range printer.conf.Requests {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, name := range keys {
		request := printer.conf.Requests[name]
		printer.showExecutable(request)
		fmt.Println("")
	}
}

func (printer *Printer) ShowEndpoints() {
	keys := make([]string, len(printer.conf.Endpoints))
	i := 0
	for k := range printer.conf.Endpoints {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, name := range keys {
		endpoint := printer.conf.Endpoints[name]
		printer.showExecutable(endpoint)
		fmt.Println("")
	}
}

func (printer *Printer) ShowRequestOrEndpoint(requestName string) {
	request := printer.conf.Requests[requestName]
	if request != nil {
		printer.showExecutable(request)
		return
	}

	endpoint := printer.conf.Endpoints[requestName]
	if endpoint != nil {
		printer.showExecutable(endpoint)
	}
}

func (printer *Printer) showExecutable(executable Executable) {
	t := template.Must(template.New("curlTemplate").Parse(showCurlTemplate))
	fmt.Printf("Endpoint %v:\n", executable.GetName())
	t.Execute(os.Stdout, executable)
}
