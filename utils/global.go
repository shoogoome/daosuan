package utils

import (
	"daosuan/models"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var GlobalConfig models.SystemConfiguration

func InitGlobal() {
	yamlFile, err := ioutil.ReadFile("../config.yaml")
	if err != nil {
		panic(fmt.Errorf("failed to load configuration: %s", err.Error()))
	}
	err = yaml.Unmarshal(yamlFile, &GlobalConfig)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshal configuration: %s", err.Error()))
	}
}