package ws

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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

func OpenConnection(wsConn *Conn) func(*gin.Context) {
	return func(c *gin.Context) {
		pollID := c.Query("poll_id")

		if pollID == "" {
			c.JSON(400, gin.H{"message": "Poll ID is required."})
			return
		}

		conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		addrString := conn.RemoteAddr().String()

		for _, connection := range wsConn.Table[pollID] {
			connection.WriteMessage(websocket.TextMessage, []byte(strconv.Itoa(len(wsConn.Table[pollID])+1)))
		}

		if wsConn.Table[pollID] == nil {
			wsConn.Table[pollID] = make(map[string]*websocket.Conn)
		}
		wsConn.Table[pollID][addrString] = conn

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				delete(wsConn.Table[pollID], addrString)
				if flag.Lookup("test.v") != nil {
					wsConn.Done <- true
				}
				return
			}
		}
	}
}
