package lobby

import (
	"context"
	"time"
)

func (player *Player) StartPlayer(ctx context.Context, mm *MatchMaker) Match {

	for {
		select {
		case <-time.After(time.Second * 30):
			player.relaxRequirements *= 1.03
			mm.Add(ctx, player.Id, player.responseQueue)

		case player2 := <-player.playersQueue:
			if player2.Id != player.Id && player.isValidMatch(player2) && player2.isValidMatch(player) {
				continue
			}

			player.mtx.Lock()
			player2_available := player2.mtx.TryLock()

			if !player2_available {
				player.mtx.Unlock()
				continue
			}

			// err = wsjson.Write(r.Context(), conn, player2)
			// if err != nil {
			// 	fmt.Printf("Player %v failed processing player %v to check if worthy candidate", id, player2.Id)
			// }

			if player.isWaiting && player2.isWaiting {
				match := Match{
					player1: player,
					player2: player2,
				}
				mm.createMatch(match)
				player2.setIsInGame()
				player.setIsInGame()
				player2.matchQueue <- &match

				player2.mtx.Unlock()
				player.mtx.Unlock()
				return match
			} else {
				mm.Add(ctx, player2.Id, player2.responseQueue)
			}

			player2.mtx.Unlock()
			player.mtx.Unlock()

		case matchedGame := <-player.matchQueue:
			return *matchedGame
		}
	}
}

func (player *Player) isValidMatch(player2 *Player) bool {
	return Abs(player.elo-int(float64(player2.elo)*player.relaxRequirements)) >= 50
}

func (player *Player) setIsInGame() {
	player.isWaiting = false
}

// Abs returns the absolute value of x.
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
