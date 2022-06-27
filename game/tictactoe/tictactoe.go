package tictactoe

import "fmt"

const (
	PLAY_O = 'O'
	PLAY_X = 'X'
)

const (
	STATE_TIE        = "Tied"
	STATE_WON_X      = "Player X won"
	STATE_WON_O      = "Player O won"
	STATE_CANCELLED  = "Cancelled"
	STATE_UNFINISHED = "Unfinished"
)

// TicTacToe represents a tic-tac-toe game state.
type TicTacToe struct {
	Board       []rune
	currentTurn rune
	state       string
	turnsPlayed int
}

// NewTicTacToe creates and initializes a new TicTacToe game instance.
func NewTicTacToe() TicTacToe {
	return TicTacToe{
		state:       STATE_UNFINISHED,
		turnsPlayed: 0,
		currentTurn: PLAY_X,
		Board:       []rune{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
	}
}

// Copy returns an independent copy of the tic-tac-toe instance.
func (t *TicTacToe) Copy() *TicTacToe {
	new := TicTacToe{
		state:       t.state,
		turnsPlayed: t.turnsPlayed,
		currentTurn: t.currentTurn,
		Board:       make([]rune, len(t.Board)),
	}
	copy(new.Board, t.Board)
	return &new
}

// AddPlay executes the logic of a player trying to add a play to
// the corresponding position in the board. The position must be
// a number between 0 and 8, as if the board was a single array.
// Play must be either `PLAY_X` or `PLAY_O`.
func (t *TicTacToe) AddPlay(play rune, position int) error {
	if play != PLAY_O && play != PLAY_X {
		return fmt.Errorf("Invalid play selected")
	}

	if play != t.currentTurn {
		return fmt.Errorf("Not your turn")
	}

	if t.state != STATE_UNFINISHED {
		return fmt.Errorf("Game already finished")
	}

	if t.isEntryOccupied(position) {
		return fmt.Errorf("Entry '%v' occupied. Select a different one", position)
	}

	if position < 0 || position > 8 {
		return fmt.Errorf("Column position must be [0, 8]")
	}

	t.Board[position] = play
	t.turnsPlayed += 1
	t.updateStatus()

	t.currentTurn = t.getNextTurnPlay()
	return nil
}

// Undo undoes the play in the selected position.
func (t *TicTacToe) Undo(position int) {
	t.turnsPlayed -= 1
	t.Board[position] = ' '
	t.state = STATE_UNFINISHED
	t.currentTurn = t.getNextTurnPlay()
}

// GetState returns the game's state.
func (t TicTacToe) GetState() string {
	return t.state
}

// GetCurrentTurn returns the player that should make the next move.
func (t TicTacToe) GetCurrentTurn() rune {
	return t.currentTurn
}

// GetCurrentTurn returns the player that should make the next move.
func (t *TicTacToe) Cancel() {
	t.state = STATE_CANCELLED
}

// updateStatus will update the game status based on the
// current board and player. It will check whether the game
// is a tie (board completed) or the current player won by
// completing a full row, column or diagonal.
func (t *TicTacToe) updateStatus() {
	if t.state != STATE_UNFINISHED {
		return
	}

	if t.PlayerWon(t.currentTurn) {
		if t.currentTurn == PLAY_X {
			t.state = STATE_WON_X
		} else {
			t.state = STATE_WON_O
		}
	} else if t.turnsPlayed == 9 {
		// if the board is full, then we call it a tie
		t.state = STATE_TIE
	}
}

// PlayerWon checks if a particular player won the game by
// completing a row, column or diagonal and returns `true if
// it did.
func (t TicTacToe) PlayerWon(play rune) bool {
	// check a win in the rows
	for row := 0; row < 3; row++ {
		if t.rowEquals(play, row) {
			return true
		}
	}

	// check a win in the columns
	for col := 0; col < 3; col++ {
		if t.colEquals(play, col) {
			return true
		}
	}

	// check diagonals
	// start by verifying if the middle entry is owned by this player
	if t.Board[4] == play {
		if t.Board[0] == play && t.Board[8] == play {
			return true
		}
		if t.Board[2] == play && t.Board[6] == play {
			return true
		}
	}

	// if it reached this point, then this player has not
	// completed a row
	return false
}

// rowEquals checks if a particular row in the 3x3 board is completed
// by the given player. The row must be an integer between 0 and 2
func (t TicTacToe) rowEquals(play rune, row int) bool {
	if row < 0 || row > 2 {
		panic("Invalid row number")
	}

	position := row * 3
	for i := 0; i < 3; i++ {
		if t.Board[position+i] != play {
			return false
		}
	}

	// at this point, all entries in the row are equal to the
	// expected play, so it's a hit
	return true
}

// colEquals checks if a particular column in the 3x3 board is completed
// by the given player. The column must be an integer between 0 and 2
func (t TicTacToe) colEquals(play rune, col int) bool {
	if col < 0 || col > 2 {
		panic("Invalid column number")
	}

	for i := 0; i < 3; i++ {
		if t.Board[col+3*i] != play {
			return false
		}
	}

	// at this point, all entries in the column are equal to the
	// expected play, so it's a hit
	return true
}

// isEntryOccupied will check if the specified position in the board
// is already occupied (i.e. already selected before by any player).
func (t *TicTacToe) isEntryOccupied(position int) bool {
	if t.Board[position] == ' ' {
		return false
	}
	return true
}

// getNextTurnPlay returns the player that should go in the next
// turn. The return value can be either PLAY_X or PLAY_O.
func (t TicTacToe) getNextTurnPlay() rune {
	if t.currentTurn == PLAY_O {
		return PLAY_X
	} else {
		return PLAY_O
	}
}
