package main

import (
	"devstack/websockets"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/logrusorgru/aurora"
	"golang.org/x/net/websocket"
)

func main() {
	connections := websockets.New()

	// API Router
	restAPI := newRestAPI(connections)
	restAPI.RunServer()
}

// RestServer is the implementation of the rest API
type RestServer struct {
	echoServer *echo.Echo
}

// NewRestAPI initialize an empty
func newRestAPI(connections *websockets.Connections) *RestServer {

	e := echo.New()

	// e.Use(common.Exceptions.PanicMiddlewareWebsocket)

	e.GET("/ws", websocketHandler(connections))
	e.GET("/healthcheck", func(c echo.Context) error {
		return c.String(http.StatusOK, "")
	})

	return &RestServer{
		echoServer: e,
	}
}

// RunServer starts the rest server on a specific port
func (server *RestServer) RunServer() {
	// go func() {
	err := server.echoServer.Start(":9111")
	if err != nil && err.Error() != "http: Server closed" {
		// server.common.Exceptions.CaptureFatalException(err)
		fmt.Println(aurora.Red(err))

	}
	// }()

	// server.common.Graceful.OnShutdown(func() {
	// 	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	// 	defer cancel()
	// 	if err := server.echoServer.Shutdown(ctx); err != nil {
	// 		server.echoServer.Logger.Fatal(err)
	// 	}
	// })
}

func websocketHandler(connections *websockets.Connections) func(c echo.Context) error {
	return func(c echo.Context) error {
		websocket.Handler(func(ws *websocket.Conn) {
			defer ws.Close()
			for {
				var err error

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

				connections.Connect(ws)

				// var event connectionEvent
				// err = json.Unmarshal([]byte(msg), &event)
				// if err != nil {
				// 	common.Exceptions.CaptureException(err)
				// 	return
				// }

				// if event.Event == "connected" {
				// 	userID =
				// }

				// if event.Event == "disconnected" {
				// 	connections.LoggedOut(ws, userID)
				// }
			}
		}).ServeHTTP(c.Response(), c.Request())
		return nil
	}

}
