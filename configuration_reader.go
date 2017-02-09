package main

import (
	"os"
	"strings"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"io"
	"bytes"
)

type ConfigurationReader struct {
	writer          *io.Writer
	directory       string
	file            string
	configurations  map[string][]byte
}

func NewDefaultConfigurationReader(directory string, file string) *ConfigurationReader {
	return NewConfigurationReader(os.Stdout, directory, file)
}
func NewSilentConfigurationReader(directory string, file string) *ConfigurationReader {
	var out bytes.Buffer
	return NewConfigurationReader(&out, directory, file)
}

func NewConfigurationReader(writer io.Writer, directory string, file string) *ConfigurationReader {
	confReader := &ConfigurationReader{
		writer: &writer,
		directory: directory,
		file: file,
		configurations: make(map[string][]byte),
	}
	return confReader
}

func (confReader *ConfigurationReader) Read() error {
	return confReader.loadConfigurationAndEndpoints()
}

func (confReader *ConfigurationReader) Configuration() map[string][]byte {
	return confReader.configurations
}

func (confReader *ConfigurationReader) Directory() string {
	return confReader.directory
}

func (confReader *ConfigurationReader) loadConfigurationAndEndpoints() error {
	if !strings.HasSuffix(confReader.file, ".yaml") {
		confReader.file = confReader.file + ".yaml"
	}
	source, err := ioutil.ReadFile(confReader.directory + "/" + confReader.file)
	if err != nil {
		return err
	}
	confReader.configurations[confReader.file] = source
	return nil
}

func (confReader *ConfigurationReader) visit(path string, f os.FileInfo, err error) error {
	if !f.IsDir() {
		if strings.HasSuffix(path, ".yaml") {
			fmt.Fprintf(*confReader.writer, "Loading: %s\n", path)
			source, err := ioutil.ReadFile(path)
			if err != nil {
				panic(err)
			}
			confReader.configurations[path] = source
		}
	}
	return nil
}

func (confReader *ConfigurationReader) loadAll() {
	filepath.Walk(confReader.directory, confReader.visit)
}


