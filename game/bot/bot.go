package bot

import (
	"fmt"
	"matchmaker/game/msg"
	"matchmaker/game/room"
	"matchmaker/game/tictactoe"
)

// Bot that plays tic-tac-toe.
type Bot struct {
	Game *tictactoe.TicTacToe
}

// Start makes the bot listen and respond to game events.
func (bot *Bot) Start(player *room.Player) {
	for gameState := range player.GameEvents {
		// check if it's the bot's turn
		if gameState.CurrentTurn == string(player.Play) {
			if gameState.State == tictactoe.STATE_UNFINISHED {
				player.PlayerEvents <- msg.SelectPosition{
					Position: bot.selectPosition(gameState, player.Play),
				}
			}
		}
	}
}

// selectPosition returns a position to play for the bot.
func (bot *Bot) selectPosition(state msg.GameState, play rune) int {
	position, err := MiniMax(bot.Game, play)
	if err != nil {
		panic(fmt.Sprintf("Bot failed: %v", err))
	}

	return position
}
