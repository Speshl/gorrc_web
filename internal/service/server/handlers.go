package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/Speshl/gorrc_web/internal/service/auth"
	"github.com/Speshl/gorrc_web/internal/service/server/models"
	"github.com/Speshl/gorrc_web/internal/service/stores/v1gorrc"
	"github.com/google/uuid"
)

func (s *Server) homeHandler(w http.ResponseWriter, req *http.Request) {
	var tempBuffer bytes.Buffer

	cookie, err := req.Cookie(auth.TokenName)
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		log.Printf("failed parsing cookie for token: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if cookie.Valid() == nil {
		_, err := auth.ValidateJWT(cookie.Value)
		if err == nil {
			trackListData, err := s.GetTrackList(req.Context())
			if err != nil {
				log.Printf("failed generating track list: %s\n", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			err = s.templates["track_list"].Execute(&tempBuffer, trackListData)
			if err != nil {
				log.Printf("failed executing template: %s\n", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			mainData := models.MainTMPLData{
				Body: template.HTML(tempBuffer.String()),
			}

			err = s.templates["main"].Execute(w, mainData)
			if err != nil {
				log.Printf("failed executing template: %s\n", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}

	}

	err = s.templates["login_or_register"].Execute(&tempBuffer, nil)
	if err != nil {
		log.Printf("failed executing template: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	mainData := models.MainTMPLData{
		Body: template.HTML(tempBuffer.String()),
	}

	err = s.templates["main"].Execute(w, mainData)
	if err != nil {
		log.Printf("failed executing template: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) loginOrRegisterHandler(w http.ResponseWriter, req *http.Request) {
	err := s.templates["login_or_register"].Execute(w, nil)
	if err != nil {
		log.Printf("failed executing template: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) loginHandler(w http.ResponseWriter, req *http.Request) {
	var creds models.Credentials
	log.Printf("Login Handler Recieved Headers: %+v", req.Header)
	err := json.NewDecoder(req.Body).Decode(&creds)
	if err != nil {
		log.Printf("error decoding json body: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.validateCredentials(req.Context(), creds)
	if err != nil {
		log.Printf("error validating credentials: %s", err.Error())
		err = s.templates["login_or_register"].Execute(w, models.LoginTMPLData{
			LoginError: "Invalid username or password",
		})
		if err != nil {
			log.Printf("failed executing template: %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	tokenString, err := auth.GenerateJWT(creds.Username)
	if err != nil {
		log.Printf("error generating JWT: %s", err.Error())
		return
	}

	log.Println("setting user token")
	w.Header().Set("Token", tokenString)
	http.SetCookie(w, &http.Cookie{
		Name:    auth.TokenName,
		Value:   tokenString,
		Expires: auth.AuthTime,
	})

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

func (s *Server) showRegisterHandler(w http.ResponseWriter, req *http.Request) {
	err := s.templates["register"].Execute(w, models.RegisterTMPLData{})
	if err != nil {
		log.Printf("failed executing template: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) registerHandler(w http.ResponseWriter, req *http.Request) {
	body := models.RegisterBody{}

	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		log.Printf("error decoding json body: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if body.Password != body.VerifyPassword {
		err := s.templates["register"].Execute(w, models.RegisterTMPLData{
			PasswordError: "Passwords don't match",
			UserName:      body.Username,
			RealName:      body.RealName,
			Email:         body.Email,
		})
		if err != nil {
			log.Printf("failed executing template: %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	_, err = s.store.GetUserByUserName(req.Context(), body.Username)
	if err != nil && err != v1gorrc.ErrNotFound {
		log.Printf("db error: %s\n", err.Error())
		err := s.templates["register"].Execute(w, models.RegisterTMPLData{
			DisplayNameError: "Username already exists",
			UserName:         body.Username,
			RealName:         body.RealName,
			Email:            body.Email,
		})
		if err != nil {
			log.Printf("failed executing template: %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	err = s.store.CreateUser(req.Context(), v1gorrc.User{
		Id:       uuid.New(),
		UserName: body.Username,
		Password: body.Password, //TODO encrypt password
		RealName: body.RealName,
		Email:    body.Email,
	})
	if err != nil {
		log.Printf("db error: %s\n", err.Error())
		err := s.templates["register"].Execute(w, models.RegisterTMPLData{
			DisplayNameError: "Error creating user",
			UserName:         body.Username,
			RealName:         body.RealName,
			Email:            body.Email,
		})
		if err != nil {
			log.Printf("failed executing template: %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	err = s.templates["register_success"].Execute(w, nil)
	if err != nil {
		log.Printf("failed executing template: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) validateCredentials(ctx context.Context, creds models.Credentials) error {
	user, err := s.store.GetUserByUserName(ctx, creds.Username)
	if err != nil {
		return fmt.Errorf("error getting user: %w", err)
	}

	if user.Password != creds.Password {
		return fmt.Errorf("invalid username and password")
	}
	return nil
}
