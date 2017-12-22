package router_test

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/richardpanda/quick-poll/server/pkg/choice"
	"github.com/richardpanda/quick-poll/server/pkg/httperror"
	"github.com/richardpanda/quick-poll/server/pkg/poll"
	"github.com/richardpanda/quick-poll/server/pkg/postgres"
	. "github.com/richardpanda/quick-poll/server/pkg/router"
	"github.com/richardpanda/quick-poll/server/pkg/ws"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetWS(t *testing.T) {
	db := postgres.ConnectTest(t)
	defer db.Close()
	poll.CreateTable(db)
	choice.CreateTable(db)
	defer poll.DropTable(db)
	defer choice.DropTable(db)

	p := poll.Poll{
		ID:       uuid.NewV4().String(),
		Question: "Favorite color?",
		Choices: []choice.Choice{
			choice.New("blue"),
			choice.New("red"),
			choice.New("yellow"),
		},
	}
	err := db.Create(&p).Error
	assert.NoError(t, err)

	wsConn := ws.NewConn()
	server := httptest.NewServer(NewTestRouter(db, wsConn))
	defer server.Close()

	url := fmt.Sprintf("ws%s/v1/polls/%s/ws", strings.TrimPrefix(server.URL, "http"), p.ID)
	fmt.Println(url)
	d := websocket.DefaultDialer
	conn, _, err := d.Dial(url, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(wsConn.Table))

	clients, ok := wsConn.Table[p.ID]
	assert.True(t, ok)
	assert.Equal(t, 1, len(clients))

	conn.Close()
	<-wsConn.Done
	assert.Equal(t, 1, len(wsConn.Table))
	clients, ok = wsConn.Table[p.ID]
	assert.True(t, ok)
	assert.Equal(t, 0, len(clients))
}

func TestGetWSWithInvalidPollID(t *testing.T) {
	db := postgres.ConnectTest(t)
	defer db.Close()

	wsConn := ws.NewConn()
	server := httptest.NewServer(NewTestRouter(db, wsConn))
	defer server.Close()

	invalidPollID := uuid.NewV4().String()
	url := fmt.Sprintf("ws%s/v1/polls/%s/ws", strings.TrimPrefix(server.URL, "http"), invalidPollID)
	d := websocket.DefaultDialer
	conn, resp, err := d.Dial(url, nil)
	assert.Nil(t, conn)
	assert.Error(t, err)

	var responseBody httperror.ResponseBody
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	defer resp.Body.Close()

	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 400)
	assert.Equal(t, responseBody.Message, "Invalid poll ID.")
}
