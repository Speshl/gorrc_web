package socketio

import (
	"context"
	"log"

	"github.com/Speshl/gorrc_web/internal/service/auth"
	"github.com/Speshl/gorrc_web/internal/service/server/models"
	socketio "github.com/googollee/go-socket.io"
)

func (s *SocketIOServer) ConnectUser(socketConn socketio.Conn, msg string) {
	log.Printf("connect user event: %s", socketConn.ID())

	details := models.UserConnect{}
	err := decode(msg, &details)
	if err != nil {
		log.Printf("failed decoding connect user event msg - %s\n", err.Error())
		return
	}

	userName, err := auth.ValidateJWT(details.Token)
	if err != nil {
		log.Printf("failed validating token on socket connect user: %s\n", err.Error())
		return
	}

	user, err := s.store.GetUserByUserName(context.Background(), userName)
	if err != nil {
		log.Printf("error getting user from database - %s\n", err.Error())
		return
	}

	log.Printf("user %s connected", user.UserName)

	s.UserConns.NewUserConnection(socketConn, user)
	s.UserConns.emit(socketConn.ID(), "connected")
}

func (s *SocketIOServer) HealthyUser(socketConn socketio.Conn, msg string) {
	health := ""
	err := decode(msg, &health)
	if err != nil {
		log.Printf("failed decoding healthy user event msg - %s\n", err.Error())
		return
	}
	log.Printf("%s user health: %s\n", socketConn.ID(), health)
}

func (s *SocketIOServer) OnUserOffer(socketConn socketio.Conn, msg string) {
	log.Printf("user offer recieved: %s\n", socketConn.ID())
	offer := models.Offer{}
	err := decode(msg, &offer)
	if err != nil {
		log.Printf("offer from %s failed unmarshaling: %s\n", socketConn.ID(), string(msg))
		return
	}

	//check if the car exists
	carKey := s.CarConns.GetKeyByCarShortName(offer.CarShortName)
	car := s.CarConns.getCar(carKey)
	if car == nil {
		log.Printf("no car found for carshortname %s with carKey %s\n", offer.CarShortName, carKey)
		return
	}

	track := s.CarConns.getTrack(carKey)
	if track == nil {
		log.Printf("no track found for carshortname %s with carKey %s\n", offer.CarShortName, carKey)
		return
	}

	user := s.UserConns.getUser(socketConn.ID())
	if user == nil {
		log.Printf("no user found in userConns: %s\n", socketConn.ID())
		return
	}
	s.UserConns.setCar(socketConn.ID(), car)
	s.UserConns.setTrack(socketConn.ID(), track)

	encodedOffer, err := encode(offer)
	if err != nil {
		log.Printf("failed encoding offer: %s\n", err.Error())
		return
	}

	log.Printf("user %s sending offer to car %s @ %s\n", user.UserName, car.Name, track.Name)
	s.CarConns.emit(carKey, "offer", encodedOffer)
	log.Printf("user %s sent offer to car %s @ %s\n", user.UserName, car.Name, track.Name)
}

func (s *SocketIOServer) OnUserCandidate(socketConn socketio.Conn, msg string) {
	log.Printf("User Candidate Recieved (%s): %s\n", socketConn.ID(), msg)
	// carOffer := CarOffer{}
	// err := decode(msg, &carOffer)
	// if err != nil {
	// 	log.Printf("Offer from %s failed unmarshaling: %s\n", socketConn.ID(), string(msg))
	// 	return
	// }

	// //check if the car exists
	// carConn := s.getCarConnByName(carOffer.CarName)

	// s.userLock.Lock()
	// userConn, ok := s.userConnections[socketConn.ID()]
	// if !ok {
	// 	log.Printf("user for id %s is missing from user connections\n", socketConn.ID())
	// 	return
	// }
	// userConn.RequestedCar = carOffer.CarName
	// s.userLock.Unlock()

	// log.Printf("Sending offer to car: %s\n", carConn.Connection.Socket.ID())
	// carConn.Connection.Socket.Emit("offer", msg)
}
