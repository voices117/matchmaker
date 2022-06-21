package bot

import (
	"matchmaker/game/msg"
	"matchmaker/game/room"
	"matchmaker/game/tictactoe"
)

// Bot that plays tic-tac-toe.
type Bot struct {
}

// Start makes the bot listen and respond to game events.
func (bot *Bot) Start(player *room.Player) {
	for gameState := range player.GameEvents {
		// check if it's the bot's turn
		if gameState.CurrentTurn == string(player.Play) {
			if gameState.State == tictactoe.STATE_UNFINISHED {
				player.PlayerEvents <- msg.SelectPosition{
					Position: bot.selectPosition(gameState),
				}
			}
		}
	}
}

// selectPosition returns a position to play for the bot.
func (bot *Bot) selectPosition(state msg.GameState) int {
	for i := range state.Board {
		if state.Board[i] == ' ' {
			return i
		}
	}
	// unreachable?
	panic("Bot didn't find an empty play!")
}
