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
)

func Server() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/event", Event)
	mux.HandleFunc("/profile", Profile)

	return mux
}

func Event(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetUserFromCookies(r)
	if err != nil {
		errorJson(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodGet {
		getEvents(w, r, user)
	} else if r.Method == http.MethodPost {
		postEvent(w, r, user)
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
	var events []types.Event
	for _, event := range dbEvents {
		events = append(events, types.Event{
			ID:          fmt.Sprint(event.ID),
			Name:        event.Name,
			Description: event.Description,
		})
	}

	// optionally filter by event ID
	eventID := r.URL.Query().Get("id")
	if eventID != "" {
		var filteredEvents []types.Event
		for _, event := range events {
			if event.ID == eventID {
				filteredEvents = append(filteredEvents, event)
			}
		}
		events = filteredEvents
	}

	// write response
	err = json.NewEncoder(w).Encode(types.EventResponse{
		Events: events,
	})
	if err != nil {
		log.Println("ERROR: encoding response:", err)
	}
}

func postEvent(w http.ResponseWriter, r *http.Request, user *gen.User) {
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
