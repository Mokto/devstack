package runner

import (
	"devstack/config"
	"devstack/websockets"
	"encoding/json"
	"os"
	"os/exec"

	"github.com/fsnotify/fsnotify"
	"github.com/logrusorgru/aurora"
)

type Data struct {
	Message string          `json:"message"`
	Service *config.Service `json:"service"`
}

type Message struct {
	EventName string `json:"eventName"`
	Data      Data   `json:"data"`
}

type ServiceRunner struct {
	runner              *Runner
	cmd                 *exec.Cmd
	watcher             *fsnotify.Watcher
	stopWatchingChannel chan bool
	service             *config.Service
	IsWatching          bool `json:"isWatching"`
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
	serviceRunner.SendLog(aurora.Yellow("Restarting service...").String())
	if serviceRunner.cmd.Process == nil {
		panic("Process is nil")
	}
	err := serviceRunner.cmd.Process.Signal(os.Kill)
	if err != nil {
		panic(err)
	}
	serviceRunner.Init()
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
