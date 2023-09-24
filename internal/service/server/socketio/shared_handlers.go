package socketio

import (
	"log"

	"github.com/google/uuid"
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
	carId := s.CarConns.Disconnect(id)
	if carId != uuid.Nil {
		//do anything related to removing a car from the connections
	} else {
		userId := s.UserConns.Disconnect(id)
		if userId == uuid.Nil { //id wasn't a car or a user
			log.Printf("error: socketio conn %s was not a car or a user\n", id)
		} else {
			//do anything related to disconnecting a user
			s.CarConns.clearUser(userId)
		}
	}
}
