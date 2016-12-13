package main

import (
	"fmt"
	"os"
	"text/template"
)

type Printer struct {
	conf *Configuration
}

func (printer *Printer) ShowRequests() {
	for name := range printer.conf.Requests {
		request := printer.conf.Requests[name]
		printer.showExecutable(request)
		fmt.Println("")
	}
}

func (printer *Printer) ShowEndpoints() {
	for name := range printer.conf.Endpoints {
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
