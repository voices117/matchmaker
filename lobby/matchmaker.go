package lobby

import (
	"context"
	"fmt"
	"time"
)

// PlayerId type.
type PlayerId string

// Player contains the data associated to each player
// connected to the service.
type Player struct {
	// Player ID (must be unique).
	Id PlayerId
	// Name of the player.
	Name string
	// ELO score.
	ELO int
	// Time when the player joined the waiting queue.
	Time time.Time
}

// MatchMaker keeps track of the players enqueued waiting for
// a match and tries to find the most suitable matches between
// them.
type MatchMaker struct {
	// players in the waiting queue.
	players map[PlayerId]*Player

	// channel to receive new player connection messages.
	Join chan PlayerId
}

// NewMatchMaker creates, initializes and returns a MatchMaker
// instance.
func NewMatchMaker() MatchMaker {
	return MatchMaker{
		players: make(map[PlayerId]*Player),
		Join:    make(chan PlayerId),
	}
}

// Start the match making algorithm. This function should be
// called in a separate goroutine.
func (mm *MatchMaker) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			// stop the goroutine
			return ctx.Err()

		case playerId := <-mm.Join:
			// TODO: implement actual matchmaking logic
			fmt.Printf("Player %v joined\n", playerId)
		}
	}
}

// Add a client to the matchmaking waiting queue.
func (mm *MatchMaker) Add(ctx context.Context, id PlayerId) error {
	select {
	case mm.Join <- PlayerId(id):
	case <-ctx.Done():
		return ctx.Err()
	}

	// successfully inserted in queue
	return nil
}
