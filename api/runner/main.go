package runner

import (
	"devstack/config"
	"devstack/websockets"
)

type Runner struct {
	config      *config.ConfigurationFile
	connections *websockets.Connections
	services    map[string]*ServiceRunner
	Logs        []Data
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

func (r *Runner) IsWatching(serviceName string) bool {
	return r.services[serviceName].IsWatching
}

func (r *Runner) SetWatching(serviceName string, isWatching bool) {
	if isWatching {
		r.services[serviceName].watch()
	} else {
		r.services[serviceName].stopWatching()
	}
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
