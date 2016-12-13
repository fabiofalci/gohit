package main

import (
	"regexp"
	"strings"
)

const showCurlTemplate = `curl '{{.Url}}{{.Path}}{{if .Query}}?{{.Query}}{{end}}' \
{{- if .Headers}}
        {{- range $key, $value := .Headers }}
        -H '{{$key}}' \
        {{- end}}
{{- end}}
{{- if .Options}}
        {{- range $key, $value := .Options }}
        {{$key}} \
        {{- end}}
{{- end}}
        -X{{.Method}}
`

// A separated template for running as it needs to transform the command to an array fo string.
// It splits on newline.
const runCurlTemplate = `{{.Url}}{{.Path}}{{if .Query}}?{{.Query}}{{end}}
{{- if .Headers}}
        {{- range $key, $value := .Headers }}
-H
{{$key}}
        {{- end}}
{{- end}}
{{- if .Options}}
{{.OptionsAsToken}}
{{- end}}
-X{{.Method}}`

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
