package lobby

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
)

func (player *Player) StartPlayer(ctx context.Context, mm *MatchMaker, matchResponse chan *Match) {
	log.Printf("Start player %+v\n", player)

	for {
		select {
		case player2 := <-player.playersQueue:
			log.Printf("Player %v entering Case...\n", player.Id)
			if player2.Id == player.Id || !player.isValidMatch(player2) || !player2.isValidMatch(player) {
				log.Printf("Player %v leaving Case...\n", player.Id)
				continue
			}

			player.mtx.Lock()
			player2_available := player2.mtx.TryLock()

			if !player2_available {
				log.Printf("Player %v leaving Case 2...\n", player.Id)
				player.mtx.Unlock()
				continue
			}
			log.Println("Not Continuing...")

			if player.isWaiting && player2.isWaiting {
				match := Match{
					player1:  player,
					player2:  player2,
					GameRoom: uuid.NewString(),
				}

				log.Printf("Created match: %+v\n", match)

				player2.setIsInGame()
				player.setIsInGame()
				log.Printf("Set both players in match\n")
				log.Printf("Player2 status %+v\n", player2)
				select {
				case player2.matchQueue <- &match:
				case <-time.After(time.Second * 15):
					log.Panicf("Failed sending game Id to player '%v'", player2.Id)
				}
				log.Printf("player2.matchQueue <- &match\n")

				player2.mtx.Unlock()
				player.mtx.Unlock()
				log.Printf("Unlocked both players\n")
				select {
				case matchResponse <- &match:
				case <-time.After(time.Second * 15):
					log.Panicf("Failed sending match to matchResponse ")
				}
			} else {
				mm.Add(ctx, player2.Id)
				player2.mtx.Unlock()
				player.mtx.Unlock()
			}

			log.Printf("Player %v leaving Case correctly ...\n", player.Id)

		case matchedGame := <-player.matchQueue:
			select {
			case matchResponse <- matchedGame:
			case <-time.After(time.Second * 15):
				log.Printf("Failed sending match to response\n")
			}

		case <-time.After(time.Second * 30):
			log.Printf("Player %v Relaxing Requirements...\n", player.Id)
			player.relaxRequirements *= 1.03
			mm.Add(ctx, player.Id)
		}

	}
}

func (player *Player) isValidMatch(player2 *Player) bool {
	return Abs(player.elo-player2.elo) <= int(50*player.relaxRequirements)
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
