package playerdb

import (
	"sync" 
	"fmt"
)

var (
	PlayerDB = PlayerData{
		playerList: map[string]int{},
	}
	intialElo = 1000
)


// SafeCounter is safe to use concurrently.
type PlayerData struct {
	mu sync.Mutex
	playerList map[string]int
}

func (c *PlayerData) Add(key string, value int) {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.playerList.
	// TODO: check for multiple user connections and defer errors
	c.playerList[key] = value
	c.mu.Unlock()
}

func (c *PlayerData) Update(key string, value int) {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.playerList.
	defer c.mu.Unlock()
	c.playerList[key] = value
}

func (c *PlayerData) UpdateAfterMatch(player0 string, playerX string, result string) {
	var elo0 = c.GetData(player0)
	var eloX = c.GetData(playerX)

	if result == "Player O won" {
		c.Update(player0, elo0 + 100)
		c.Update(playerX, eloX - 100)
	} else {
		c.Update(player0, elo0 - 100)
		c.Update(playerX, eloX + 100)
	}

	fmt.Println("%d , %d -> %d, %d", elo0, eloX, c.GetData(player0), c.GetData(playerX))
}

func (c *PlayerData) GetData(key string) int {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.playerList.
	defer c.mu.Unlock()
	if val, ok := c.playerList[key]; ok {
	    return val
	} 
	// Need to initialize player
	c.playerList[key] = intialElo
	return intialElo
}