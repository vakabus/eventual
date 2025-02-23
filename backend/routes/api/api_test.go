package api_test

import (
	"events/backend/database"
	"events/backend/database/gen"
	"events/backend/routes/api"
	"events/backend/routes/api/types"
	"fmt"
	"strconv"

	"events/backend/routes/auth"
	"flag"
	"net/http"
	"os"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	flag.Parse()
	database.DangerousSwitchDefaultToMemory()
	os.Exit(m.Run())
}

func setupSession(t *testing.T) string {
	user, err := database.Default().NewUser(t.Context(), gen.NewUserParams{
		TextID:     "test",
		Email:      "test@test.com",
		Name:       "Test User",
		PictureUrl: "https://example.com/test.png",
	})
	if err != nil {
		t.Fatal(err)
	}
	sessionToken, err := auth.NewSession(t.Context(), &user)
	if err != nil {
		t.Fatal(err)
	}
	return sessionToken
}

func TestAddParticipant(t *testing.T) {
	sessionToken := setupSession(t)

	// Create an event
	var event types.Event
	apitest.New().
		Handler(api.Server()).
		Post("/event").
		JSON(`{"name": "Test Event", "description": "Test Description"}`).
		Cookie(auth.SessionCookieName, sessionToken).
		Expect(t).
		Status(http.StatusCreated).
		End().
		JSON(&event)
	eventID, err := strconv.ParseInt(event.ID, 10, 64)
	if err != nil {
		t.Fatal(err)
	}

	// Get participants
	var participants types.Participants
	apitest.New().
		Handler(api.Server()).
		Get(fmt.Sprintf("/event/%s/participant", event.ID)).
		Cookie(auth.SessionCookieName, sessionToken).
		Expect(t).
		Status(http.StatusOK).
		End().
		JSON(&participants)

	// Check that there are no participants
	assert.Equal(t, 0, len(participants.Participants))
	assert.NotNil(t, participants.Participants)

	// Add an participant
	var participant types.Participant
	apitest.New().
		Handler(api.Server()).
		Post(fmt.Sprintf("/event/%s/participant", event.ID)).
		Cookie(auth.SessionCookieName, sessionToken).
		JSON(`{"data": {"name": "Test Participant", "email": "test@test.com"}}`).
		Expect(t).
		Status(http.StatusCreated).
		End().
		JSON(&participant)

	// Check that it is in DB
	dbParticipants, err := database.Default().Participants(t.Context(), eventID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(dbParticipants))

	// Get participants
	apitest.New().
		Handler(api.Server()).
		Get(fmt.Sprintf("/event/%s/participant", event.ID)).
		Cookie(auth.SessionCookieName, sessionToken).
		Expect(t).
		Status(http.StatusOK).
		End().
		JSON(&participants)

	// Test that we read the same participant as we added
	assert.Equal(t, 1, len(participants.Participants))
	assert.Equal(t, "Test Participant", participants.Participants[0].Data["name"])
	assert.Equal(t, "test@test.com", participants.Participants[0].Data["email"])
}
