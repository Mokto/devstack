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

func (r *Runner) IsRunning(serviceName string) bool {
	return r.services[serviceName].IsRunning
}

func (r *Runner) SetWatching(serviceName string, isWatching bool) {
	if isWatching {
		r.services[serviceName].watch()
	} else {
		r.services[serviceName].StopWatching()
	}
}

func (r *Runner) SetIsRunning(serviceName string, isRunning bool) {
	if isRunning {
		r.services[serviceName].Restart()
	} else {
		r.services[serviceName].StopWatching()
		r.services[serviceName].Stop()
	}
}

func (r *Runner) StopAll() {
	for _, service := range r.services {
		service.Stop()
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
