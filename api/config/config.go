package config

import (
	"devstack/errors"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/logrusorgru/aurora"
)

type Service struct {
	Command          string            `json:"command"`
	Name             string            `json:"name"`
	Color            string            `json:"color"`
	Cwd              string            `json:"cwd"`
	WatchDirectories []string          `json:"watchDirectories"`
	Env              map[string]string `json:"env"`
	// dynamic
	IsWatching bool `json:"isWatching"`
}

type ConfigurationFile struct {
	BasePath string
	Env      map[string]string `json:"env"`
	Services []*Service        `json:"services"`
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

func ReadConfigurationFile() (configFile *ConfigurationFile, err error) {
	directory := os.Getenv("PROJECT_PATH")
	if directory == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, errors.Wrap(err)
		}
		directory = cwd
	}
	filePath := directory + "/.services.json"
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	err = json.Unmarshal(content, &configFile)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	configFile.BasePath = directory

	return configFile, nil
}
