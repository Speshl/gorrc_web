package models

import (
	"github.com/google/uuid"
	"github.com/pion/webrtc/v3"
)

type CarConnectReq struct {
	Key       string `json:"key"`
	Password  string `json:"password"`
	SeatCount int    `json:"seat_count"`
}

type CarConnectResp struct {
	Car   Car
	Track Track
}

type CarAnswerReq struct {
	Answer     *webrtc.SessionDescription `json:"answer"`
	SeatNumber int                        `json:"seat_number"`
}

type Car struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	ShortName string    `json:"short_name"`
	Type      string    `json:"type"`
	SeatCount int       `json:"seat_count"`
}

type Track struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	ShortName string    `json:"short_name"`
	Type      string    `json:"type"`
}
