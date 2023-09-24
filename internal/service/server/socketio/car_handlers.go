package socketio

import (
	"context"
	"log"

	"github.com/Speshl/gorrc_web/internal/service/server/models"
	"github.com/google/uuid"
	socketio "github.com/googollee/go-socket.io"
)

// func (s *SocketIOServer) OnConnect(socketConn socketio.Conn) error {
// 	log.Printf("socketio connected %s - Local: %s - Remote: %s\n", socketConn.ID(), socketConn.LocalAddr().String(), socketConn.RemoteAddr().String())
// 	// id := socketConn.ID()

// 	// Create a new Client for the connected socket
// 	// conn, err := s.NewClientConn(socketConn)
// 	// if err != nil {
// 	// 	return fmt.Errorf("failed creating new client: %w", err)
// 	// }

// 	// s.connectionsLock.Lock()
// 	// s.baseConnections[id] = conn
// 	// s.connectionsLock.Unlock()
// 	return nil
// }

// func (s *SocketIOServer) OnDisconnect(socketConn socketio.Conn, reason string) {
// 	log.Printf("socketio connection disconnected (%s): %s\n", reason, socketConn.ID())
// 	s.RemoveClient(socketConn.ID())
// }

// func (s *SocketIOServer) OnError(socketConn socketio.Conn, err error) {
// 	log.Println("Got Error")
// 	log.Printf("Error: %s", err.Error())
// 	// 	log.Printf("Id:%s", socketConn.ID())
// 	// 	log.Printf("socketio connection %s error: %s\n", socketConn.ID(), err.Error())
// }

// func (s *SocketIOServer) OnICECandidate(socketConn socketio.Conn, msg []byte) {
// 	log.Printf("candidate recieved from client: %s", socketConn.ID())
// }

//Custom events below

// TODO: Make sure replace current car conn with same id if exists
func (s *SocketIOServer) ConnectCar(socketConn socketio.Conn, msg string) {
	log.Printf("connect car event ID: %s", socketConn.ID())

	details := models.CarConnectReq{}
	err := decode(msg, &details)
	if err != nil {
		log.Printf("failed decoding connect car event msg - %s\n", err.Error())
		return
	}

	id, err := uuid.Parse(details.Key)
	if err != nil {
		log.Printf("invalid car uuid - %s\n", err.Error())
		return
	}

	car, err := s.store.GetCarByID(context.Background(), id)
	if err != nil {
		log.Printf("car uuid not found - %s\n", err.Error())
		return
	}

	track, err := s.store.GetTrackByID(context.Background(), car.Track)
	if err != nil {
		log.Printf("car uuid not found - %s\n", err.Error())
		return
	}

	key := s.CarConns.GetKeyByCarId(id) //See if connection for this car is already established
	if key != "" {
		log.Printf("old connection for car found (%s), removing old connection\n", key)
		s.CarConns.Disconnect(key)
	}

	s.CarConns.NewCarConnection(socketConn, car, track)

	log.Printf("car connected: %s(%s) @ %s(%s)\n", car.Name, car.ShortName, track.Name, track.ShortName)

	resp := models.CarConnectResp{
		Car: models.Car{
			Id:        car.Id,
			Name:      car.Name,
			ShortName: car.ShortName,
			Type:      car.Type,
		},
		Track: models.Track{
			Id:        track.Id,
			Name:      track.Name,
			ShortName: track.ShortName,
			Type:      track.Type,
		},
	}
	encodedResp, err := encode(resp)
	if err != nil {
		log.Printf("failed encoding car connect response(%s): %+v", socketConn.ID(), resp)
	}
	s.CarConns.emit(socketConn.ID(), "register_success", encodedResp)
}

func (s *SocketIOServer) HealthyCar(socketConn socketio.Conn, msg string) {
	s.CarConns.resetHealth(socketConn.ID())
	log.Printf("car conn %s reports healthy\n", socketConn.ID())
}

func (s *SocketIOServer) OnCarAnswer(socketConn socketio.Conn, msg string) {
	log.Println("Answer Recieved")
	//get name of this car

	car := s.CarConns.getCar(socketConn.ID())
	track := s.CarConns.getTrack(socketConn.ID())

	//figure out what user was requesting this car
	userKey := s.UserConns.GetKeyByCarId(car.Id)

	user := s.UserConns.getUser(userKey)

	s.CarConns.setUser(socketConn.ID(), user)

	//pass the answer through to the user
	s.UserConns.emit(userKey, "answer", msg)
	log.Printf("car %s @ %s sent answer to user %s\n", car.Name, track.Name, user.UserName)

}
