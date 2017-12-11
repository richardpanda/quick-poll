package router_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/richardpanda/quick-poll/server/pkg/ws"

	"github.com/richardpanda/quick-poll/server/pkg/httperror"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/richardpanda/quick-poll/server/pkg/choice"
	"github.com/richardpanda/quick-poll/server/pkg/poll"
	. "github.com/richardpanda/quick-poll/server/pkg/router"
	"github.com/richardpanda/quick-poll/server/pkg/test"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestPOSTChoice(t *testing.T) {
	db, close := test.DBConnection(t)
	defer close()
	test.CreatePollsTable(db)
	test.CreateChoicesTable(db)
	defer test.DropPollsTable(db)
	defer test.DropChoicesTable(db)

	poll := poll.Poll{
		ID:       uuid.NewV4().String(),
		Question: "Favorite color?",
		Choices: []choice.Choice{
			choice.New("blue"),
			choice.New("red"),
			choice.New("yellow"),
		},
	}

	err := db.Create(&poll).Error
	assert.NoError(t, err)

	choiceID := poll.Choices[0].ID
	router := NewTestRouter(db, ws.NewConn())
	endpoint := fmt.Sprintf("/v1/choices/%s", choiceID)
	req, err := http.NewRequest("POST", endpoint, nil)
	assert.NoError(t, err)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var responseBody choice.POSTChoiceResponseBody
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, choiceID, responseBody.ID)
	assert.Equal(t, "blue", responseBody.Text)
	assert.Equal(t, 1, responseBody.NumVotes)
}

func TestPOSTChoiceWithInvalidID(t *testing.T) {
	db, close := test.DBConnection(t)
	defer close()
	test.CreatePollsTable(db)
	test.CreateChoicesTable(db)
	defer test.DropPollsTable(db)
	defer test.DropChoicesTable(db)

	poll := poll.Poll{
		ID:       uuid.NewV4().String(),
		Question: "Favorite color?",
		Choices: []choice.Choice{
			choice.New("blue"),
			choice.New("red"),
			choice.New("yellow"),
		},
	}

	err := db.Create(&poll).Error
	assert.NoError(t, err)

	invalidID := uuid.NewV4()
	router := NewTestRouter(db, ws.NewConn())
	endpoint := fmt.Sprintf("/v1/choices/%s", invalidID)
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
