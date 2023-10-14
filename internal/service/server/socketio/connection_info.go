package socketio

// //Returns [uuid]count
// func (c *Connections) GetActiveTracks() map[uuid.UUID]int {
// 	c.lock.RLock()
// 	defer c.lock.RUnlock()

// 	activeTracks := make(map[uuid.UUID]int, len(c.connections))
// 	for i := range c.connections {
// 		trackUUID := c.connections[i].Track.Id
// 		activeTracks[trackUUID] += 1
// 	}
// 	return activeTracks
// }

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
