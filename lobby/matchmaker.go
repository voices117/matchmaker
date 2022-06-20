package lobby

import (
	"context"
	"fmt"
	"log"
	"time"
)

var (
	LoeThreshold int = 10
)

// PlayerId type.
type PlayerId string

// Player contains the data associated to each player
// connected to the service.
type Player struct {
	// Player ID (must be unique).
	Id PlayerId
	// Player LOE
	Loe int
	// queue used to signal events to the player's connection.
	responseQueue chan<- Game
}

// Match represents a pair of players that has been
// selected to play against each other (i.e. has been matched).
type Match struct {
	player1 *Player
	player2 *Player
}

// Represents the created game for a particular match.
type Game struct {
	// Id is the game Id that corresponds the match.
	Id string
}

// MatchMaker keeps track of the players enqueued waiting for
// a match and tries to find the most suitable matches between
// them.
type MatchMaker struct {
	// players in the waiting queue.
	players map[PlayerId]*Player

	// channel to receive new player connection messages.
	join chan *Player

	// channel were matched pairs are sent to
	output chan Match
}

// NewMatchMaker creates, initializes and returns a MatchMaker
// instance.
func NewMatchMaker() MatchMaker {
	return MatchMaker{
		players: make(map[PlayerId]*Player),
		join:    make(chan *Player),
	}
}

// Start the match making algorithm. This function should be
// called in a separate goroutine.
func (mm *MatchMaker) Start(ctx context.Context) error {
	// matchmaking logic: the firs player in que found that is within LoeThreshold is selected
	
	for {
		select {
		case <-ctx.Done():
			// stop the goroutine
			return ctx.Err()

		case player := <-mm.join:

			var match bool = false
			for _, p := range mm.players {
				var loe_dif = player.Loe - p.Loe
				if loe_dif < LoeThreshold && loe_dif > -LoeThreshold {
					fmt.Printf("Matched %v (%v) against %v (%v)\n [QUEUE: %v]", p.Id, p.Loe, player.Id, player.Loe, len(mm.players)-1)
					go mm.createMatch(Match{
						player1: p,
						player2: player,
					})
					delete(mm.players, p.Id)
					match = true
				}
			}
			if !match {
				mm.players[player.Id] = player
			}
		}
	}
}

// Add a client to the matchmaking waiting queue.
func (mm *MatchMaker) Add(ctx context.Context, id PlayerId, loe int, response chan<- Game) error {
	select {
	case mm.join <- &Player{Id: id, Loe: loe, responseQueue: response}:
	case <-ctx.Done():
		return ctx.Err()
	}

	// successfully inserted in queue
	return nil
}

// createMatch creates a new match and informs the players
// about the event.
func (mm *MatchMaker) createMatch(match Match) {
	// TODO: create game and get real ID
	game := Game{
		Id: "<fake game ID>",
	}

	send := func(player *Player) {
		select {
		case player.responseQueue <- game:
		case <-time.After(time.Second * 5):
			log.Printf("Failed sending game Id to player '%v'", player.Id)
		}
	}

	go send(match.player1)
	go send(match.player2)
}
