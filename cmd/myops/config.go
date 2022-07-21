package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

const (
	configFile    string = "/var/myops/myops_config.json"
	oldConfigFile string = "/var/myops/previous_myops_config.json"
)

type Config struct {
	DomainMatcher string            `json:"domainMatcher"`
	RepoUrl       string            `json:"repoUrl"`
	Branch        string            `json:"branch"`
	Dockerfile    string            `json:"dockerfile"`
	EnvVars       map[string]string `json:"envVars"`
	Port          string            `json:"port"`
	VolumePath    string            `json:"volumePath"`
}

type Configs map[string]Config

func getConfigs() Configs {
	_, err := os.Stat(configFile)
	if errors.Is(err, os.ErrNotExist) {
		writeSampleConfig()
	}

	return parseConfig(configFile)
}

func getOldConfigs() Configs {
	_, err := os.Stat(oldConfigFile)
	if errors.Is(err, os.ErrNotExist) {
		return Configs{}
	}

	return parseConfig(oldConfigFile)
}

func parseConfig(fileName string) Configs {

	var configs map[string]Config
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	for name, config := range configs {
		if config.DomainMatcher == "" {
			panic("Config parse error: domainMatcher can't be empty, project " + name)
		}
		if config.RepoUrl == "" {
			panic("Config parse error: repoUrl can't be empty, project " + name)
		}
		if config.Branch == "" {
			config.Branch = "main"
		}
		if config.Dockerfile == "" {
			config.Dockerfile = "Dockerfile"
		}
		if config.Port == "" {
			config.Port = "80"
		}
		// If volumePath is empty, just don't mount volumes
		// If envVars are empty, just don't add env vars
	}

	err = json.Unmarshal(content, &configs)
	if err != nil {
		panic(err)
	}

	return configs
}

func writeSampleConfig() {
	sample := map[string]Config{
		"ping": {
			DomainMatcher: "localhost:3000",
			RepoUrl:       "https://github.com/jessm/ping.git",
			Branch:        "ping",
			Dockerfile:    "go/Dockerfile",
			EnvVars: map[string]string{
				"LISTENPORT": "8000",
			},
			Port:       "8000",
			VolumePath: "/app",
		},
	}

	jsonSample, err := json.MarshalIndent(sample, "", "	")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(configFile, jsonSample, 0777)
	if err != nil {
		panic(err)
	}
}

func writeConfigToOldConfig() {
	bytesRead, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	//Copy all the contents to the desitination file
	err = ioutil.WriteFile(oldConfigFile, bytesRead, 0777)
	if err != nil {
		panic(err)
	}
}
