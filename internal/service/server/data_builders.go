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
	userTrackIds := s.socketIOServer.UserConns.GetDataForActiveTracks()

	//get a map of tracks with connected cars
	carTrackIds := s.socketIOServer.CarConns.GetDataForActiveTracks()

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
			TrackData: models.TrackData{
				UserCount: userTrackIds[allTracks[i].Id].UserCount,
				CarCount:  carTrackIds[allTracks[i].Id].CarCount,
				SeatCount: carTrackIds[allTracks[i].Id].SeatCount,
			},
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

	//build array with all the data
	sortCarList := make([]models.CarTMPLData, 0, len(trackCars))
	for i := range trackCars {

		//check which cars are currently occupied
		seatStatuses, err := s.socketIOServer.CarConns.GetSeatStatusByCarId(trackCars[i].Id)
		if err != nil {
			//log.Printf("error: failed getting seat status for car %s, skipping...\n", trackCars[i].Name)
			continue
		}

		carKey, err := s.socketIOServer.CarConns.GetKeyByCarId(trackCars[i].Id)
		if err != nil {
			continue
		}

		password, err := s.socketIOServer.CarConns.GetPassword(carKey)
		if err != nil {
			continue
		}

		hasPassword := false
		if password != "" {
			hasPassword = true
		}

		//Add an entry in the car list per seat
		for j := range seatStatuses {
			seatType := "driver"
			if j > 0 {
				seatType = "passenger"
			}

			carToAdd := models.CarTMPLData{
				Name:         trackCars[i].Name,
				CarShortName: trackCars[i].ShortName,
				Type:         trackCars[i].Type,
				//Logo:        allTracks[i].Logo,
				Logo:           "./img/golang_drive.png", //TODO Dynamic logo
				Description:    trackCars[i].Description,
				HasPassword:    hasPassword,
				TrackShortName: track.ShortName,
				SeatNumber:     j,
				SeatStatus:     seatStatuses[j],
				SeatType:       seatType,
			}

			sortCarList = append(sortCarList, carToAdd)
		}

	}
	returnValue.Cars = sortCarList

	sort.Sort(returnValue)
	return returnValue, nil
}
