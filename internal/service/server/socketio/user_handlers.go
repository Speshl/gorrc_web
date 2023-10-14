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
		log.Printf("error: failed decoding connect user event msg - %s\n", err.Error())
		return
	}

	userName, err := auth.ValidateJWT(details.Token)
	if err != nil {
		log.Printf("error: failed validating token on socket connect user: %s\n", err.Error())
		return
	}

	user, err := s.store.GetUserByUserName(context.Background(), userName)
	if err != nil {
		log.Printf("error: failed getting user from database - %s\n", err.Error())
		return
	}

	log.Printf("user %s connected", user.UserName)

	s.UserConns.NewUserConnection(socketConn, user)
	s.UserConns.Emit(socketConn.ID(), "connected")
}

func (s *SocketIOServer) HealthyUser(socketConn socketio.Conn, msg string) {
	health := ""
	err := decode(msg, &health)
	if err != nil {
		log.Printf("error: failed decoding healthy user event msg - %s\n", err.Error())
		return
	}
	log.Printf("%s user health: %s\n", socketConn.ID(), health)
}

func (s *SocketIOServer) OnUserOffer(socketConn socketio.Conn, msg string) {
	log.Printf("user offer recieved: %s\n", socketConn.ID())
	offer := models.Offer{}
	err := decode(msg, &offer)
	if err != nil {
		log.Printf("error: offer from %s failed unmarshaling: %s\n", socketConn.ID(), err.Error())
		return
	}

	//check if the car exists
	carKey, err := s.CarConns.GetKeyByCarShortName(offer.CarShortName)
	if err != nil {
		log.Printf("error: car not found with shortname %s\n", offer.CarShortName)
		return
	}

	car, err := s.CarConns.GetCar(carKey)
	if err != nil {
		log.Printf("error: no car found for car shortname %s with carKey %s\n", offer.CarShortName, carKey)
		return
	}

	track, err := s.CarConns.GetTrack(carKey)
	if err != nil {
		log.Printf("error: no track found for car shortname %s with carKey %s\n", offer.CarShortName, carKey)
		return
	}

	user, err := s.UserConns.GetUser(socketConn.ID())
	if err != nil {
		log.Printf("error: no connection found for user: %s\n", socketConn.ID())
		return
	}

	err = s.UserConns.SetCar(socketConn.ID(), car, offer.SeatNum)
	if err != nil {
		log.Printf("error: failed setting car for user: %s\n", err.Error())
		return
	}

	err = s.UserConns.SetTrack(socketConn.ID(), track)
	if err != nil {
		log.Printf("error: failed setting track for user: %s\n", err.Error())
		return
	}

	encodedOffer, err := encode(offer)
	if err != nil {
		log.Printf("error: failed encoding offer: %s\n", err.Error())
		return
	}

	log.Printf("user %s sending offer to car %s @ %s seat %d\n", user.UserName, car.Name, track.Name, offer.SeatNum)
	s.CarConns.Emit(carKey, "offer", encodedOffer)
}

func (s *SocketIOServer) OnUserCandidate(socketConn socketio.Conn, msg string) {
	log.Printf("User Candidate Recieved (%s): %s\n", socketConn.ID(), msg)
}
