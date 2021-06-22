package runner

import (
	"devstack/config"
	"devstack/websockets"
)

type Runner struct {
	Logs        []Data
	config      *config.ConfigurationFile
	connections *websockets.Connections
	services    map[string]*ServiceRunner
}

func (r *Runner) InitAll() {
	for _, service := range r.config.Services {
		r.services[service.Name] = &ServiceRunner{
			service: service,
			runner:  r,
		}
		r.services[service.Name].Init()
	}
}

func (r *Runner) Restart(serviceName string) {
	r.services[serviceName].Restart()
}

func Start(configFile *config.ConfigurationFile, connections *websockets.Connections) *Runner {

	logsHolder := &Runner{
		Logs:        []Data{},
		config:      configFile,
		connections: connections,
		services:    map[string]*ServiceRunner{},
	}

	logsHolder.InitAll()

	return logsHolder
}
