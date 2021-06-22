package websockets

import (
	"devstack/errors"

	"golang.org/x/net/websocket"
)

func Send(ws *websocket.Conn, message string) error {

	err := websocket.Message.Send(ws, message)
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}
