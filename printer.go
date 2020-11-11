package main

import (
	"fmt"
	"text/template"
	"sort"
	"io"
	"regexp"
	"strings"
)

const showCurlTemplate = `curl '{{.Url}}{{.Path}}{{if .QueryRaw}}?{{.QueryRaw}}{{end}}' \
{{- if .Headers}}
        {{- range $key, $value := .Headers }}
        -H '{{$key}}' \
        {{- end}}
{{- end}}
{{- if .QueryList}}
		-G \
        {{- range $key, $value := .QueryList }}
		--data-urlencode '{{$key}}={{$value}}' \
        {{- end}}
{{- end}}
{{- if .Options}}
        {{- range $key, $value := .Options }}
        {{$key}} \
        {{- end}}
{{- end}}
        -X{{.Method}}
`

// A separated template for running as it needs to transform the command to an array of string.
// It splits on newline.
const runCurlTemplate = `{{.Url}}{{.Path}}{{if .QueryRaw}}?{{.QueryRaw}}{{end}}
{{- if .Headers}}
        {{- range $key, $value := .Headers }}
-H
{{$key}}
        {{- end}}
{{- end}}
{{- if .QueryList}}
-G
        {{- range $key, $value := .QueryList }}
--data-urlencode
'{{$key}}={{$value}}'
        {{- end}}
{{- end}}
{{- if .Options}}
{{.OptionsAsToken}}
{{- end}}
-X{{.Method}}`

type Printer struct {
	conf    *Configuration
	writer  io.Writer
	oneLine bool
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
		fmt.Fprintln(printer.writer, "")
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
		fmt.Fprintln(printer.writer, "")
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
	t := template.Must(template.New("curlTemplate").Parse(printer.getTemplate(executable)))
	fmt.Fprintf(printer.writer, "Endpoint %v:\n", executable.GetName())
	t.Execute(printer.writer, executable)
}


func (printer *Printer) getTemplate(executable Executable) string {
	if !printer.oneLine {
		return showCurlTemplate
	}
	t := strings.Replace(showCurlTemplate, "\\", "", -1)
	t = strings.Replace(t, "\n", "", -1)
	return regexp.MustCompile(`[\s\p{Zs}]{2,}`).ReplaceAllString(t, " ") + "\n"
}

// methods used by the templates

func (request *Request) OptionsAsToken() string {
	return executableOptionsAsToken(request)
}

func (endpoint *Endpoint) OptionsAsToken() string {
	return executableOptionsAsToken(endpoint)
}

func executableOptionsAsToken(executable Executable) string {
	oneLineOptions := ""
	for option := range executable.GetOptions() {
		re := regexp.MustCompile("[^\\s\"']+|\"([^\"]*)\"|'([^']*)'")
		for _, v := range re.FindAllString(option, -1) {
			oneLineOptions = oneLineOptions + "\n" + v
		}
	}

	return strings.TrimPrefix(oneLineOptions, "\n")
}

