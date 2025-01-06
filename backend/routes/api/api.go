package api

import (
	"encoding/json"
	"events/backend/database"
	"events/backend/database/gen"
	"events/backend/routes/api/types"
	"events/backend/routes/auth"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func Server() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/invite/{invite_code}", Invite)
	mux.HandleFunc("/event", Events)
	mux.HandleFunc("/event/{event_id}", Event)
	mux.HandleFunc("/event/{event_id}/organizer", EventOrganizers)
	mux.HandleFunc("/event/{event_id}/organizer/{org_id}", EventOrganizer)
	mux.HandleFunc("/event/{event_id}/participant", EventParticipants)
	mux.HandleFunc("/event/{event_id}/participant/{participant_id}", EventParticipant)
	mux.HandleFunc("/profile", Profile)

	return mux
}

func Invite(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetUserFromCookies(r)
	if err != nil {
		errorJson(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	inviteCode := r.PathValue("invite_code")
	if inviteCode == "" {
		errorJson(w, "missing invite code", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodGet {
		getInvite(w, r, user, inviteCode)
	} else {
		errorJson(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func getInvite(w http.ResponseWriter, r *http.Request, user *gen.User, inviteCode string) {
	// add the current user to an event as an organizer using the invite code and the method db.AddEventOrganizerByInviteCode
	ctx := r.Context()
	eventID, err := database.Default().AddEventOrganizerByInviteCode(ctx, gen.AddEventOrganizerByInviteCodeParams{
		InviteCode: inviteCode,
		UserID:     user.ID,
	})
	if err != nil {
		log.Printf("database error: %v", err)
		errorJson(w, "database error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/app/#/event/%d", eventID), http.StatusSeeOther)
}

func Events(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetUserFromCookies(r)
	if err != nil {
		errorJson(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodGet {
		getEvents(w, r, user)
	} else if r.Method == http.MethodPost {
		postEvents(w, r, user)
	} else {
		errorJson(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func getEvents(w http.ResponseWriter, r *http.Request, user *gen.User) {
	ctx := r.Context()
	dbEvents, err := database.Default().ListEvents(ctx, user.ID)
	if err != nil {
		log.Printf("database error: %v", err)
		errorJson(w, "database error", http.StatusInternalServerError)
		return
	}

	// transform database events to API events
	events := make([]types.Event, 0)
	for _, event := range dbEvents {
		events = append(events, types.Event{
			ID:          fmt.Sprint(event.ID),
			Name:        event.Name,
			Description: event.Description,
		})
	}

	// write response
	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		log.Println("ERROR: encoding response:", err)
	}
}

func postEvents(w http.ResponseWriter, r *http.Request, user *gen.User) {
	ctx := r.Context()
	var req types.Event
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errorJson(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Add the given event to the database
	eventID, err := database.Default().NewEvent(ctx, gen.NewEventParams{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		log.Printf("database error: %v", err)
		errorJson(w, "database error", http.StatusInternalServerError)
		return
	}
	err = database.Default().AddEventOrganizer(ctx, gen.AddEventOrganizerParams{
		EventID: eventID,
		UserID:  user.ID,
	})
	if err != nil {
		log.Printf("database error: %v", err)
		errorJson(w, "database error", http.StatusInternalServerError)
		return
	}

	// write response
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(types.Event{ID: fmt.Sprint(eventID), Name: req.Name, Description: req.Description})
	if err != nil {
		log.Println("ERROR: encoding response:", err)
	}
}

func userAndEvent(r *http.Request, w http.ResponseWriter) (user *gen.User, event_id int64, errorWritten bool) {
	user, err := auth.GetUserFromCookies(r)
	if err != nil {
		errorJson(w, "unauthorized", http.StatusUnauthorized)
		return nil, -1, true
	}

	event_id, err = getEventID(r)
	if err != nil {
		errorJson(w, err.Error(), http.StatusBadRequest)
		return nil, -1, true
	}
	return user, event_id, false
}

func Event(w http.ResponseWriter, r *http.Request) {
	user, event_id, errorWritten := userAndEvent(r, w)
	if errorWritten {
		return
	}

	if r.Method == http.MethodGet {
		getEvent(w, r, user, event_id)
	} else if r.Method == http.MethodPost {
		postEvent(w, r, user, event_id)
	} else {
		errorJson(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func getEvent(w http.ResponseWriter, r *http.Request, user *gen.User, event_id int64) {
	ctx := r.Context()
	events, err := database.Default().EventsByIds(ctx, gen.EventsByIdsParams{
		EventIds: []int64{event_id},
		UserID:   user.ID,
	})
	if err != nil {
		log.Printf("database error: %v", err)
		errorJson(w, "database error", http.StatusInternalServerError)
		return
	}
	if len(events) == 0 {
		errorJson(w, "event not found", http.StatusNotFound)
		return
	}
	event := events[0]

	// write response
	err = json.NewEncoder(w).Encode(types.Event{
		ID:          fmt.Sprint(event.ID),
		Name:        event.Name,
		Description: event.Description,
	})
	if err != nil {
		log.Println("ERROR: encoding response:", err)
	}
}

func postEvent(w http.ResponseWriter, r *http.Request, user *gen.User, event_id int64) {
	ctx := r.Context()

	// Parse request body
	var req types.Event
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errorJson(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Update the given event in the database
	event, err := database.Default().UpdateEvent(ctx, gen.UpdateEventParams{
		ID:          event_id,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		log.Printf("database error: %v", err)
		errorJson(w, "database error", http.StatusInternalServerError)
		return
	}

	// write response
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(types.Event{
		ID:          fmt.Sprint(event.ID),
		Name:        event.Name,
		Description: event.Description,
	})
	if err != nil {
		log.Println("ERROR: encoding response:", err)
	}
}

func Profile(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetUserFromCookies(r)
	if err != nil {
		errorJson(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodGet {
		getProfile(w, r, user)
	} else {
		errorJson(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func getProfile(w http.ResponseWriter, r *http.Request, user *gen.User) {
	err := json.NewEncoder(w).Encode(types.Profile{
		Email:      user.Email,
		Name:       user.Name,
		PictureURL: user.PictureUrl,
	})
	if err != nil {
		log.Println("ERROR: encoding response:", err)
	}
}

// Gets the event ID integer from the request path, requires the path to have a parameter named "event_id".
func getEventID(r *http.Request) (int64, error) {
	event_id_str := r.PathValue("event_id")
	if event_id_str == "" {
		return 0, fmt.Errorf("missing event_id")
	}
	event_id, err := strconv.ParseInt(event_id_str, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid event_id")
	}
	return event_id, nil
}

func EventParticipants(w http.ResponseWriter, r *http.Request) {
	user, event_id, errorWritten := userAndEvent(r, w)
	if errorWritten {
		return
	}

	if r.Method == http.MethodGet {
		getParticipants(w, r, user, event_id)
	} else {
		errorJson(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func getParticipants(w http.ResponseWriter, r *http.Request, user *gen.User, event_id int64) {
	ctx := r.Context()
	participants, err := database.Default().Participants(ctx, event_id)
	if err != nil {
		log.Printf("database error: %v", err)
		errorJson(w, "database error", http.StatusInternalServerError)
		return
	}

	// transform database participants to API participants
	var parts []types.Participant
	for _, part := range participants {
		parts = append(parts, types.Participant{
			ID:    fmt.Sprint(part.ID),
			Name:  part.Name.String,
			Email: part.Email,
		})
	}

	// write response
	err = json.NewEncoder(w).Encode(parts)
	if err != nil {
		log.Println("ERROR: encoding response:", err)
	}
}

func EventParticipant(w http.ResponseWriter, r *http.Request) {
	user, event_id, errorWritten := userAndEvent(r, w)
	if errorWritten {
		return
	}

	participant_id_str := r.PathValue("participant_id")
	if participant_id_str == "" {
		errorJson(w, "missing participant_id", http.StatusBadRequest)
		return
	}
	participant_id, err := strconv.ParseInt(participant_id_str, 10, 64)
	if err != nil {
		errorJson(w, "invalid participant_id", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodDelete {
		deleteParticipant(w, r, user, event_id, participant_id)
	} else {
		errorJson(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func deleteParticipant(w http.ResponseWriter, r *http.Request, user *gen.User, event_id, participant_id int64) {
	ctx := r.Context()
	err := database.Default().RemoveParticipant(ctx, participant_id)
	if err != nil {
		log.Printf("database error: %v", err)
		errorJson(w, "database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func EventOrganizers(w http.ResponseWriter, r *http.Request) {
	user, event_id, errorWritten := userAndEvent(r, w)
	if errorWritten {
		return
	}

	if r.Method == http.MethodGet {
		getOrganizers(w, r, user, event_id)
	} else {
		errorJson(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func getOrganizers(w http.ResponseWriter, r *http.Request, user *gen.User, event_id int64) {
	ctx := r.Context()
	organizers, err := database.Default().EventOrganizers(ctx, event_id)
	if err != nil {
		log.Printf("database error: %v", err)
		errorJson(w, "database error", http.StatusInternalServerError)
		return
	}

	// transform database organizers to API organizers
	var orgs []types.Organizer
	for _, org := range organizers {
		orgs = append(orgs, types.Organizer{
			ID:    fmt.Sprint(org.ID),
			Name:  org.Name,
			Email: org.Email,
		})
	}

	// write response
	err = json.NewEncoder(w).Encode(orgs)
	if err != nil {
		log.Println("ERROR: encoding response:", err)
	}
}

func EventOrganizer(w http.ResponseWriter, r *http.Request) {
	user, event_id, errorWritten := userAndEvent(r, w)
	if errorWritten {
		return
	}

	org_id_str := r.PathValue("org_id")
	if org_id_str == "" {
		errorJson(w, "missing org_id", http.StatusBadRequest)
		return
	}
	org_id, err := strconv.ParseInt(org_id_str, 10, 64)
	if err != nil {
		errorJson(w, "invalid org_id", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodDelete {
		deleteOrganizer(w, r, user, event_id, org_id)
	} else {
		errorJson(w, "method not allowed", http.StatusMethodNotAllowed)
	}
	errorJson(w, "method not allowed", http.StatusMethodNotAllowed)
}

func deleteOrganizer(w http.ResponseWriter, r *http.Request, user *gen.User, event_id, org_id int64) {
	ctx := r.Context()
	err := database.Default().RemoveEventOrganizer(ctx, gen.RemoveEventOrganizerParams{
		EventID: event_id,
		UserID:  org_id,
	})
	if err != nil {
		log.Printf("database error: %v", err)
		errorJson(w, "database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
