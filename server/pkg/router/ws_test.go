package router_test

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/richardpanda/quick-poll/server/pkg/httperror"

	"github.com/stretchr/testify/assert"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"

	. "github.com/richardpanda/quick-poll/server/pkg/router"
	"github.com/richardpanda/quick-poll/server/pkg/test"
	"github.com/richardpanda/quick-poll/server/pkg/ws"
)

func TestGetWS(t *testing.T) {
	db, close := test.DBConnection(t)
	defer close()

	wsConn := ws.NewConn()
	server := httptest.NewServer(NewTestRouter(db, wsConn))
	defer server.Close()

	id := uuid.NewV4()
	url := fmt.Sprintf("ws%s/v1/ws?poll_id=%s", strings.TrimPrefix(server.URL, "http"), id.String())
	d := websocket.DefaultDialer

	conn, _, err := d.Dial(url, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(wsConn.Table))

	clients, ok := wsConn.Table[id.String()]
	assert.True(t, ok)
	assert.Equal(t, 1, len(clients))

	conn.Close()
	<-wsConn.Done

	assert.Equal(t, 1, len(wsConn.Table))

	clients, ok = wsConn.Table[id.String()]
	assert.True(t, ok)
	assert.Equal(t, 0, len(clients))
}

func TestGetWSWithoutPollID(t *testing.T) {
	db, close := test.DBConnection(t)
	defer close()

	wsConn := ws.NewConn()
	server := httptest.NewServer(NewTestRouter(db, wsConn))
	defer server.Close()

	url := fmt.Sprintf("ws%s/v1/ws", strings.TrimPrefix(server.URL, "http"))
	d := websocket.DefaultDialer

	conn, resp, err := d.Dial(url, nil)

	assert.Nil(t, conn)
	assert.Error(t, err)

	var responseBody httperror.ResponseBody
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	defer resp.Body.Close()

	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 400)
	assert.Equal(t, responseBody.Message, "Poll ID is required.")
}
