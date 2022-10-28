package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/mneumi/websocket-demo/connection"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandleFunc(w http.ResponseWriter, r *http.Request) {
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		wsConn.Close()
		return
	}
	a := 10

	fmt.Print(a)

	a = 20
	a = 30
	a = 40
	a = 50

	conn, err := connection.New(wsConn)
	if err != nil {
		wsConn.Close()
		return
	}
	defer wsConn.Close()

	for {
		data, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			return
		}
		fmt.Printf("%s\n", data)
		conn.WriteMessage(data)
	}
}

func main() {
	http.HandleFunc("/ws", wsHandleFunc)
	a := 10
	fmt.Println("Server Start")
	a = 20
	fmt.Print(a)
	a = 30
	fmt.Print(a)

	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
