package router_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/richardpanda/quick-poll/server/choice"
	"github.com/richardpanda/quick-poll/server/httperror"
	"github.com/richardpanda/quick-poll/server/poll"
	"github.com/richardpanda/quick-poll/server/postgres"
	. "github.com/richardpanda/quick-poll/server/router"
	"github.com/richardpanda/quick-poll/server/vote"
	"github.com/richardpanda/quick-poll/server/ws"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestPostChoice(t *testing.T) {
	db := postgres.ConnectTest(t)
	defer db.Close()
	poll.CreateTable(db)
	choice.CreateTable(db)
	vote.CreateTable(db)
	defer poll.DropTable(db)
	defer choice.DropTable(db)
	defer vote.DropTable(db)

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

	choiceID := p.Choices[0].ID
	router := NewTestRouter(db, ws.NewConn())
	server := httptest.NewServer(router)
	wsURL := fmt.Sprintf("ws%s/v1/polls/%s/ws", strings.TrimPrefix(server.URL, "http"), p.ID)
	d := websocket.DefaultDialer
	conn, _, err := d.Dial(wsURL, nil)
	assert.NoError(t, err)
	defer conn.Close()

	choiceURL := fmt.Sprintf("%s/v1/polls/%s/choices/%s", server.URL, p.ID, choiceID)
	resp, err := http.Post(choiceURL, "", nil)
	assert.NoError(t, err)
	defer resp.Body.Close()

	var responseBody choice.PostChoiceResponseBody
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)

	cu := ws.ChoiceUpdate{}
	err = conn.ReadJSON(&cu)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, choiceID, responseBody.ID)
	assert.Equal(t, "blue", responseBody.Text)
	assert.Equal(t, 1, responseBody.NumVotes)
	assert.Equal(t, choiceID, cu.ID)
	assert.Equal(t, 1, cu.NumVotes)
}

func TestPostChoiceTwiceToSamePoll(t *testing.T) {
	db := postgres.ConnectTest(t)
	defer db.Close()
	poll.CreateTable(db)
	choice.CreateTable(db)
	vote.CreateTable(db)
	defer poll.DropTable(db)
	defer choice.DropTable(db)
	defer vote.DropTable(db)

	p := poll.Poll{
		ID:       uuid.NewV4().String(),
		Question: "Favorite color?",
		Choices: []choice.Choice{
			choice.New("blue"),
			choice.New("red"),
			choice.New("yellow"),
		},
		CheckIP: true,
	}
	err := db.Create(&p).Error
	assert.NoError(t, err)

	firstChoiceID := p.Choices[0].ID
	secondChoiceID := p.Choices[1].ID
	router := NewTestRouter(db, ws.NewConn())
	server := httptest.NewServer(router)
	wsURL := fmt.Sprintf("ws%s/v1/polls/%s/ws", strings.TrimPrefix(server.URL, "http"), p.ID)
	d := websocket.DefaultDialer
	conn, _, err := d.Dial(wsURL, nil)
	assert.NoError(t, err)
	defer conn.Close()

	choiceURL := fmt.Sprintf("%s/v1/polls/%s/choices/%s", server.URL, p.ID, firstChoiceID)
	resp, err := http.Post(choiceURL, "", nil)
	assert.NoError(t, err)
	defer resp.Body.Close()

	var responseBody choice.PostChoiceResponseBody
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)

	cu := ws.ChoiceUpdate{}
	err = conn.ReadJSON(&cu)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, firstChoiceID, responseBody.ID)
	assert.Equal(t, "blue", responseBody.Text)
	assert.Equal(t, 1, responseBody.NumVotes)
	assert.Equal(t, firstChoiceID, cu.ID)
	assert.Equal(t, 1, cu.NumVotes)

	choiceURL = fmt.Sprintf("%s/v1/polls/%s/choices/%s", server.URL, p.ID, secondChoiceID)
	resp, err = http.Post(choiceURL, "", nil)
	assert.NoError(t, err)
	defer resp.Body.Close()

	var errorResponseBody httperror.ResponseBody
	err = json.NewDecoder(resp.Body).Decode(&errorResponseBody)

	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, "You have already voted on this poll.", errorResponseBody.Message)
}

func TestPostChoiceWithInvalidID(t *testing.T) {
	db := postgres.ConnectTest(t)
	defer db.Close()
	poll.CreateTable(db)
	choice.CreateTable(db)
	vote.CreateTable(db)
	defer poll.DropTable(db)
	defer choice.DropTable(db)
	defer vote.DropTable(db)

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

	invalidID := uuid.NewV4()
	router := NewTestRouter(db, ws.NewConn())
	endpoint := fmt.Sprintf("/v1/polls/%s/choices/%s", p.ID, invalidID)
	req, err := http.NewRequest("POST", endpoint, nil)
	assert.NoError(t, err)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var responseBody httperror.ResponseBody
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)

	assert.NoError(t, err)
	assert.Equal(t, 400, resp.Code)
	assert.Equal(t, "Invalid choice ID.", responseBody.Message)
}
