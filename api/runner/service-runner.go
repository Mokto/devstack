package runner

import (
	"bufio"
	"devstack/config"
	"devstack/websockets"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/creack/pty"
	"github.com/logrusorgru/aurora"
	"golang.org/x/term"
)

type Data struct {
	Message string         `json:"message"`
	Service config.Service `json:"service"`
}

type Message struct {
	EventName string `json:"eventName"`
	Data      Data   `json:"data"`
}

type ServiceRunner struct {
	service config.Service
	runner  *Runner
	cmd     *exec.Cmd
}

func (serviceRunner *ServiceRunner) Init() {
	go func() {
		serviceRunner.execCommand()
	}()
}

func (serviceRunner *ServiceRunner) Restart() {
	err := serviceRunner.cmd.Process.Signal(os.Kill)
	if err != nil {
		panic(err)
	}
	serviceRunner.Init()
}

func (serviceRunner *ServiceRunner) execCommand() {
	fmt.Println(aurora.Blue("Running " + serviceRunner.service.Command))
	splitted := strings.Split(serviceRunner.service.Command, " ")
	serviceRunner.cmd = exec.Command(splitted[0], splitted[1:]...)

	if serviceRunner.service.Cwd != "" {
		serviceRunner.cmd.Dir = serviceRunner.service.Cwd
	}

	ptmx, err := pty.Start(serviceRunner.cmd)
	if err != nil {
		panic(err)
	}
	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.
	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH                        // Initial resize.
	defer func() { signal.Stop(ch); close(ch) }() // Cleanup signals when done.
	// Set stdin in raw mode.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.
	// Copy stdin to the pty and the pty to stdout.
	// NOTE: The goroutine will keep reading until the next keystroke before returning.
	go func() { _, _ = io.Copy(ptmx, os.Stdin) }()

	reader := bufio.NewReader(ptmx)

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			break
		}

		data := Data{
			Message: string(line),
			Service: serviceRunner.service,
		}
		message := Message{
			EventName: "log",
			Data:      data,
		}
		serviceRunner.runner.Logs = append(serviceRunner.runner.Logs, data)

		bytes, _ := json.Marshal(message)
		for _, user := range serviceRunner.runner.connections.Users {

			err = websockets.Send(user, string(bytes))
			if err != nil {
				fmt.Println(aurora.Red(err))
			}
		}

	}
}
