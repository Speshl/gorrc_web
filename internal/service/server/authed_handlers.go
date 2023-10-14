package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Speshl/gorrc_web/internal/service/auth"
	"github.com/Speshl/gorrc_web/internal/service/server/models"
)

// TODO: Ensure user is logged in before processing these requests
func (s *Server) trackListHandler(w http.ResponseWriter, req *http.Request) {

	_, err := auth.ValidateCookie(req)
	if err != nil {
		log.Printf("error: cookie failed validation on track list: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	trackListData, err := s.GetTrackList(req.Context())
	if err != nil {
		log.Printf("error: failed generating track list: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.templates["track_list"].Execute(w, trackListData)
	if err != nil {
		log.Printf("error: failed executing template: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) trackSelectHandler(w http.ResponseWriter, req *http.Request) {
	_, err := auth.ValidateCookie(req)
	if err != nil {
		log.Printf("error: cookie failed validation on track select: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	selectedTrack := req.FormValue("track")

	err = s.templates["track_select"].Execute(w, models.TrackTMPLData{
		ShortName: selectedTrack,
	})
	if err != nil {
		log.Printf("error: failed executing template: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) carListHandler(w http.ResponseWriter, req *http.Request) {
	_, err := auth.ValidateCookie(req)
	if err != nil {
		log.Printf("error: cookie failed validation on car list: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	track, err := s.store.GetTrackByShortName(req.Context(), req.FormValue("track"))
	if err != nil {
		log.Printf("error: failed getting track by short name: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	carData, err := s.GetCarListForTrack(req.Context(), track.Id)
	if err != nil {
		log.Printf("error: failed generating car list: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.templates["car_list"].Execute(w, carData)
	if err != nil {
		log.Printf("error: failed executing template: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) carSelectHandler(w http.ResponseWriter, req *http.Request) {
	_, err := auth.ValidateCookie(req)
	if err != nil {
		log.Printf("error: cookie failed validation on car select: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	selectedTrack := req.FormValue("track")
	selectedCar := req.FormValue("car")
	selectedSeat := req.FormValue("seat")

	intSeat, err := strconv.Atoi(selectedSeat)
	if err != nil {
		log.Printf("error: failed parsing selected seat to int: %s", err.Error())
		return
	}

	carKey, err := s.socketIOServer.CarConns.GetKeyByCarShortName(selectedCar)
	if err != nil {
		log.Printf("error: failed getting car by shortname: %s\n", selectedCar)
		return
	}

	password, err := s.socketIOServer.CarConns.GetPassword(carKey)
	if err != nil {
		log.Printf("error: failed getting car password: %s\n", carKey)
		return
	}

	hasPassword := false
	if password != "" {
		hasPassword = true
	}

	err = s.templates["car_select"].Execute(w, models.CarTMPLData{
		CarShortName:   selectedCar,
		TrackShortName: selectedTrack,
		SeatNumber:     intSeat,
		HasPassword:    hasPassword,
	})
	if err != nil {
		log.Printf("failed executing template: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) driveHandler(w http.ResponseWriter, req *http.Request) {
	_, err := auth.ValidateCookie(req)
	if err != nil {
		log.Printf("error: cookie failed validation on drive: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// selectedTrack := req.FormValue("car_short_name")
	// selectedCar := req.FormValue("track_short_name")
	// selectedSeat := req.FormValue("seat_number")
	// password := req.FormValue("password")

	reqBody := models.DriveReqBody{}
	err = json.NewDecoder(req.Body).Decode(&reqBody)
	if err != nil {
		log.Printf("error decoding json body: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	intSeat, err := strconv.Atoi(reqBody.SeatNumber)
	if err != nil {
		log.Printf("error: failed parsing selected seat to int: %s", err.Error())
		return
	}

	carKey, err := s.socketIOServer.CarConns.GetKeyByCarShortName(reqBody.CarShortName)
	if err != nil {
		log.Printf("error: failed getting car by shortname: %s\n", reqBody.CarShortName)
		return
	}

	carPassword, err := s.socketIOServer.CarConns.GetPassword(carKey)
	if err != nil {
		log.Printf("error: failed getting car password: %s\n", carKey)
		return
	}

	if reqBody.Password != carPassword {
		log.Println("error: user provided password does not match car password")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	videoData := models.VideoTMPLData{
		TrackName:  reqBody.TrackShortName,
		CarName:    reqBody.CarShortName,
		SeatNumber: intSeat,
	}

	startConnecting := fmt.Sprintf(`{"startconnecting": {"track": "%s", "car": "%s", "seat":%d}}`, reqBody.TrackShortName, reqBody.CarShortName, intSeat)
	w.Header().Add("HX-Trigger-After-Settle", startConnecting)
	err = s.templates["video"].Execute(w, videoData)
	if err != nil {
		log.Printf("error: failed executing video template: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//Executes to the response writer handles the return of html to the user
}
