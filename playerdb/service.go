package playerdb

import (
	"sync"
	"fmt" 
)

var (
	intialElo = 100
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

func (c *PlayerData) Len() int {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.playerList.
	defer c.mu.Unlock()
	return len(c.playerList)
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

type PlayerQueue struct {
	mu sync.Mutex
	playerList []string
}

func (c *PlayerQueue) Add(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, v := range c.playerList {
		if v == key {
			return false
		}
	}
	c.playerList = append(c.playerList, key)
	return true
}

func (c *PlayerQueue) Remove(key string) bool {
	// Tries to remove a player from the player Queue. 
	// Returns true if it could be removed else flase
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, v := range c.playerList {
		if v == key {
			c.playerList = append(c.playerList[:i], c.playerList[i+1:]...)
			return true
		}
	}
	return false
}