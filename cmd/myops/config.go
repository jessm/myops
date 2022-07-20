package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

const configFile string = "/var/myops/myops_config.json"

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

func parseConfig() map[string]Config {
	_, err := os.Stat(configFile)
	if errors.Is(err, os.ErrNotExist) {
		writeSampleConfig()
	}

	var configs map[string]Config
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
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
			RepoUrl:       "https://github.com/jessm/myops",
			Branch:        "myops-container",
			Dockerfile:    "Dockerfile.ping",
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
