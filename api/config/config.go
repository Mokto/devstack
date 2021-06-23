package config

import (
	"devstack/errors"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/logrusorgru/aurora"
)

type Service struct {
	Command          string   `json:"command"`
	Name             string   `json:"name"`
	Color            string   `json:"color"`
	Cwd              string   `json:"cwd"`
	WatchDirectories []string `json:"watchDirectories"`
	IsWatching       bool     `json:"isWatching"`
}

func (s *Service) Log(anything string) {
	var applyColor func(str interface{}) aurora.Value
	if s.Color == "blue" {
		applyColor = aurora.Blue
	}
	if s.Color == "cyan" {
		applyColor = aurora.Cyan
	}
	if s.Color == "yellow" {
		applyColor = aurora.Yellow
	}
	if applyColor == nil {
		panic("Color " + s.Color + " not supported.")
	}
	fmt.Println(applyColor("["+s.Name+"] ").String() + anything)
}

type ConfigurationFile struct {
	Services []*Service `json:"services"`
}

func ReadConfigurationFile() (configFile *ConfigurationFile, err error) {

	content, err := ioutil.ReadFile("../example/test.json")
	if err != nil {
		return nil, errors.Wrap(err)
	}

	err = json.Unmarshal(content, &configFile)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return configFile, nil
}
