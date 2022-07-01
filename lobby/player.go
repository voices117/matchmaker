package lobby

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

func (player *Player) StartPlayer(ctx context.Context, mm *MatchMaker, matchResponse chan interface{}) {
	log.Printf("Start player %+v\n", player)

	matchResponse <- fmt.Sprintf("Start matchmaking with ELO: %v", player.elo)

	defer close(matchResponse)

	for {
		select {
		case <-ctx.Done():
			return
		case player2 := <-player.playersQueue:
			log.Printf("Player %v entering Case...\n", player.Id)

			if player2.Id == player.Id {
				continue
			}
			if !player.isValidMatch(player2) || !player2.isValidMatch(player) {
				log.Printf("Player %v leaving Case...\n", player.Id)
				matchResponse <- fmt.Sprintf("Tried to match with %v (ELO %v) but was not accepted", player2.Id, player2.elo)
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
				log.Printf("Set player %v and %v in match\n", player.Id, player2.Id)
				select {
				case player2.matchQueue <- &match:
				case <-time.After(time.Second * 5):
					log.Panicf("Failed sending game Id to player '%v'", player2.Id)
				}

				player2.mtx.Unlock()
				player.mtx.Unlock()
				log.Printf("Unlocked both players\n")
				select {
				case matchResponse <- &match:
				case <-time.After(time.Second * 5):
					log.Panicf("Failed sending match to matchResponse ")
				}
				log.Printf("Player %v created game correctly\n", player.Id)
				return
			} else {
				mm.Add(ctx, player2.Id)
				player2.mtx.Unlock()
				player.mtx.Unlock()
			}
		case matchedGame := <-player.matchQueue:
			select {
			case matchResponse <- matchedGame:
			case <-time.After(time.Second * 5):
				log.Printf("Failed sending match to response\n")
			}
			return

		case <-time.After(time.Second * 5):
			log.Printf("Player %v Relaxing Requirements...\n", player.Id)
			player.relaxRequirements *= 1.03
			mm.Add(ctx, player.Id)
			matchResponse <- fmt.Sprintf("Relaxing matchmaking requirements: %v", player.relaxRequirements)
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
