package models

import (
	"github.com/google/uuid"
	"github.com/pion/webrtc/v3"
)

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type RegisterBody struct {
	Username       string `json:"user_name"`
	Password       string `json:"password"`
	VerifyPassword string `json:"verify_password"`
	RealName       string `json:"real_name"`
	Email          string `json:"email"`
}

type UserConnect struct {
	CarShortName   string `json:"car_short_name"`
	TrackShortName string `json:"track_short_name"`
	Token          string `json:"token"`
}

type IceCandidate struct {
	Candidate    webrtc.ICECandidateInit `json:"candidate"`
	CarShortName string                  `json:"car_name"`
	SeatNum      int                     `json:"seat_number"`
	UserId       uuid.UUID               `json:"user_id"`
}

type Offer struct {
	Offer        webrtc.SessionDescription `json:"offer"`
	CarShortName string                    `json:"car_name"`
	SeatNum      int                       `json:"seat_number"`
	UserId       uuid.UUID                 `json:"user_id"`
}

type DriveReqBody struct {
	TrackShortName string `json:"track_short_name"`
	CarShortName   string `json:"car_short_name"`
	SeatNumber     string `json:"seat_number"`
	Password       string `json:"password"`
}
