package runner

import (
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/boz/go-throttle"
	"github.com/fsnotify/fsnotify"
	"github.com/logrusorgru/aurora"
)

func (serviceRunner *ServiceRunner) watch() {
	if serviceRunner.IsWatching {
		return
	}
	serviceRunner.IsWatching = true

	serviceRunner.watcher, _ = fsnotify.NewWatcher()
	defer serviceRunner.watcher.Close()

	if err := filepath.Walk("../example", func(path string, info fs.FileInfo, err error) error {
		if info.Mode().IsDir() {
			return serviceRunner.watcher.Add(path)
		}
		return nil
	}); err != nil {
		panic(err)
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
				panic(err)
			}
		}
	}()

	<-serviceRunner.stopWatchingChannel
	serviceRunner.IsWatching = false
	throttle.Stop()

}

func (serviceRunner *ServiceRunner) stopWatching() {
	serviceRunner.stopWatchingChannel <- true
}
