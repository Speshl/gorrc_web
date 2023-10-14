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

var ErrNotFound = fmt.Errorf("not found")

const AddHealthPerCheckIn = 2 * time.Minute

type CarConnections struct {
	connections map[string]*CarConnection
	lock        sync.RWMutex
}

type CarConnection struct {
	ID       string
	Socket   socketio.Conn
	Cancel   context.CancelFunc
	CTX      context.Context
	Seats    []*v1gorrc.User
	Car      *v1gorrc.Car
	Track    *v1gorrc.Track
	Password string

	HealthTTL time.Time
}

func NewCarConnections() *CarConnections {
	return &CarConnections{
		connections: make(map[string]*CarConnection, 10),
	}
}

func NewCarConnection(socketConn socketio.Conn, car *v1gorrc.Car, track *v1gorrc.Track, password string) *CarConnection {
	log.Printf("creating car connection %s\n", socketConn.ID())

	ctx, cancelCTX := context.WithCancel(context.Background())
	conn := CarConnection{
		ID:        socketConn.ID(),
		Socket:    socketConn,
		Cancel:    cancelCTX,
		CTX:       ctx,
		Seats:     make([]*v1gorrc.User, car.SeatCount),
		Car:       car,
		Track:     track,
		HealthTTL: time.Now().Add(AddHealthPerCheckIn),
		Password:  password,
	}

	return &conn
}

func (c *CarConnections) NewCarConnection(socketConn socketio.Conn, car *v1gorrc.Car, track *v1gorrc.Track, password string) string {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.connections[socketConn.ID()] = NewCarConnection(socketConn, car, track, password)
	return socketConn.ID()
}

func (c *CarConnections) Disconnect(id string) (uuid.UUID, error) {
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

func (c *CarConnection) disconnect() uuid.UUID {
	c.Cancel()
	log.Printf("car %s disconnected\n", c.Car.Name)
	return c.Car.Id
}

func (c *CarConnections) RemoveUnhealthy() {
	c.lock.Lock()
	defer c.lock.Unlock()

	currentTime := time.Now()
	for i := range c.connections {
		if currentTime.After(c.connections[i].HealthTTL) {
			log.Printf("car %s became unhealthy\n", c.connections[i].Car.Name)
			c.connections[i].disconnect()
			delete(c.connections, i)
		}
	}
}

func (c *CarConnections) ResetHealth(id string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	conn, ok := c.connections[id]
	if !ok {
		return
	}
	conn.HealthTTL = time.Now().Add(AddHealthPerCheckIn)
}

// Get car by key
func (c *CarConnections) GetPassword(id string) (string, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	conn, ok := c.connections[id]
	if !ok {
		return "", ErrNotFound
	}

	return conn.Password, nil
}

func (c *CarConnections) GetCar(id string) (*v1gorrc.Car, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	conn, ok := c.connections[id]
	if !ok {
		return nil, ErrNotFound
	}

	return conn.Car, nil
}

func (c *CarConnections) GetUserInSeat(id string, seatNum int) (*v1gorrc.User, error) {
	if seatNum < 0 {
		return nil, fmt.Errorf("unsupported seat number: %d", seatNum)
	}

	c.lock.RLock()
	defer c.lock.RUnlock()
	conn, ok := c.connections[id]
	if !ok {
		return nil, ErrNotFound
	}

	if seatNum > len(conn.Seats) {
		return nil, fmt.Errorf("unsupported seat number: %d", seatNum)
	}

	return conn.Seats[seatNum], nil
}

func (c *CarConnections) GetTrack(id string) (*v1gorrc.Track, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	conn, ok := c.connections[id]
	if !ok {
		return nil, ErrNotFound
	}

	return conn.Track, nil
}

func (c *CarConnections) SetCar(id string, car *v1gorrc.Car) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, ok := c.connections[id]
	if !ok {
		return ErrNotFound
	}

	c.connections[id].Car = car
	return nil
}

func (c *CarConnections) SetUserInSeat(id string, seatNum int, user *v1gorrc.User) error {
	if seatNum < 0 {
		return fmt.Errorf("invalid seat number: %d", seatNum)
	}

	c.lock.Lock()
	defer c.lock.Unlock()
	_, ok := c.connections[id]
	if !ok {
		return ErrNotFound
	}

	if seatNum > len(c.connections[id].Seats) {
		return fmt.Errorf("unsupported seat number: %d", seatNum)
	}

	c.connections[id].Seats[seatNum] = user
	return nil
}

func (c *CarConnections) SetTrack(id string, track *v1gorrc.Track) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, ok := c.connections[id]
	if !ok {
		return ErrNotFound
	}

	c.connections[id].Track = track
	return nil
}

func (c *CarConnections) ClearUser(userId uuid.UUID) {
	c.lock.Lock()
	defer c.lock.Unlock()
	for i := range c.connections {
		for j := range c.connections[i].Seats {
			if c.connections[i].Seats[j] != nil && c.connections[i].Seats[j].Id == userId {
				c.connections[i].Seats[j] = nil
				log.Printf("userid %s removed from car %s@%s seat %d\n", userId, c.connections[i].Car.Name, c.connections[i].Track.Name, j)
			}
		}
	}
}

func (c *CarConnections) Emit(id string, event string, msg ...interface{}) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	_, ok := c.connections[id]
	if ok {
		log.Printf("car emitting %s@%s\n", id, event)
		c.connections[id].Socket.Emit(event, msg)
	}
}

func (c *CarConnections) GetKeyByCarId(carId uuid.UUID) (string, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	for i := range c.connections {
		if c.connections[i].Car.Id == carId {
			return i, nil
		}
	}
	return "", ErrNotFound
}

func (c *CarConnections) GetKeyByCarShortName(name string) (string, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	for i := range c.connections {
		if c.connections[i].Car == nil {
			continue
		}
		if c.connections[i].Car.ShortName == name {
			return i, nil
		}
	}
	return "", ErrNotFound
}

func (c *CarConnections) GetDataForActiveTracks() map[uuid.UUID]models.TrackData {
	c.lock.Lock()
	defer c.lock.Unlock()

	activeTracks := make(map[uuid.UUID]models.TrackData, len(c.connections))
	for i := range c.connections {
		trackUUID := c.connections[i].Track.Id

		val, ok := activeTracks[trackUUID]
		if !ok {
			val = models.TrackData{}
		}

		activeSeats := 0
		for j := range c.connections[i].Seats {
			if c.connections[i].Seats[j] != nil {
				activeSeats++
			}
		}

		val.CarCount++
		val.SeatCount += len(c.connections[i].Seats)
		val.UserCount += activeSeats
		activeTracks[trackUUID] = val
	}
	return activeTracks
}

func (c *CarConnections) GetActiveCarIds() []uuid.UUID {
	c.lock.RLock()
	defer c.lock.RUnlock()

	activeCarIds := make([]uuid.UUID, len(c.connections))
	for i := range c.connections {
		activeCarIds = append(activeCarIds, c.connections[i].Car.Id)
	}
	return activeCarIds
}

func (c *CarConnections) GetSeatStatusByCarId(carId uuid.UUID) ([]string, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	seatStatus := make([]string, 0)
	for i := range c.connections {
		if c.connections[i].Car != nil && c.connections[i].Car.Id == carId {

			for j := range c.connections[i].Seats {
				if c.connections[i].Seats[j] == nil {
					seatStatus = append(seatStatus, "available")
				} else {
					seatStatus = append(seatStatus, "occupied")
				}
			}
			return seatStatus, nil
		}
	}
	return seatStatus, ErrNotFound
}
