package bot

import (
	"fmt"
	"matchmaker/game/tictactoe"
	"math/rand"
)

// MiniMax returns the optimal play for the player with the given
// `play` for the state of the `game`.
// Returns the position of the best next play or an error.
func MiniMax(current *tictactoe.TicTacToe, play rune) (int, error) {
	game := current.Copy()

	position, _, err := minimax(game, play, play, true)
	return position, err
}

// minimax executes the algorithm to find the best position to play next.
func minimax(game *tictactoe.TicTacToe, play rune, player rune, max bool) (position, score int, err error) {
	bestScore, bestPos := 0, -1
	scores, positions := []int{}, []int{}

	for i := range game.Board {
		if game.Board[i] == ' ' {
			// add a play to see how it goes...
			if err = game.AddPlay(play, i); err != nil {
				return -1, 0, err
			}

			if game.GetState() != tictactoe.STATE_UNFINISHED {
				// if the game finished, calculate the state based on the outcome
				score = calculateScore(game.GetState(), player)
			} else {
				// the game continues, so recursively calculates the score of the
				// other player's play
				_, score, err = minimax(game, game.GetCurrentTurn(), player, !max)
				if err != nil {
					return -1, 0, err
				}
			}

			if bestPos == -1 {
				bestScore, bestPos = score, i
			} else if max {
				if score > bestScore {
					bestScore, bestPos = score, i
				}
			} else {
				if score < bestScore {
					bestScore, bestPos = score, i
				}
			}

			scores, positions = append(scores, score), append(positions, i)

			// undo the temporal play
			game.Undo(i)
		}
	}

	// filter only best options
	bestScores, bestPositions := []int{}, []int{}
	for i := range scores {
		if scores[i] == bestScore {
			bestScores = append(bestScores, scores[i])
			bestPositions = append(bestPositions, positions[i])
		}
	}

	// now randomly select one of the best outcomes
	selection := rand.Intn(len(bestScores))
	return bestPositions[selection], bestScores[selection], nil
}

// calculateScore returns the score based on the game outcome
// assuming we are optimizing for the given `play`.
func calculateScore(state string, play rune) int {
	if state == tictactoe.STATE_TIE {
		return 0
	} else if state == tictactoe.STATE_WON_O {
		if play == tictactoe.PLAY_O {
			return 1
		} else {
			return -1
		}
	} else if state == tictactoe.STATE_WON_X {
		if play == tictactoe.PLAY_X {
			return 1
		} else {
			return -1
		}
	} else {
		panic(fmt.Sprintf("Unexpected state '%v'", state))
	}
}
