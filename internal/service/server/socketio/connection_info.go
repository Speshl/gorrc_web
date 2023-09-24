package socketio

import "github.com/google/uuid"

func (c *Connections) GetKeyByCarId(carId uuid.UUID) string {
	c.lock.RLock()
	defer c.lock.RUnlock()

	for i := range c.connections {
		if c.connections[i].Car.Id == carId {
			return i
		}
	}
	return ""
}

func (c *Connections) GetKeyByCarShortName(name string) string {
	c.lock.RLock()
	defer c.lock.RUnlock()

	for i := range c.connections {
		if c.connections[i].Car == nil {
			continue
		}
		if c.connections[i].Car.ShortName == name {
			return i
		}
	}
	return ""
}

//Returns [uuid]count
func (c *Connections) GetActiveTracks() map[uuid.UUID]int {
	c.lock.RLock()
	defer c.lock.RUnlock()

	activeTracks := make(map[uuid.UUID]int, len(c.connections))
	for i := range c.connections {
		trackUUID := c.connections[i].Track.Id
		activeTracks[trackUUID] += 1
	}
	return activeTracks
}

func (c *Connections) GetActiveCarIds() []uuid.UUID {
	c.lock.RLock()
	defer c.lock.RUnlock()

	activeCarIds := make([]uuid.UUID, len(c.connections))
	for i := range c.connections {
		activeCarIds = append(activeCarIds, c.connections[i].Car.Id)
	}
	return activeCarIds
}

func (c *Connections) GetOccupiedCarIds() []uuid.UUID {
	c.lock.RLock()
	defer c.lock.RUnlock()

	occupiedCarIds := make([]uuid.UUID, len(c.connections))
	for i := range c.connections {
		if c.connections[i].User != nil {
			occupiedCarIds = append(occupiedCarIds, c.connections[i].Car.Id)
		}
	}
	return occupiedCarIds
}

// func (s *SocketIOServer) getUserConnByName(name string) *UserConnection {
// 	s.userLock.RLock()
// 	defer s.userLock.RUnlock()

// 	for _, con := range s.userConnections {
// 		if con.Details.Name == name {
// 			return con
// 		}
// 	}
// 	return nil
// }

// func (s *SocketIOServer) GetTracks() []string {
// 	s.carLock.RLock()
// 	defer s.carLock.RUnlock()

// 	trackMap := make(map[string]bool, len(s.carConnections))
// 	for _, carConn := range s.carConnections {
// 		trackMap[carConn.Details.ShortTrackName] = true
// 	}

// 	trackList := make([]string, 0, len(trackMap))
// 	for track := range trackMap {
// 		trackList = append(trackList, track)
// 	}
// 	return trackList
// }

// func (s *SocketIOServer) GetCars() []string {
// 	s.carLock.RLock()
// 	defer s.carLock.RUnlock()

// 	trackMap := make(map[string]bool, len(s.carConnections))
// 	for _, carConn := range s.carConnections {
// 		trackMap[carConn.Id] = true
// 	}

// 	trackList := make([]string, 0, len(trackMap))
// 	for track := range trackMap {
// 		trackList = append(trackList, track)
// 	}
// 	return trackList
// }

// func (s *SocketIOServer) GetCarsForTrack(track string) []string {
// 	s.carLock.RLock()
// 	defer s.carLock.RUnlock()

// 	carList := make([]string, 0)
// 	for _, carConn := range s.carConnections {
// 		if carConn.Details.ShortTrackName == track {
// 			carList = append(carList, carConn.Details.ShortCarName)
// 		}
// 	}
// 	return carList
// }
