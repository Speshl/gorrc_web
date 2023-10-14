package socketio

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Speshl/gorrc_web/internal/service/server/models"
	"github.com/Speshl/gorrc_web/internal/service/stores/v1gorrc"
	"github.com/google/uuid"
	socketio "github.com/googollee/go-socket.io"
)

type UserConnections struct {
	connections map[string]*UserConnection
	lock        sync.RWMutex
}

type UserConnection struct {
	ID      string
	Socket  socketio.Conn
	Cancel  context.CancelFunc
	CTX     context.Context
	User    *v1gorrc.User
	Car     *v1gorrc.Car
	Track   *v1gorrc.Track
	SeatNum int

	HealthTTL time.Time
}

func NewUserConnections() *UserConnections {
	return &UserConnections{
		connections: make(map[string]*UserConnection, 10),
	}
}

func NewUserConnection(socketConn socketio.Conn, user *v1gorrc.User) *UserConnection {
	ctx, cancelCTX := context.WithCancel(context.Background())
	conn := UserConnection{
		ID:     socketConn.ID(),
		Socket: socketConn,
		Cancel: cancelCTX,
		CTX:    ctx,
		User:   user,
	}

	return &conn
}

func (c *UserConnections) NewUserConnection(socketConn socketio.Conn, user *v1gorrc.User) string {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.connections[socketConn.ID()] = NewUserConnection(socketConn, user)
	return socketConn.ID()
}

func (c *UserConnections) Disconnect(id string) (uuid.UUID, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	conn, ok := c.connections[id]
	if !ok {
		return uuid.Nil, ErrNotFound
	}

	dbId := conn.disconnect()
	delete(c.connections, id)
	return dbId, nil
}

func (c *UserConnection) disconnect() uuid.UUID {
	c.Cancel()
	log.Printf("user %s disconnected\n", c.User.UserName)
	return c.User.Id
}

func (c *UserConnections) GetCar(id string) (*v1gorrc.Car, int, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	conn, ok := c.connections[id]
	if !ok {
		return nil, -1, ErrNotFound
	}
	return conn.Car, conn.SeatNum, nil
}

func (c *UserConnections) GetUser(id string) (*v1gorrc.User, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	conn, ok := c.connections[id]
	if !ok {
		return nil, ErrNotFound
	}

	return conn.User, nil
}

func (c *UserConnections) GetTrack(id string) (*v1gorrc.Track, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	conn, ok := c.connections[id]
	if !ok {
		return nil, ErrNotFound
	}

	return conn.Track, nil
}

func (c *UserConnections) SetCar(id string, car *v1gorrc.Car, seatNum int) error {
	if car == nil {
		return fmt.Errorf("cannot set user to use nil car")
	}

	if seatNum < 0 || seatNum > car.SeatCount {
		return fmt.Errorf("unsupported seat number: %d", seatNum)
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	_, ok := c.connections[id]
	if !ok {
		return ErrNotFound
	}

	c.connections[id].Car = car
	c.connections[id].SeatNum = seatNum
	return nil
}

func (c *UserConnections) SetUser(id string, user *v1gorrc.User) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, ok := c.connections[id]
	if !ok {
		return ErrNotFound
	}

	c.connections[id].User = user
	return nil
}

func (c *UserConnections) SetTrack(id string, track *v1gorrc.Track) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, ok := c.connections[id]
	if !ok {
		return ErrNotFound
	}

	c.connections[id].Track = track
	return nil
}

func (c *UserConnections) Emit(id string, event string, msg ...interface{}) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	_, ok := c.connections[id]
	if ok {
		log.Printf("user emitting %s@%s\n", id, event)
		c.connections[id].Socket.Emit(event, msg)
	}
}

func (c *UserConnections) GetKeyByCarIdAndSeat(carId uuid.UUID, seatNum int) (string, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	for i := range c.connections {
		if c.connections[i].Car.Id == carId && seatNum == c.connections[i].SeatNum {
			return i, nil
		}
	}
	return "", ErrNotFound
}

func (c *UserConnections) GetDataForActiveTracks() map[uuid.UUID]models.TrackData {
	c.lock.Lock()
	defer c.lock.Unlock()

	activeTracks := make(map[uuid.UUID]models.TrackData, len(c.connections))
	for i := range c.connections {
		trackUUID := c.connections[i].Track.Id

		val, ok := activeTracks[trackUUID]
		if !ok {
			val = models.TrackData{}
		}

		val.UserCount++
		activeTracks[trackUUID] = val
	}
	return activeTracks
}
