package chat

import "github.com/gorilla/websocket"

var users = make(map[string]*websocket.Conn)

func AddUser(username string, conn *websocket.Conn) {
	users[username] = conn
}

func SendPrivateMessage(fromUsername, toUsername string, message []byte) {
	conn, ok := users[toUsername]
	if !ok {
		return
	}
	conn.WriteMessage(websocket.TextMessage, message)
}
