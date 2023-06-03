package env

import (
	"os"
	"reflect"
	"testing"
)

type Config struct {
	Port int    `env:"PORT"`
	Host string `env:"HOST"`
}

func TestParse(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("PORT", "8080")
	os.Setenv("HOST", "localhost")

	var config Config
	err := Parse(&config)
	if err != nil {
		t.Errorf("Failed to parse environment variables: %v", err)
	}

	expectedConfig := Config{Port: 8080, Host: "localhost"}
	if !reflect.DeepEqual(config, expectedConfig) {
		t.Errorf("Parsed config does not match expected config.\nExpected: %+v\nGot: %+v", expectedConfig, config)
	}
}

func TestParse_EnvironmentVariablesNotSet(t *testing.T) {
	// Clear environment variables for testing
	os.Clearenv()

	var config Config
	err := Parse(&config)
	if err == nil {
		t.Error("Expected an error while parsing unset environment variables")
	}

	// Ensure that the config remains unchanged
	expectedConfig := Config{}
	if !reflect.DeepEqual(config, expectedConfig) {
		t.Errorf("Parsed config does not match expected config.\nExpected: %+v\nGot: %+v", expectedConfig, config)
	}
}

func TestParse_InvalidEnvironmentVariable(t *testing.T) {
	os.Setenv("PORT", "invalid")

	var config Config
	err := Parse(&config)
	if err == nil {
		t.Error("Expected an error while parsing invalid environment variable")
	}
}

func TestParse_NestedStruct(t *testing.T) {
	os.Setenv("PORT", "9090")
	os.Setenv("DSN", "localhost")

	type NestedConfig struct {
		DSN string `env:"DSN"`
	}

	type Config struct {
		Port   int `env:"PORT"`
		Nested NestedConfig
	}

	var config Config
	err := Parse(&config)
	if err != nil {
		t.Errorf("Failed to parse environment variables: %v", err)
	}

	expectedConfig := Config{Port: 9090, Nested: NestedConfig{DSN: "localhost"}}
	if !reflect.DeepEqual(config, expectedConfig) {
		t.Errorf("Parsed config does not match expected config.\nExpected: %+v\nGot: %+v", expectedConfig, config)
	}
}
