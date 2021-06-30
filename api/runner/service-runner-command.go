package runner

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/creack/pty"
	"github.com/logrusorgru/aurora"
)

func getEnvValue(key string, value string) string {
	r, _ := regexp.Compile(`{{\$([A-Z]+)}}`)
	matches := r.FindStringSubmatch(value)
	if len(matches) > 1 {
		value = os.Getenv(matches[1])
	}
	return key + "=" + value
}

func (serviceRunner *ServiceRunner) execCommand() {
	if serviceRunner.IsRunning {
		serviceRunner.SendLog(aurora.Yellow("Waiting for old process to be killed...").String())
		time.Sleep(time.Millisecond * 100)
		serviceRunner.execCommand()
		return
	}

	serviceRunner.SendLog(aurora.Yellow("Starting process...").String())
	splitted := strings.Split(serviceRunner.service.Command, " ")
	serviceRunner.cmd = exec.Command(splitted[0], splitted[1:]...)
	serviceRunner.IsRunning = true
	serviceRunner.SendIsRunning(true)
	// serviceRunner.cmd = exec.Command("/bin/bash", "-c", serviceRunner.service.Command)

	if serviceRunner.service.Cwd != "" {
		serviceRunner.cmd.Dir = serviceRunner.runner.config.BasePath + "/" + serviceRunner.service.Cwd
	}
	if serviceRunner.service.Env != nil || serviceRunner.runner.config.Env != nil {
		envs := []string{}
		for key, value := range serviceRunner.service.Env {
			envs = append(envs, getEnvValue(key, value))
		}
		for key, value := range serviceRunner.runner.config.Env {
			envs = append(envs, getEnvValue(key, value))
		}

		serviceRunner.cmd.Env = append(os.Environ(), envs...)
	}

	ptmx, err := pty.Start(serviceRunner.cmd)
	if err != nil {
		fmt.Println(aurora.Red("Panic on service " + serviceRunner.service.Name + ". Command is " + serviceRunner.service.Command))
		panic(err)
	}
	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.

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

		message := string(line)
		// if serviceRunner.service.Name == "segmentation" {
		// 	fmt.Println(message)
		// }
		serviceRunner.runner.Logs = append(serviceRunner.runner.Logs, Data{
			Message: message,
			Service: serviceRunner.service,
		})
		serviceRunner.SendLog(message)
	}
	serviceRunner.SendIsRunning(false)
	serviceRunner.IsRunning = false

}
