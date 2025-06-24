package config

import (
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/spf13/viper"
)

func TestConfig(t *testing.T) {
	// Get all field names from ConfigDefinition struct to verify coverage
	configType := reflect.TypeOf(ConfigDefinition{})
	allFields := make(map[string]bool)
	for i := 0; i < configType.NumField(); i++ {
		allFields[configType.Field(i).Name] = false
	}

	// Set test values
	viper.Set("appName", "TestApp")
	viper.Set("appVersion", "1.0.0")
	viper.Set("ffmpeg", "/usr/bin/ffmpeg")
	viper.Set("port", uint(8080))
	viper.Set("tray", true)
	viper.Set("database", "/path/to/db.sqlite")
	viper.Set("debug", "true")
	viper.Set("loglevel", "trace")
	viper.Set("maxConcurrentTasks", uint(4))
	viper.Set("sendTelemetry", true)
	viper.Set("noUI", true)

	Init()
	c := Config()

	// Test all fields
	tests := []struct {
		name   string
		got    interface{}
		want   interface{}
		errMsg string
	}{
		{"AppName", c.AppName, "TestApp", "AppName mismatch"},
		{"AppVersion", c.AppVersion, "1.0.0", "AppVersion mismatch"},
		{"FFMpeg", c.FFMpeg, "/usr/bin/ffmpeg", "FFMpeg path mismatch"},
		{"Port", c.Port, uint(8080), "Port mismatch"},
		{"Tray", c.Tray, true, "Tray setting mismatch"},
		{"Database", c.Database, "/path/to/db.sqlite", "Database path mismatch"},
		{"Debug", c.Debug, "true", "Debug setting mismatch"},
		{"Loglevel", c.Loglevel, "trace", "Loglevel mismatch"},
		{"MaxConcurrentTasks", c.MaxConcurrentTasks, uint(4), "MaxConcurrentTasks mismatch"},
		{"SendTelemetry", c.SendTelemetry, true, "SendTelemetry mismatch"},
		{"NoUI", c.NoUI, true, "NoUI mismatch"},
		{"Mutex", c.Mutex, sync.RWMutex{}, "Mutex mismatch"},
	}

	// Run tests and track covered fields
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("%s: got %v, want %v", tt.errMsg, tt.got, tt.want)
			}
		})
		delete(allFields, tt.name)
	}

	// Check if any fields are not tested
	if len(allFields) > 0 {
		var missingTests []string
		for field := range allFields {
			missingTests = append(missingTests, field)
		}
		t.Errorf("Missing tests for config fields: %v", missingTests)
	}
}

func TestDebugEnvironmentOverride(t *testing.T) {
	os.Setenv("DEBUGO", "debug_from_env")
	viper.Set("debug", "")

	Init()
	c := Config()

	if c.Debug != "debug_from_env" {
		t.Errorf("Debug environment override failed: got %v, want debug_from_env", c.Debug)
	}

	// Cleanup
	os.Unsetenv("DEBUGO")
}
