package connection

import (
	"context"
	"fmt"

	"github.com/gorilla/websocket"
)

type Connection struct {
	wsConn     *websocket.Conn
	inChan     chan []byte
	outChan    chan []byte
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func New(wsConn *websocket.Conn) (*Connection, error) {
	ctx, cancel := context.WithCancel(context.Background())

	conn := &Connection{
		wsConn:     wsConn,
		inChan:     make(chan []byte, 1000),
		outChan:    make(chan []byte, 1000),
		ctx:        ctx,
		cancelFunc: cancel,
	}

	go conn.readLoop()
	go conn.writeLoop()

	return conn, nil
}

func (conn *Connection) ReadMessage() ([]byte, error) {
	var data []byte
	select {
	case data = <-conn.inChan:
	case <-conn.ctx.Done():
		return data, fmt.Errorf("connection is closed")
	}
	return data, nil
}

func (conn *Connection) WriteMessage(data []byte) error {
	select {
	case conn.outChan <- data:
	case <-conn.ctx.Done():
		return fmt.Errorf("connection is closed")
	}
	return nil
}

func (conn *Connection) Close() {
	conn.wsConn.Close()
	conn.cancelFunc()
}

func (conn *Connection) readLoop() {
	for {
		_, data, err := conn.wsConn.ReadMessage()
		if err != nil {
			conn.Close()
			return
		}

		select {
		case conn.inChan <- data:
		case <-conn.ctx.Done():
			return
		}
	}
}

func (conn *Connection) writeLoop() {
	for {
		select {
		case data := <-conn.outChan:
			err := conn.wsConn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				conn.Close()
				return
			}
		case <-conn.ctx.Done():
			conn.Close()
			return
		}
	}
}
