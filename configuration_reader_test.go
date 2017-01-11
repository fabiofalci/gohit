package main

import (
	"testing"
)

func TestLoadAllFilesConfigurationReader(t *testing.T) {
	confReader := NewSilentConfigurationReader(true, "_resources/valid", "")

	if err := confReader.Read(); err != nil {
		t.Errorf("Should not throw an error '%v'", err)
	}

	if confReader.Directory() != "_resources/valid" {
		t.Error("Directory should be _resources/valid")
	}

	configuration := confReader.Configuration()
	if len(configuration) !=2 {
		t.Error("Should have read 2 configuration files")
	}
}

func TestLoadFilesConfigurationReader(t *testing.T) {
	confReader := NewSilentConfigurationReader(false, "_resources/valid", "api-requests.yaml")

	if err := confReader.Read(); err != nil {
		t.Errorf("Should not throw an error '%v'", err)
	}

	if confReader.Directory() != "_resources/valid" {
		t.Error("Directory should be _resources/valid")
	}

	configuration := confReader.Configuration()
	if len(configuration) != 1 {
		t.Error("Should have read 1 configuration file")
	}
}

func TestOptionalYamlConfigurationReader(t *testing.T) {
	confReader := NewSilentConfigurationReader(false, "_resources/valid", "api-requests")

	if err := confReader.Read(); err != nil {
		t.Errorf("Should not throw an error '%v'", err)
	}

	if confReader.Directory() != "_resources/valid" {
		t.Error("Directory should be _resources/valid")
	}

	configuration := confReader.Configuration()
	if len(configuration) != 1 {
		t.Error("Should have read 1 configuration file")
	}
}


