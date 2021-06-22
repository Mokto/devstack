package websockets

import (
	"golang.org/x/net/websocket"
)

func removeConn(connections []*websocket.Conn, conn *websocket.Conn) []*websocket.Conn {
	index := -1
	for i, connection := range connections {
		if connection == conn {
			index = i
		}
	}

	if index != -1 {
		connections[index] = connections[len(connections)-1]
		connections[len(connections)-1] = nil
		connections = connections[:len(connections)-1]
	}

	return connections
}
