package controller

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Erro ao atualizar a conexão WebSocket:", err)
	}
	defer conn.Close()

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Erro ao ler a mensagem WebSocket:", err)
		}

		if messageType == websocket.TextMessage {
			err := conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				fmt.Fprintf(w, err.Error())
			}
		}
	}
}