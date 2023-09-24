package models

import "github.com/google/uuid"

type CarConnectReq struct {
	Key      string `json:"key"`
	Password string `json:"password"`
}

type CarConnectResp struct {
	Car   Car
	Track Track
}
type Car struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	ShortName string    `json:"short_name"`
	Type      string    `json:"type"`
}

type Track struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	ShortName string    `json:"short_name"`
	Type      string    `json:"type"`
}
