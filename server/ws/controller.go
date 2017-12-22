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
		id := c.Params.ByName("id")

		conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		addrString := conn.RemoteAddr().String()

		if wsConn.Table[id] == nil {
			wsConn.Table[id] = make(map[string]*websocket.Conn)
		}
		wsConn.Table[id][addrString] = conn

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				delete(wsConn.Table[id], addrString)
				if len(wsConn.Table[id]) == 0 {
					delete(wsConn.Table, id)
				}
				if flag.Lookup("test.v") != nil {
					wsConn.Done <- true
				}
				return
			}
		}
	}
}
