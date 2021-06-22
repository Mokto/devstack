package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/logrusorgru/aurora"
)

func main() {
	// // API Router
	// restAPI := newRestAPI()
	// restAPI.RunServer()
	index := 0
	for {
		time.Sleep(time.Second)
		fmt.Println(aurora.Green(strconv.Itoa(index)))
		index++
	}
}

// RestServer is the implementation of the rest API
type RestServer struct {
	echoServer *echo.Echo
}

// NewRestAPI initialize an empty
func newRestAPI() *RestServer {

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Test")
	})

	return &RestServer{
		echoServer: e,
	}
}

// RunServer starts the rest server on a specific port
func (server *RestServer) RunServer() {
	// go func() {
	err := server.echoServer.Start(":9112")
	if err != nil && err.Error() != "http: Server closed" {
		// server.common.Exceptions.CaptureFatalException(err)
		fmt.Println(aurora.Red(err))
	}
}
