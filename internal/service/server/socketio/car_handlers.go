package socketio

import (
	"context"
	"log"

	"github.com/Speshl/gorrc_web/internal/service/server/models"
	"github.com/google/uuid"
	socketio "github.com/googollee/go-socket.io"
)

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

	if details.SeatCount < 1 || details.SeatCount > car.SeatCount {
		log.Printf("car %s requested connection with more seats than it supports: %d<%d\n", car.ShortName, car.SeatCount, details.SeatCount)
		return
	}

	track, err := s.store.GetTrackByID(context.Background(), car.Track)
	if err != nil {
		log.Printf("car uuid not found - %s\n", err.Error())
		return
	}

	key, err := s.CarConns.GetKeyByCarId(id) //See if connection for this car is already established
	if err == nil {
		log.Printf("old connection for car found (%s), removing old connection\n", key)
		s.CarConns.Disconnect(key)
	}

	s.CarConns.NewCarConnection(socketConn, car, track, details.Password)

	log.Printf("car connected: %s(%s) @ %s(%s) with %d seats\n", car.Name, car.ShortName, track.Name, track.ShortName, details.SeatCount)

	resp := models.CarConnectResp{
		Car: models.Car{
			Id:        car.Id,
			Name:      car.Name,
			ShortName: car.ShortName,
			Type:      car.Type,
			SeatCount: details.SeatCount,
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
	s.CarConns.Emit(socketConn.ID(), "register_success", encodedResp)
}

func (s *SocketIOServer) HealthyCar(socketConn socketio.Conn, msg string) {
	s.CarConns.ResetHealth(socketConn.ID())
	log.Printf("car conn %s reports healthy\n", socketConn.ID())
}

func (s *SocketIOServer) OnCarAnswer(socketConn socketio.Conn, msg string) {
	log.Println("Answer Recieved")

	answerReq := models.CarAnswerReq{}
	err := decode(msg, &answerReq)
	if err != nil {
		log.Printf("failed decoding answer request - %s\n", err.Error())
		return
	}

	car, err := s.CarConns.GetCar(socketConn.ID())
	if err != nil {
		log.Printf("car not found for socket connection: %s\n", socketConn.ID())
		return
	}

	track, err := s.CarConns.GetTrack(socketConn.ID())
	if err != nil {
		log.Printf("track not found for socket connection: %s\n", socketConn.ID())
		return
	}

	//figure out what user was requesting this car
	userKey, err := s.UserConns.GetKeyByCarIdAndSeat(car.Id, answerReq.SeatNumber)
	if err != nil {
		log.Printf("no user fround for car %s in seat %d\n", car.Name, answerReq.SeatNumber)
		return
	}

	user, err := s.UserConns.GetUser(userKey)
	if err != nil {
		log.Printf("no user found for user key %s\n", userKey)
		return
	}

	s.CarConns.SetUserInSeat(socketConn.ID(), answerReq.SeatNumber, user)

	//pass the answer through to the user
	encodedResp, err := encode(answerReq.Answer)
	if err != nil {
		log.Printf("failed encoding car connect response(%s): %+v", socketConn.ID(), encodedResp)
	}

	s.UserConns.Emit(userKey, "answer", encodedResp)
	log.Printf("car %s @ %s seat %d sent answer to user %s\n", car.Name, track.Name, answerReq.SeatNumber, user.UserName)

}
