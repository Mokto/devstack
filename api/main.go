package main

import (
	"context"
	"devstack/config"
	"devstack/errors"
	"devstack/graceful"
	"devstack/runner"
	"devstack/ulimit"
	"devstack/utility"
	"devstack/websockets"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/logrusorgru/aurora"
	"golang.org/x/net/websocket"
)

var Prod string

//go:embed build
var content embed.FS

func clientHandler() http.Handler {
	fsys := fs.FS(content)
	contentStatic, _ := fs.Sub(fsys, "build")
	return http.FileServer(http.FS(contentStatic))
}

func main() {
	ulimit.Increase()

	isProd := Prod == "true"

	graceful := graceful.Start()

	connections := websockets.New()
	configFile, err := config.ReadConfigurationFile()
	if err != nil {
		panic(err)
	}

	servicesRunner := runner.Start(configFile, connections)

	graceful.OnShutdown(func() {
		servicesRunner.StopAll()
	})

	if isProd {
		go func() {
			mux := http.NewServeMux()
			mux.Handle("/", clientHandler())
			fmt.Println(aurora.Green("Server started on port :9999"))
			err := http.ListenAndServe(":9999", mux)
			if err != nil {
				panic(err)
			}
		}()

		go func() {
			time.Sleep(time.Second)
			utility.OpenBrowser("http://localhost:9999")
		}()
	}

	// API Router
	restAPI := newRestAPI(connections, configFile, servicesRunner)
	restAPI.RunServer(graceful, isProd)

	graceful.Wait()
}

// RestServer is the implementation of the rest API
type RestServer struct {
	echoServer *echo.Echo
}

type SetWatchingBody struct {
	IsWatching bool `json:"isWatching"`
}
type SetRunningBody struct {
	IsRunning bool `json:"isRunning"`
}

// NewRestAPI initialize an empty
func newRestAPI(connections *websockets.Connections, configFile *config.ConfigurationFile, servicesRunner *runner.Runner) *RestServer {

	e := echo.New()

	e.Use(middleware.CORS())
	e.HTTPErrorHandler = errors.HTTPErrorHandler
	e.Use(errors.PanicMiddleware)
	e.GET("/ws", websocketHandler(connections))
	e.GET("/healthcheck", func(c echo.Context) error {
		return c.String(http.StatusOK, "")
	})
	e.GET("/state", func(c echo.Context) error {
		for _, service := range configFile.Services {
			service.IsWatching = servicesRunner.IsWatching(service.Name)
			service.IsRunning = servicesRunner.IsRunning(service.Name)
		}
		return c.JSON(http.StatusOK, configFile)
	})
	e.GET("/logs", func(c echo.Context) error {
		// we only return the last 1000 logs
		minIndex := len(servicesRunner.Logs) - 1000
		if minIndex < 0 {
			minIndex = 0
		}
		return c.JSON(http.StatusOK, servicesRunner.Logs[minIndex:])
	})
	e.POST("/restart/:name", func(c echo.Context) error {
		serviceName := c.Param("name")
		servicesRunner.Restart(serviceName)
		return c.NoContent(http.StatusOK)
	})
	e.POST("/setWatching/:name", func(c echo.Context) error {
		serviceName := c.Param("name")
		var body SetWatchingBody
		if err := c.Bind(&body); err != nil {
			return errors.Wrap(err)
		}
		servicesRunner.SetWatching(serviceName, body.IsWatching)
		return c.NoContent(http.StatusOK)
	})

	e.POST("/setRunning/:name", func(c echo.Context) error {
		serviceName := c.Param("name")
		var body SetRunningBody
		if err := c.Bind(&body); err != nil {
			return errors.Wrap(err)
		}
		servicesRunner.SetIsRunning(serviceName, body.IsRunning)
		return c.NoContent(http.StatusOK)
	})

	return &RestServer{
		echoServer: e,
	}
}

// RunServer starts the rest server on a specific port
func (server *RestServer) RunServer(graceful *graceful.Manager, isProd bool) {
	go func() {
		server.echoServer.HideBanner = true
		server.echoServer.HidePort = isProd
		err := server.echoServer.Start(":9111")
		if err != nil && err.Error() != "http: Server closed" {
			// server.common.Exceptions.CaptureFatalException(err)
			fmt.Println(aurora.Red(err))

		}
	}()

	graceful.OnShutdown(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		if err := server.echoServer.Shutdown(ctx); err != nil {
			server.echoServer.Logger.Fatal(err)
		}
	})
}

func websocketHandler(connections *websockets.Connections) func(c echo.Context) error {
	return func(c echo.Context) error {
		websocket.Handler(func(ws *websocket.Conn) {
			defer ws.Close()
			for {
				var err error
				connections.Connect(ws)

				// Read
				msg := ""
				err = websocket.Message.Receive(ws, &msg)
				if err != nil {
					if err.Error() == "EOF" {
						connections.Disconnect(ws)
						return
					}
					panic(err)
				}

			}
		}).ServeHTTP(c.Response(), c.Request())
		return nil
	}

}
