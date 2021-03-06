package lobby

import (
	"context"
	"log"
	"sync"
	"time"

	"matchmaker/playerdb"
)

// PlayerId type.
type PlayerId string

// Player contains the data associated to each player
// connected to the service.
type Player struct {
	// Player ID (must be unique).
	Id PlayerId
	// queue used to signal events to the player's connection.
	// responseQueue chan<- Game
	// channel that will receive players to ask if it's a valid candidate to generate a game
	playersQueue chan *Player

	MatchQueue chan *Match

	elo int

	isWaiting bool

	relaxRequirements float64

	mtx sync.Mutex
}

func NewPlayer(id PlayerId) Player {
	return Player{
		isWaiting:         true,
		elo:               playerdb.PlayerDB.GetData(string(id)),
		playersQueue:      make(chan *Player),
		MatchQueue:        make(chan *Match),
		Id:                id,
		relaxRequirements: 1.0,
	}
}

// Match represents a pair of players that has been
// selected to play against each other (i.e. has been matched).
type Match struct {
	GameRoom string
	player1  *Player
	player2  *Player
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
	for {
		select {
		case <-ctx.Done():
			// stop the goroutine
			return ctx.Err()

		case player := <-mm.join:
			log.Printf("Inserted player in queue %+v %+v\n", mm, player)

			go func(player *Player) {
				// notify all players about player
				for _, p := range mm.players {
					if p.Id != player.Id && p.isWaiting {
						log.Println("Player joining...")
						player.playersQueue <- p
					}
				}
			}(player)

		case <-time.After(time.Second * 30):
			log.Printf("Match maker status %+v", mm)
		}
	}
}

func (mm *MatchMaker) AddPlayer(player *Player) {
	mm.players[player.Id] = player
}

// Add a client to the matchmaking waiting queue.
func (mm *MatchMaker) Add(ctx context.Context, id PlayerId) error {
	log.Printf("Add of %+v %v\n", mm, id)
	player, _ := mm.players[id]

	log.Printf("Player %+v\n", player)

	select {
	case mm.join <- player:
	case <-ctx.Done():
		return ctx.Err()
	}

	// successfully inserted in queue
	return nil
}

// createMatch creates a new match and informs the players
// about the event.
