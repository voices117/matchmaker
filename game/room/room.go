package room

import (
	"context"
	"fmt"
	"log"
	"matchmaker/game/msg"
	"matchmaker/game/tictactoe"
	"sync"
	"time"

	// TODO: Uncoment
	// "matchmaker/playerdb"
)

// Room represents a game instance where two players are
// connected to and send/receive events to synchronize state.
type Room struct {
	id string

	joined int

	PlayerX Player
	PlayerO Player

	Game tictactoe.TicTacToe

	mtx sync.Mutex

	onFinishCallback func(string)
}

// NewRoom creates and initializes a new Room instance.
func NewRoom(id string, onFinishCallback func(string)) *Room {
	return &Room{
		id:               id,
		joined:           0,
		PlayerX:          NewPlayer(tictactoe.PLAY_X),
		PlayerO:          NewPlayer(tictactoe.PLAY_O),
		Game:             tictactoe.NewTicTacToe(),
		onFinishCallback: onFinishCallback,
	}
}

// Join adds a player to the game room and returns the corresponding
// player connection so the client can interact with the game. If the
// game room is full, returns an error.
func (room *Room) Join() (*Player, error) {
	room.mtx.Lock()
	defer room.mtx.Unlock()

	if room.joined == 0 {
		room.joined += 1
		return &room.PlayerX, nil
	} else if room.joined == 1 {
		room.joined += 1
		return &room.PlayerO, nil
	} else {
		return nil, fmt.Errorf("Room is full")
	}
}

// RunGame starts a loop that lasts until the game ends or
// the context is cancelled. It will poll for events from
// the players and send game updates.
func (room *Room) RunGame(ctx context.Context) {
	// close gameEvents channels to indicate we are done
	defer close(room.PlayerX.GameEvents)
	defer close(room.PlayerO.GameEvents)

	// trigger the finish callback
	defer room.onFinishCallback(room.id)

	// send the initial game state to the players. This will also
	// work as a barrier for the players to be ready
	// send update status to the players
	if room.sendGameState(&room.PlayerX) != nil || room.sendGameState(&room.PlayerO) != nil {
		log.Printf("Failed sending initial game status")
		return
	}

	for room.Game.GetState() == tictactoe.STATE_UNFINISHED {
		select {
		case <-ctx.Done():
			room.Game.Cancel()
		case play, ok := <-room.PlayerX.PlayerEvents:
			if !ok {
				room.Game.Cancel()
			} else {
				room.setPlay(room.PlayerX, play.Position)
			}
		case play, ok := <-room.PlayerO.PlayerEvents:
			if !ok {
				room.Game.Cancel()
			} else {
				room.setPlay(room.PlayerO, play.Position)
			}
		}

		// send update status to the players
		if room.sendGameState(&room.PlayerX) != nil || room.sendGameState(&room.PlayerO) != nil {
			log.Printf("Failed sending initial game status")
			room.Game.Cancel()
			return
		}
	}

	// Update player ELOs
	// TODO: need to get player IDs of X and O to update on playerDB
	// playerdb.PlayerDB.UpdateAfterMatch(string(room.PlayerO), string(room.PlayerO), room.Game.GetState())
}

// sendGameState sends the current game state to the corresponding
// player.
func (room *Room) sendGameState(player *Player) error {
	select {
	case player.GameEvents <- room.getGameStateMsg():
		return nil
	case <-time.After(time.Second * 10):
		return fmt.Errorf("Timeout while sending game state to player '%v'", player.Play)
	}
}

// getGameStateMsg returns a message with the game state.
func (room *Room) getGameStateMsg() msg.GameState {
	return msg.GameState{
		Error:       "",
		Board:       string(room.Game.Board),
		State:       room.Game.GetState(),
		CurrentTurn: string(room.Game.GetCurrentTurn()),
	}
}

// setPlay handles the input action sent by the player, applies
// it to the game and sends the response to the player.
func (room *Room) setPlay(player Player, position int) error {
	if err := room.Game.AddPlay(player.Play, position); err != nil {
		player.GameEvents <- msg.GameState{Error: err.Error()}
	}
	return nil
}
