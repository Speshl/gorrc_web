package socketio

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Speshl/gorrc_web/internal/service/stores/v1gorrc"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
)

type SocketIOServer struct {
	socketio  *socketio.Server
	UserConns *UserConnections
	CarConns  *CarConnections
	store     v1gorrc.StoreAPI
	cfg       SocketIOServerCfg
}

type SocketIOServerCfg struct {
}

var allowOriginFunc = func(r *http.Request) bool {
	return true
}

func NewSocketServer(cfg SocketIOServerCfg, store v1gorrc.StoreAPI) *SocketIOServer {
	socketioServer := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
	})

	return &SocketIOServer{
		socketio:  socketioServer,
		UserConns: NewUserConnections(),
		CarConns:  NewCarConnections(),
		store:     store,
		cfg:       cfg,
	}
}

func (s *SocketIOServer) RegisterSocketIOHandlers() {
	s.socketio.OnConnect("/", s.onConnect)
	s.socketio.OnDisconnect("/", s.onDisconnect)
	s.socketio.OnError("/", s.onError)

	s.socketio.OnEvent("/", "", s.ConnectUser)

	s.socketio.OnEvent("/", "car_connect", s.ConnectCar)
	s.socketio.OnEvent("/", "user_connect", s.ConnectUser)

	s.socketio.OnEvent("/", "car_healthy", s.HealthyCar)
	s.socketio.OnEvent("/", "user_healthy", s.HealthyUser)

	s.socketio.OnEvent("/", "offer", s.OnUserOffer)
	s.socketio.OnEvent("/", "answer", s.OnCarAnswer)

}

func (s *SocketIOServer) StartHealthChecker(ctx context.Context) {
	go func() {
		healthTicker := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-ctx.Done():
				log.Printf("stopping health checker: %s\n", ctx.Err())
				return
			case <-healthTicker.C:
				s.CarConns.RemoveUnhealthy()
				//s.UserConns.removeUnhealthy() TODO: Need health checks for user?
			}
		}
	}()
}

func (s *SocketIOServer) Close() error {
	return s.socketio.Close()
}

func (s *SocketIOServer) Serve() error {
	return s.socketio.Serve()
}

func (s *SocketIOServer) GetHandler() *socketio.Server {
	return s.socketio
}

// func (s *SocketIOServer) NewClientConn(socketConn socketio.Conn) (*Connection, error) {
// 	clientConn, err := NewConnection(socketConn)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// err = clientConn.RegisterHandlers(s.carAudioTrack, s.carVideoTrack, s.memeSoundChannel)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	// Set the handler for Peer connection state
// 	// This will notify you when the peer has connected/disconnected
// 	clientConn.PeerConnection.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
// 		log.Printf("Peer Connection State has changed: %s\n", state.String())
// 		if state == webrtc.PeerConnectionStateFailed {
// 			// Wait until PeerConnection has had no network activity for 30 seconds or another failure. It may be reconnected using an ICE Restart.
// 			// Use webrtc.PeerConnectionStateDisconnected if you are interested in detecting faster timeout.
// 			// Note that the PeerConnection may come back from PeerConnectionStateDisconnected.
// 			log.Println("Peer Connection has gone to failed")
// 			s.RemoveClient(socketConn.ID())
// 		}
// 	})

// 	return clientConn, nil
// }
