package socketio

import (
	"log"

	socketio "github.com/googollee/go-socket.io"
)

func (s *SocketIOServer) onConnect(socketConn socketio.Conn) error {
	log.Printf("socketio connected %s - Local: %s - Remote: %s\n", socketConn.ID(), socketConn.LocalAddr().String(), socketConn.RemoteAddr().String())
	return nil
}

func (s *SocketIOServer) onDisconnect(socketConn socketio.Conn, reason string) {
	log.Printf("socketio connection disconnected (%s): %s\n", socketConn.ID(), reason)
	s.RemoveClient(socketConn.ID())
}

func (s *SocketIOServer) onError(socketConn socketio.Conn, err error) {
	log.Printf("socketio connection %s error: %s\n", socketConn.ID(), err.Error())
}

func (s *SocketIOServer) RemoveClient(id string) {
	log.Printf("Removing Client: %s\n", id)
	//Should only match one of the types of connections, but not sure which so check both
	_, err := s.CarConns.Disconnect(id)
	if err == nil {
		log.Printf("socketio conn %s was a car and was removed\n", id)
		return
	}

	userId, err := s.UserConns.Disconnect(id)
	if err == nil {
		log.Printf("socketio conn %s was a user and was removed\n", id)
		s.CarConns.ClearUser(userId)
		return
	}

	log.Printf("error: socketio conn %s was neither a car or a user\n", id)
}
