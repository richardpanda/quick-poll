package ws

import (
	"flag"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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
