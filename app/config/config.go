// Package config defines application's configurations
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/pottava/docker-webui/app/misc"
)

func defaultConfig() Config {
	gopath := os.Getenv("GOPATH")
	return Config{
		Name:                   "docker web-ui",
		Port:                   9000,
		LogLevel:               4,
		DockerEndpoint:         "unix:///var/run/docker.sock",
		DockerAPIVersion:       "1.17",
		DockerPullBeginTimeout: 3 * time.Minute,
		DockerPullTimeout:      2 * time.Hour,
		DockerStatTimeout:      1 * time.Second,
		DockerStartTimeout:     10 * time.Second,
		DockerStopTimeout:      10 * time.Second,
		DockerRestartTimeout:   10 * time.Second,
		DockerKillTimeout:      10 * time.Second,
		DockerRmTimeout:        5 * time.Minute,
		DockerCommitTimeout:    30 * time.Second,
		StaticFileHost:         "",
		StaticFilePath:         gopath + "/src/github.com/pottava/docker-webui/app",
		PreventSelfStop:        true,
	}
}

// NewConfig returns a config struct created by
// merging environment variables and a config file.
func NewConfig() *Config {
	temp := environmentConfig()
	config := &temp

	if !config.complete() {
		config.merge(fileConfig())
	}
	defer func() {
		config.merge(defaultConfig())
		config.trimWhitespace()
	}()
	return config
}

func environmentConfig() Config {
	return Config{
		Name:                   os.Getenv("APP_NAME"),
		Port:                   misc.ParseUint16(os.Getenv("APP_PORT")),
		LogLevel:               misc.Atoi(os.Getenv("APP_LOG_LEVEL")),
		DockerEndpoint:         os.Getenv("APP_DOCKER_ENDPOINT"),
		DockerAPIVersion:       os.Getenv("APP_DOCKER_API_VERSION"),
		DockerPullBeginTimeout: misc.ParseDuration(os.Getenv("APP_DOCKER_PULL_BEGIN_TIMEOUT")),
		DockerPullTimeout:      misc.ParseDuration(os.Getenv("APP_DOCKER_PULL_TIMEOUT")),
		DockerStatTimeout:      misc.ParseDuration(os.Getenv("APP_DOCKER_STAT_TIMEOUT")),
		DockerStartTimeout:     misc.ParseDuration(os.Getenv("APP_DOCKER_START_TIMEOUT")),
		DockerStopTimeout:      misc.ParseDuration(os.Getenv("APP_DOCKER_STOP_TIMEOUT")),
		DockerRestartTimeout:   misc.ParseDuration(os.Getenv("APP_DOCKER_RESTART_TIMEOUT")),
		DockerKillTimeout:      misc.ParseDuration(os.Getenv("APP_DOCKER_KILL_TIMEOUT")),
		DockerRmTimeout:        misc.ParseDuration(os.Getenv("APP_DOCKER_RM_TIMEOUT")),
		DockerCommitTimeout:    misc.ParseDuration(os.Getenv("APP_DOCKER_COMMIT_TIMEOUT")),
		StaticFileHost:         os.Getenv("APP_STATIC_FILE_HOST"),
		StaticFilePath:         os.Getenv("APP_STATIC_FILE_PATH"),
		PreventSelfStop:        misc.ParseBool(os.Getenv("APP_PREVENT_SELF_STOP")),
	}
}

func fileConfig() Config {
	path := misc.NVL(os.Getenv("CONFIG_FILE_PATH"), "/etc/docker-webui/config.json")
	file, err := os.Open(path)
	if err != nil {
		return Config{}
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("Unable to read config file", "err:", err)
		return Config{}
	}
	if strings.TrimSpace(string(data)) == "" {
		return Config{}
	}
	config := Config{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("Error reading config json data. [message] ", err)
	}
	return config
}

func (config *Config) merge(arg Config) *Config {
	mine := reflect.ValueOf(config).Elem()
	theirs := reflect.ValueOf(&arg).Elem()

	for i := 0; i < mine.NumField(); i++ {
		myField := mine.Field(i)
		if misc.ZeroOrNil(myField.Interface()) {
			myField.Set(reflect.ValueOf(theirs.Field(i).Interface()))
		}
	}
	return config
}

func (config *Config) complete() bool {
	cfgs := reflect.ValueOf(config).Elem()

	for i := 0; i < cfgs.NumField(); i++ {
		if misc.ZeroOrNil(cfgs.Field(i).Interface()) {
			return false
		}
	}
	return true
}

func (config *Config) trimWhitespace() {
	cfgs := reflect.ValueOf(config).Elem()
	cfgAttrs := reflect.Indirect(reflect.ValueOf(config)).Type()

	for i := 0; i < cfgs.NumField(); i++ {
		field := cfgs.Field(i)
		if !field.CanInterface() {
			continue
		}
		attr := cfgAttrs.Field(i).Tag.Get("trim")
		if len(attr) == 0 {
			continue
		}
		if field.Kind() != reflect.String {
			continue
		}
		str := field.Interface().(string)
		field.SetString(strings.TrimSpace(str))
	}
}

// String returns a string representation of the config.
func (config *Config) String() string {
	return fmt.Sprintf(
		"Name: %v, Port: %v, LogLevel: %v, DockerEndpoint: %v, DockerAPIVersion: %v, "+
			"DockerPullBeginTimeout: %v, DockerPullTimeout: %v, DockerStatTimeout: %v, DockerStartTimeout: %v, "+
			"DockerStopTimeout: %v, DockerRestartTimeout: %v, DockerKillTimeout: %v, DockerRmTimeout: %v, "+
			"DockerCommitTimeout: %v, StaticFileHost: %v, StaticFilePath: %v, PreventSelfStop: %v",
		config.Name, config.Port, config.LogLevel,
		config.DockerEndpoint, config.DockerAPIVersion, config.DockerPullBeginTimeout,
		config.DockerPullTimeout, config.DockerStatTimeout, config.DockerStartTimeout, config.DockerStopTimeout,
		config.DockerRestartTimeout, config.DockerKillTimeout, config.DockerRmTimeout,
		config.DockerCommitTimeout, config.StaticFileHost, config.StaticFilePath, config.PreventSelfStop)
}