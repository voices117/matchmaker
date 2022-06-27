package room

import (
	"context"
	"log"
	"sync"
)

// RoomManager stores all the existing rooms and allows to
// create or destroy then.
type RoomManager struct {
	rooms map[string]*Room
	mtx   sync.Mutex

	gamesCtx    context.Context
	gamesCancel context.CancelFunc
}

// NewRoomManager creates and initializes a room manager.
func NewRoomManager() RoomManager {
	ctx, cancel := context.WithCancel(context.Background())
	return RoomManager{
		rooms:       make(map[string]*Room),
		gamesCtx:    ctx,
		gamesCancel: cancel,
	}
}

// Stop cancels all ongoing games.
func (rm *RoomManager) Stop() {
	rm.gamesCancel()
}

// GetOrCreate will return a game room associated to the given Id.
// If it does not exist, it is created.
func (rm *RoomManager) GetOrCreate(id string) *Room {
	rm.mtx.Lock()
	defer rm.mtx.Unlock()

	room, exists := rm.rooms[id]
	if !exists {
		room = NewRoom(id, rm.RemoveGame)
		rm.rooms[id] = room

		go room.RunGame(rm.gamesCtx)
	}
	return room
}

// RemoveGame removes a game from the list of existing rooms.
func (rm *RoomManager) RemoveGame(id string) {
	log.Printf("Removing game '%v'", id)

	rm.mtx.Lock()
	defer rm.mtx.Unlock()

	delete(rm.rooms, id)
}
