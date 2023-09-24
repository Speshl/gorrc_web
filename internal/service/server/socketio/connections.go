package socketio

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Speshl/gorrc_web/internal/service/stores/v1gorrc"
	"github.com/google/uuid"
	socketio "github.com/googollee/go-socket.io"
)

const AddHealthPerCheckIn = 2 * time.Minute

type Connections struct {
	connections map[string]*Connection
	lock        sync.RWMutex
}

type Connection struct {
	ID     string
	Type   string //TODO convert to enum (Car/User)
	Socket socketio.Conn
	Cancel context.CancelFunc
	CTX    context.Context
	User   *v1gorrc.User
	Car    *v1gorrc.Car
	Track  *v1gorrc.Track

	HealthTTL time.Time
}

func NewConnections() *Connections {
	return &Connections{
		connections: make(map[string]*Connection, 10),
	}
}

func NewCarConnection(socketConn socketio.Conn, car *v1gorrc.Car, track *v1gorrc.Track) *Connection {
	log.Printf("creating car connection %s\n", socketConn.ID())

	ctx, cancelCTX := context.WithCancel(context.Background())
	conn := Connection{
		ID:        socketConn.ID(),
		Socket:    socketConn,
		Cancel:    cancelCTX,
		CTX:       ctx,
		Car:       car,
		Track:     track,
		Type:      "car",
		HealthTTL: time.Now().Add(AddHealthPerCheckIn),
	}

	return &conn
}

func NewUserConnection(socketConn socketio.Conn, user *v1gorrc.User) *Connection {
	ctx, cancelCTX := context.WithCancel(context.Background())
	conn := Connection{
		ID:     socketConn.ID(),
		Socket: socketConn,
		Cancel: cancelCTX,
		CTX:    ctx,
		User:   user,
		Type:   "user",
	}

	return &conn
}

func (c *Connections) NewCarConnection(socketConn socketio.Conn, car *v1gorrc.Car, track *v1gorrc.Track) string {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.connections[socketConn.ID()] = NewCarConnection(socketConn, car, track)
	return socketConn.ID()
}

func (c *Connections) NewUserConnection(socketConn socketio.Conn, user *v1gorrc.User) string {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.connections[socketConn.ID()] = NewUserConnection(socketConn, user)
	return socketConn.ID()
}

func (c *Connections) Disconnect(id string) uuid.UUID {
	c.lock.Lock()
	defer c.lock.Unlock()
	conn, ok := c.connections[id]
	if ok {
		dbId := conn.disconnect()
		delete(c.connections, id)
		return dbId
	}
	return uuid.Nil
}

func (c *Connection) disconnect() uuid.UUID {
	c.Cancel()
	switch c.Type {
	case "car":
		log.Printf("car %s disconnected\n", c.Car.Id)
		return c.Car.Id
	case "user":
		log.Printf("user %s disconnected\n", c.User.UserName)
		return c.User.Id
	default:
		log.Printf("connection has unsupported type")
		return uuid.Nil
	}
}

func (c *Connections) removeUnhealthy() {
	c.lock.Lock()
	defer c.lock.Unlock()

	currentTime := time.Now()
	for i := range c.connections {
		if currentTime.After(c.connections[i].HealthTTL) {
			log.Printf("connection %s became unhealthy\n", i)
			c.connections[i].disconnect()
			delete(c.connections, i)
		}
	}
}

func (c *Connections) resetHealth(id string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	conn, ok := c.connections[id]
	if !ok {
		return
	}
	conn.HealthTTL = time.Now().Add(AddHealthPerCheckIn)
}

func (c *Connections) getCar(id string) *v1gorrc.Car {
	c.lock.RLock()
	defer c.lock.RUnlock()
	conn, ok := c.connections[id]
	if !ok {
		return nil
	}

	return conn.Car
}

func (c *Connections) getUser(id string) *v1gorrc.User {
	c.lock.RLock()
	defer c.lock.RUnlock()
	conn, ok := c.connections[id]
	if !ok {
		return nil
	}

	return conn.User
}

func (c *Connections) getTrack(id string) *v1gorrc.Track {
	c.lock.RLock()
	defer c.lock.RUnlock()
	conn, ok := c.connections[id]
	if !ok {
		return nil
	}

	return conn.Track
}

func (c *Connections) setCar(id string, car *v1gorrc.Car) {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, ok := c.connections[id]
	if ok {
		c.connections[id].Car = car
	}
}

func (c *Connections) setUser(id string, user *v1gorrc.User) {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, ok := c.connections[id]
	if ok {
		c.connections[id].User = user
	}
}

func (c *Connections) setTrack(id string, track *v1gorrc.Track) {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, ok := c.connections[id]
	if ok {
		c.connections[id].Track = track
	}
}

func (c *Connections) clearUser(userId uuid.UUID) {
	c.lock.Lock()
	defer c.lock.Unlock()
	for i := range c.connections {

		if c.connections[i].User != nil && c.connections[i].User.Id == userId {
			c.connections[i].User = nil
			log.Printf("userid %s removed from car connection %s\n", userId, i)
		}
	}
}

func (c *Connections) emit(id string, event string, msg ...interface{}) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	_, ok := c.connections[id]
	if ok {
		log.Printf("Emitting %s@%s\n", id, event)
		c.connections[id].Socket.Emit(event, msg)
	}
}
