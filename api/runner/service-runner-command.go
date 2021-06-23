package runner

import (
	"bufio"
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

		message := string(line)
		serviceRunner.runner.Logs = append(serviceRunner.runner.Logs, Data{
			Message: message,
			Service: serviceRunner.service,
		})
		serviceRunner.SendLog(message)
	}
}
