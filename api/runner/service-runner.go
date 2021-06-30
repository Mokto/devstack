package runner

import (
	"devstack/config"
	"devstack/websockets"
	"encoding/json"
	"os"
	"os/exec"

	"github.com/fsnotify/fsnotify"
)

type Data struct {
	Service *config.Service `json:"service"`
	Message string          `json:"message"`
}

type Message struct {
	Data      Data   `json:"data"`
	EventName string `json:"eventName"`
}

type ServiceRunner struct {
	runner              *Runner
	cmd                 *exec.Cmd
	watcher             *fsnotify.Watcher
	stopWatchingChannel chan bool
	service             *config.Service
	IsWatching          bool
	IsRunning           bool
}

func (serviceRunner *ServiceRunner) Init() {
	go func() {
		serviceRunner.execCommand()
	}()
	if len(serviceRunner.service.WatchDirectories) > 0 {
		go func() {
			serviceRunner.watch()
		}()
	}
}

func (serviceRunner *ServiceRunner) Restart() {
	serviceRunner.Stop()
	serviceRunner.Init()
}

func (serviceRunner *ServiceRunner) Stop() {
	if serviceRunner.cmd.Process == nil {
		panic("Process is nil")
	}
	err := serviceRunner.cmd.Process.Signal(os.Kill)
	if err != nil {
		panic(err)
	}
}

func (serviceRunner *ServiceRunner) SendLog(message string) {

	bytes, _ := json.Marshal(Message{
		EventName: "log",
		Data: Data{
			Message: message,
			Service: serviceRunner.service,
		},
	})
	for _, user := range serviceRunner.runner.connections.Users {
		err := websockets.Send(user, string(bytes))
		if err != nil {
			panic(err)
		}
	}
}

type State struct {
	ServiceName string `json:"serviceName"`
	IsRunning   bool   `json:"isRunning"`
}
type StateWebsocket struct {
	EventName string `json:"eventName"`
	State     State  `json:"data"`
}

func (serviceRunner *ServiceRunner) SendIsRunning(isRunning bool) {

	bytes, _ := json.Marshal(StateWebsocket{
		EventName: "isRunning",
		State: State{
			IsRunning:   isRunning,
			ServiceName: serviceRunner.service.Name,
		},
	})
	for _, user := range serviceRunner.runner.connections.Users {
		err := websockets.Send(user, string(bytes))
		if err != nil {
			panic(err)
		}
	}
}
