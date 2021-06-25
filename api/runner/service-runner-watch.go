package runner

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	p "path"

	"github.com/boz/go-throttle"
	"github.com/fsnotify/fsnotify"
	"github.com/logrusorgru/aurora"
)

func (serviceRunner *ServiceRunner) watch() {
	projectPath := os.Getenv("PROJECT_PATH")

	if serviceRunner.IsWatching {
		return
	}
	serviceRunner.IsWatching = true

	serviceRunner.watcher, _ = fsnotify.NewWatcher()
	defer serviceRunner.watcher.Close()

	for _, path := range serviceRunner.service.WatchDirectories {
		filePath := p.Clean(serviceRunner.service.Cwd + "/" + path)
		if projectPath != "" {
			filePath = p.Clean(projectPath + "/" + filePath)
		}
		if err := filepath.Walk(filePath, func(path string, info fs.FileInfo, err error) error {
			if info != nil && info.Mode().IsDir() {
				err := serviceRunner.watcher.Add(path)
				if err != nil {
					fmt.Println(aurora.Red(err), path)
				}
				return nil
			}
			return nil
		}); err != nil {
			fmt.Println(aurora.Red(filePath))
			panic(err)
		}
	}

	serviceRunner.stopWatchingChannel = make(chan bool)
	throttle := throttle.ThrottleFunc(time.Millisecond*100, true, func() {
		serviceRunner.SendLog(aurora.Yellow("Changed a file. Reloading...").String())
		serviceRunner.Restart()
	})

	go func() {
		for {
			select {
			// watch for events
			case event := <-serviceRunner.watcher.Events:
				if event.Op == fsnotify.Create {
					fi, err := os.Stat(event.Name)
					if err != nil {
						panic(err)
					}
					if fi.IsDir() {
						serviceRunner.watcher.Add(event.Name)
					}
				}
				if event.Op != fsnotify.Chmod {
					throttle.Trigger()
				}

				// watch for errors
			case err := <-serviceRunner.watcher.Errors:
				if err == nil {
					return
				}
				panic(err)
			}
		}
	}()

	<-serviceRunner.stopWatchingChannel

	serviceRunner.watcher.Close()
	serviceRunner.IsWatching = false
	throttle.Stop()

}

func (serviceRunner *ServiceRunner) StopWatching() {
	if serviceRunner.IsWatching {
		serviceRunner.stopWatchingChannel <- true
	}
}
