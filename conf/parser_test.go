// Package conf implements functions to read configuration data
package conf

import (
	"flag"
	"os"
	"reflect"
	"testing"
)

func TestNewInitializedConfigData(t *testing.T) {
	testConfigData := NewInitializedConfigData()

	if reflect.TypeOf(testConfigData).Name() != "ConfigData" {
		t.Error("Type error getting initialized configuration data struct")
	}
}

func TestInitialize(t *testing.T) {
	Initialize()
	if Data.APIKey != "YOUR API KEY" {
		t.Error("Data struct not initilized")
	}

}

func TestLoadFromFile(t *testing.T) {
	LoadFromFile("../test/dataset/test.conf.json")
	if Data.APIKey != "abcd1234567890" {
		t.Error("Data struct has not been initialized")
	}
}

// TestWriteDefaultConfigFile tests if the config file with default values could be written
func TestWriteDefaultConfigFile(t *testing.T) {
	WriteDefaultConfigFile("/tmp/splitio.test.config.json")
	if _, err := os.Stat("/tmp/splitio.test.config.json"); os.IsNotExist(err) {
		t.Error("The default config file could not be written")
	}
}

func TestCliParametersToRegister(t *testing.T) {
	parametersToRegister := CliParametersToRegister()

	if len(parametersToRegister) == 0 {
		t.Error("The parameters to be registered have not been found ")
	}
}

func TestLoadFromArgs(t *testing.T) {
	// dinamically configuration parameters
	cliParameters := CliParametersToRegister()
	cliParametersMap := make(map[string]interface{}, len(cliParameters))
	for _, param := range cliParameters {
		switch param.AttributeType {
		case "string":
			cliParametersMap[param.Command] = flag.String(param.Command, "some_value", param.Description)
			break
		case "int":
			cliParametersMap[param.Command] = flag.Int(param.Command, 999, param.Description)
			break
		case "int64":
			cliParametersMap[param.Command] = flag.Int64(param.Command, 999, param.Description)
			break
		case "bool":
			cliParametersMap[param.Command] = flag.Bool(param.Command, true, param.Description)
			break
		}
	}

	flag.Parse()

	LoadFromArgs(cliParametersMap)

}
