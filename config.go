package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Filters      []JobDef          `yaml:"jobFilter"`
	JobPipelines map[string]string `yaml:"jobPipeline"`
	JobStage     []JobStage        `yaml:"jobStage"`
}

type JobDef struct {
	JobName string `yaml:"jobName"`
	JobType string `yaml:"jobType"`
}

type JobStage struct {
	Suffix string
	Stage  string
}

func loadConfig(c *Config, f string) {

	b, err := ioutil.ReadFile(f)

	if err != nil {
		fmt.Printf("Error reading file: %s\n")
	}

	err = yaml.Unmarshal(b, c)

	if err != nil {
		fmt.Printf("Error marshaling file: %v\n", err)
	}
	return
}
