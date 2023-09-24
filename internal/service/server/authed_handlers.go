package server

import (
	"log"
	"net/http"

	"github.com/Speshl/gorrc_web/internal/service/auth"
	"github.com/Speshl/gorrc_web/internal/service/server/models"
)

// TODO: Ensure user is logged in before processing these requests
func (s *Server) trackListHandler(w http.ResponseWriter, req *http.Request) {

	_, err := auth.ValidateCookie(req)
	if err != nil {
		log.Printf("cookie failed validation on track list: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	trackListData, err := s.GetTrackList(req.Context())
	if err != nil {
		log.Printf("failed generating track list: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.templates["track_list"].Execute(w, trackListData)
	if err != nil {
		log.Printf("failed executing template: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) trackSelectHandler(w http.ResponseWriter, req *http.Request) {
	_, err := auth.ValidateCookie(req)
	if err != nil {
		log.Printf("cookie failed validation on track select: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	selectedTrack := req.FormValue("track")

	err = s.templates["track_select"].Execute(w, models.TrackTMPLData{
		ShortName: selectedTrack,
	})
	if err != nil {
		log.Printf("failed executing template: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) carListHandler(w http.ResponseWriter, req *http.Request) {
	_, err := auth.ValidateCookie(req)
	if err != nil {
		log.Printf("cookie failed validation on car list: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	track, err := s.store.GetTrackByShortName(req.Context(), req.FormValue("track"))
	if err != nil {
		log.Printf("failed getting track by short name: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	carData, err := s.GetCarListForTrack(req.Context(), track.Id)
	if err != nil {
		log.Printf("failed generating car list: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.templates["car_list"].Execute(w, carData)
	if err != nil {
		log.Printf("failed executing template: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) carSelectHandler(w http.ResponseWriter, req *http.Request) {
	_, err := auth.ValidateCookie(req)
	if err != nil {
		log.Printf("cookie failed validation on car select: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	selectedTrack := req.FormValue("track")
	selectedCar := req.FormValue("car")

	err = s.templates["car_select"].Execute(w, models.CarTMPLData{
		ShortName:      selectedCar,
		TrackShortName: selectedTrack,
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
		log.Printf("cookie failed validation on drive: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	selectedTrack := req.FormValue("track")
	selectedCar := req.FormValue("car")
	//user, err := s.store.GetUserByUserName(req.Context(), userName)

	videoData := models.VideoTMPLData{
		TrackName: selectedTrack,
		CarName:   selectedCar,
	}

	w.Header().Add("HX-Trigger-After-Settle", "startconnecting")
	err = s.templates["video"].Execute(w, videoData)
	if err != nil {
		log.Printf("failed executing video template: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//Executes to the response writer handles the return of html to the user
}
