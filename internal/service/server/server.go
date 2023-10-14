package server

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Speshl/gorrc_web/internal/service/config"
	"github.com/Speshl/gorrc_web/internal/service/server/socketio"
	"github.com/Speshl/gorrc_web/internal/service/stores/v1gorrc"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	ctx       context.Context
	ctxCancel context.CancelFunc

	cfg            config.ServerCfg
	socketIOServer *socketio.SocketIOServer
	store          v1gorrc.StoreAPI
	templates      map[string]*template.Template
}

func NewServer(ctx context.Context, cfg config.ServerCfg, store v1gorrc.StoreAPI) *Server {
	ctx, cancel := context.WithCancel(ctx)
	return &Server{
		ctx:       ctx,
		ctxCancel: cancel,
		cfg:       cfg,
		store:     store,
	}
}

func (s *Server) StartServing() error {
	s.socketIOServer = socketio.NewSocketServer(s.cfg.SocketIOCfg, s.store)
	s.socketIOServer.RegisterSocketIOHandlers()

	s.RegisterHTTPHandlers()
	s.ParseHTMLTemplates()

	group, groupCtx := errgroup.WithContext(s.ctx)

	group.Go(func() error {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

		select {
		case sig := <-signalChannel:
			fmt.Printf("received signal: %s\n", sig)
			s.ctxCancel()
		case <-groupCtx.Done():
			fmt.Printf("closing signal goroutine\n")
			return groupCtx.Err()
		}

		return nil
	})

	group.Go(func() error {
		log.Println("start serving socketio...")
		defer log.Println("stop serving socketio...")

		err := s.socketIOServer.Serve()
		if err != nil {
			log.Printf("socketio listen error: %s\n", err)
		}
		s.ctxCancel() //stop anything else on this context because the socker server stopped

		return err
	})

	group.Go(func() error {
		log.Println("start serving http...")
		defer log.Println("stop serving http...")
		addr := fmt.Sprintf(":%s", s.cfg.Port)
		err := http.ListenAndServe(addr, nil)
		if !errors.Is(err, http.ErrServerClosed) {
			log.Printf("http server error: %v", err)
		}
		s.ctxCancel() //stop anything else on this context because the http server stopped
		return err
	})

	err := group.Wait()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Println("context was cancelled")
			return nil
		} else {
			return fmt.Errorf("server stopping due to error - %w", err)
		}
	}
	return nil
}

func (s *Server) RegisterHTTPHandlers() {
	http.HandleFunc("/home", s.homeHandler)
	http.HandleFunc("/loginorregister", s.loginOrRegisterHandler)
	http.HandleFunc("/login", s.loginHandler)
	http.HandleFunc("/register", s.registerHandler)
	http.HandleFunc("/show_register", s.showRegisterHandler)

	http.HandleFunc("/track_list", s.trackListHandler)
	http.HandleFunc("/track_select", s.trackSelectHandler)
	http.HandleFunc("/car_list", s.carListHandler)
	http.HandleFunc("/car_select", s.carSelectHandler)
	http.HandleFunc("/drive", s.driveHandler)

	//serves js and static html
	http.Handle("/", http.FileServer(http.Dir("public/")))
	//sets up socket connections for video/commands
	http.Handle("/socket.io/", s.socketIOServer.GetHandler())

}

func (s *Server) ParseHTMLTemplates() {
	s.templates = make(map[string]*template.Template)
	s.templates["car_list"] = template.Must(template.ParseFiles("public/templates/car_list.tmpl"))
	s.templates["car_select"] = template.Must(template.ParseFiles("public/templates/car_select.tmpl"))
	s.templates["main"] = template.Must(template.ParseFiles("public/templates/main.tmpl"))
	s.templates["login_or_register"] = template.Must(template.ParseFiles("public/templates/login_or_register.tmpl"))
	s.templates["track_list"] = template.Must(template.ParseFiles("public/templates/track_list.tmpl"))
	s.templates["track_select"] = template.Must(template.ParseFiles("public/templates/track_select.tmpl"))
	s.templates["video"] = template.Must(template.ParseFiles("public/templates/video.tmpl"))
	s.templates["register"] = template.Must(template.ParseFiles("public/templates/register.tmpl"))
	s.templates["register_success"] = template.Must(template.ParseFiles("public/templates/register_success.tmpl"))
}
