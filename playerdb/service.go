package playerdb

import (
	"sync" 
	"fmt"
	"strings"
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

func (c *PlayerData) UpdateAfterMatch(player string, rune string, result string) {

	var playerElo = c.GetData(player)

	if strings.Contains(result, "Tied") {
		fmt.Println("Players Tied: no ELO update!")
		return
	}
	if strings.Contains(result, "won") && strings.Contains(result, rune) {
		c.Update(player, playerElo + 50)
	} else {
		c.Update(player, playerElo - 25)
	}

	fmt.Println("ELO UPDATE! ", player, ":", playerElo, " ---> ", c.GetData(player))
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