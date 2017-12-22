package ws

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type ChoiceUpdate struct {
	ID       string `json:"id"`
	NumVotes int    `json:"num_votes"`
}

type Conn struct {
	Done  chan bool
	Table map[string]map[string]*websocket.Conn
}

func NewConn() *Conn {
	return &Conn{
		Done:  make(chan bool),
		Table: make(map[string]map[string]*websocket.Conn),
	}
}

func (c *Conn) BroadcastUpdate(pollID, choiceID string, numVotes int) {
	clients := c.Table[pollID]
	cu := ChoiceUpdate{ID: choiceID, NumVotes: numVotes}
	for _, conn := range clients {
		err := conn.WriteJSON(cu)
		if err != nil {
			fmt.Println(err)
		}
	}
}
