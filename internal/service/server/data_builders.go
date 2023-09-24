package server

import (
	"context"
	"fmt"
	"sort"

	"github.com/Speshl/gorrc_web/internal/service/server/models"
	"github.com/google/uuid"
)

func (s *Server) GetTrackList(ctx context.Context) (models.TrackListTMPLData, error) {
	returnValue := models.TrackListTMPLData{}

	//Get all the tracks from the database
	allTracks, err := s.store.GetTracks(ctx)
	if err != nil {
		return returnValue, fmt.Errorf("error getting tracks: %w", err)
	}

	//Get all the track IDs that currently have connected users
	userTrackIds := s.socketIOServer.UserConns.GetActiveTracks()

	//get a map of tracks with connected cars
	carTrackIds := s.socketIOServer.CarConns.GetActiveTracks()

	//build array with all the data
	sortTrackList := make([]models.TrackTMPLData, 0, len(allTracks))
	for i := range allTracks {
		sortTrackList = append(sortTrackList, models.TrackTMPLData{
			Name:      allTracks[i].Name,
			ShortName: allTracks[i].ShortName,
			Type:      allTracks[i].Type,
			//Logo:        allTracks[i].Logo,
			Logo:        "./img/golang_drive.png", //TODO Dynamic logo
			Description: allTracks[i].Description,
			UserCount:   userTrackIds[allTracks[i].Id],
			CarCount:    carTrackIds[allTracks[i].Id],
		})
	}

	trackList := models.TrackListTMPLData{
		Tracks: sortTrackList,
	}

	sort.Sort(trackList)

	return trackList, nil
}

func (s *Server) GetCarListForTrack(ctx context.Context, trackId uuid.UUID) (models.CarListTMPLData, error) {
	returnValue := models.CarListTMPLData{}

	//get all cars registered to that track
	trackCars, err := s.store.GetCarsByTrack(ctx, trackId)
	if err != nil {
		return returnValue, fmt.Errorf("error getting cars: %w", err)
	}

	track, err := s.store.GetTrackByID(ctx, trackId)
	if err != nil {
		return returnValue, fmt.Errorf("error getting track: %w", err)
	}

	//check which cars are currently connected
	activeCarKeys := s.socketIOServer.CarConns.GetActiveCarIds()

	//check which cars are currently occupied
	occupiedCarKeys := s.socketIOServer.CarConns.GetOccupiedCarIds()

	//build array with all the data
	sortCarList := make([]models.CarTMPLData, 0, len(trackCars))
	for i := range trackCars {
		carToAdd := models.CarTMPLData{
			Name:      trackCars[i].Name,
			ShortName: trackCars[i].ShortName,
			Type:      trackCars[i].Type,
			//Logo:        allTracks[i].Logo,
			Logo:           "./img/golang_drive.png", //TODO Dynamic logo
			Description:    trackCars[i].Description,
			HasPassword:    trackCars[i].Password.Valid,
			TrackShortName: track.ShortName,
		}

		if contains(occupiedCarKeys, trackCars[i].Id) {
			carToAdd.Status = "occupied"
		} else if contains(activeCarKeys, trackCars[i].Id) {
			carToAdd.Status = "available"
		} else {
			carToAdd.Status = "unavailable"
		}
		sortCarList = append(sortCarList, carToAdd)
	}
	returnValue.Cars = sortCarList

	sort.Sort(returnValue)
	return returnValue, nil
}

func contains(s []uuid.UUID, e uuid.UUID) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
