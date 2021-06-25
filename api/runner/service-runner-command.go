package runner

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

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
	// Handle pty size.
	// ch := make(chan os.Signal, 1)
	// signal.Notify(ch, syscall.SIGWINCH)
	// go func() {
	// 	for range ch {
	// 		if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
	// 			log.Printf("error resizing pty: %s", err)
	// 		}
	// 	}
	// }()
	// ch <- syscall.SIGWINCH                        // Initial resize.
	// defer func() { signal.Stop(ch); close(ch) }() // Cleanup signals when done.
	// // Set stdin in raw mode.
	// oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	// if err != nil {
	// 	panic(err)
	// }
	// defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.
	// // Copy stdin to the pty and the pty to stdout.
	// // NOTE: The goroutine will keep reading until the next keystroke before returning.
	// go func() { _, _ = io.Copy(ptmx, os.Stdin) }()

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
