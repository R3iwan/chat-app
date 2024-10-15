package chat

import "github.com/gorilla/websocket"

type Room struct {
	Name      string
	Users     map[*websocket.Conn]bool
	Broadcast chan []byte
}

func CreateRoom(name string) *Room {
	return &Room{
		Name:      name,
		Users:     make(map[*websocket.Conn]bool),
		Broadcast: make(chan []byte),
	}
}

func (r *Room) JoinRoom(conn *websocket.Conn) {
	r.Users[conn] = true
}

func (r *Room) LeaveRoom(conn *websocket.Conn) {
	delete(r.Users, conn)
	conn.Close()
}

func (r *Room) BroadcastMessage(message []byte) {
	for user := range r.Users {
		user.WriteMessage(websocket.TextMessage, message)
	}
}
