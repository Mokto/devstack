package websockets

import (
	"golang.org/x/net/websocket"
)

// Connections contains all currently logged in clients
type Connections struct {
	Users []*websocket.Conn
}

// New Initialize empty connections
func New() *Connections {
	return &Connections{
		Users: []*websocket.Conn{},
	}
}

// LoggedOut manages the loggedout state
func (connections *Connections) Connect(ws *websocket.Conn) {
	connections.cleanPreviousConnection(ws)

	connections.Users = append(connections.Users, ws)
	// connections.debug("logged out")
}

// Disconnect manages the disconnect state
func (connections *Connections) Disconnect(ws *websocket.Conn) {
	connections.cleanPreviousConnection(ws)
	// connections.debug("disconnect")
}

func (connections *Connections) cleanPreviousConnection(ws *websocket.Conn) {
	connections.Users = removeConn(connections.Users, ws)
}

// func (connections *Connections) debug(category string) {
// 	fmt.Println(category, connections.Users, connections.LoggedOutUsers)
// }
