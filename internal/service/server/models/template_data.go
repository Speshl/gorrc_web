package models

import (
	"html/template"
	"strings"
)

type TrackListTMPLData struct {
	Tracks []TrackTMPLData
}

type TrackTMPLData struct {
	Name           string
	ShortName      string
	Type           string
	Logo           string
	Description    string
	CarCount       int
	UserCount      int
	SpectatorCount int
}

func (a TrackListTMPLData) Len() int      { return len(a.Tracks) }
func (a TrackListTMPLData) Swap(i, j int) { a.Tracks[i], a.Tracks[j] = a.Tracks[j], a.Tracks[i] }
func (a TrackListTMPLData) Less(i, j int) bool {
	if a.Tracks[i].UserCount > a.Tracks[j].UserCount {
		return true
	} else if a.Tracks[i].UserCount == a.Tracks[j].UserCount {
		if a.Tracks[i].CarCount > a.Tracks[j].CarCount {
			return true
		} else if a.Tracks[i].CarCount == a.Tracks[j].CarCount {
			return strings.Compare(a.Tracks[i].Name, a.Tracks[j].Name) == -1
		}
	}
	return false
}

type CarListTMPLData struct {
	Cars []CarTMPLData
}

type CarTMPLData struct {
	UserId         string
	Name           string
	ShortName      string
	Type           string
	Logo           string
	Description    string
	TrackShortName string
	HasPassword    bool
	Status         string
}

func (a CarListTMPLData) Len() int      { return len(a.Cars) }
func (a CarListTMPLData) Swap(i, j int) { a.Cars[i], a.Cars[j] = a.Cars[j], a.Cars[i] }
func (a CarListTMPLData) Less(i, j int) bool {
	if StatusToInt(a.Cars[i].Status) < StatusToInt(a.Cars[j].Status) {
		return true
	} else if StatusToInt(a.Cars[i].Status) == StatusToInt(a.Cars[j].Status) {
		return strings.Compare(a.Cars[i].Name, a.Cars[j].Name) == -1
	}
	return false
}

type DriveButtonTMPLData struct {
	TrackName      string
	TrackShortName string
	CarName        string
	CarShortName   string
}

type VideoTMPLData struct {
	UserId    string
	TrackName string
	CarName   string
}

type LoginTMPLData struct {
	LoginError string
}

type RegisterTMPLData struct {
	DisplayNameError string
	PasswordError    string
	UserName         string
	RealName         string
	Email            string
}

type MainTMPLData struct {
	Body template.HTML
}

func StatusToInt(status string) int {
	switch status {
	case "available":
		return 0
	case "occupied":
		return 1
	default:
		return 2
	}
}
