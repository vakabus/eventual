package api

import (
	"encoding/json"
	"events/backend/database"
	"events/backend/database/gen"
	"events/backend/email"
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
	mux.HandleFunc("/event/{event_id}/template", EventTemplates)
	mux.HandleFunc("/event/{event_id}/template/{template_id}", EventTemplate)
	mux.HandleFunc("/event/{event_id}/template/{template_id}/render", EventTemplateRender)
	mux.HandleFunc("/event/{event_id}/template/{template_id}/send", EventTemplateSend)
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
	} else if r.Method == http.MethodPost {
		postParticipants(w, r, user, event_id)
	} else {
		errorJson(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// Add a new participant
func postParticipants(w http.ResponseWriter, r *http.Request, user *gen.User, event_id int64) {
	var req types.Participant
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errorJson(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	jsonData, err := json.Marshal(req.Data)
	if err != nil {
		errorJson(w, "invalid request body", http.StatusBadRequest)
		return
	}
	participant, err := database.Default().AddParticipant(ctx, gen.AddParticipantParams{
		EventID: event_id,
		Json:    string(jsonData),
	})
	if err != nil {
		log.Printf("database error: %v", err)
		errorJson(w, "database error", http.StatusInternalServerError)
		return
	}

	apiParticipant, err := toAPIType(participant)
	if err != nil {
		log.Printf("error converting participant to API type: %v", err)
		errorJson(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(apiParticipant)
	if err != nil {
		log.Println("ERROR: encoding response:", err)
	}
}

func toDBType(p types.Participant, event_id int64) (gen.Participant, error) {
	jsonData, err := json.Marshal(p.Data)
	if err != nil {
		return gen.Participant{}, err
	}

	id, err := strconv.ParseInt(p.ID, 10, 64)
	if err != nil {
		return gen.Participant{}, fmt.Errorf("invalid ID: %v", err)
	}

	return gen.Participant{
		ID:      id,
		EventID: event_id,
		Json:    string(jsonData),
	}, nil
}

func toAPIType(p gen.Participant) (types.Participant, error) {
	var keyValueData map[string]string
	err := json.Unmarshal([]byte(p.Json), &keyValueData)
	if err != nil {
		return types.Participant{}, fmt.Errorf("deserialization of DB JSON failed: %v", err)
	}
	return types.Participant{
		ID:   fmt.Sprint(p.ID),
		Data: keyValueData,
	}, nil
}

func toAPIType2(p gen.ParticipantsRow) (types.Participant, error) {
	var keyValueData map[string]string
	err := json.Unmarshal([]byte(p.Json.(string)), &keyValueData)
	if err != nil {
		return types.Participant{}, fmt.Errorf("deserialization of DB JSON failed: %v", err)
	}
	return types.Participant{
		ID:   fmt.Sprint(p.ID),
		Data: keyValueData,
	}, nil
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
	var parts types.Participants
	parts.Participants = make([]types.Participant, 0)
	for _, part := range participants {
		p, err := toAPIType2(part)
		if err != nil {
			log.Printf("error converting participant to API type: %v", err)
			errorJson(w, "internal server error", http.StatusInternalServerError)
			return
		}

		parts.Participants = append(parts.Participants, p)
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
	} else if r.Method == http.MethodPost {
		postParticipant(w, r, user, event_id, participant_id)
	} else {
		errorJson(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func postParticipant(w http.ResponseWriter, r *http.Request, user *gen.User, event_id, participant_id int64) {
	// read request body
	ctx := r.Context()
	var req types.Participant
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errorJson(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// update the DB
	extraJson, err := json.Marshal(req.Data)
	if err != nil {
		errorJson(w, "invalid request body", http.StatusBadRequest)
		return
	}
	participant, err := database.Default().UpdateParticipant(ctx, gen.UpdateParticipantParams{
		ID:      participant_id,
		EventID: event_id,
		Json:    string(extraJson),
	})
	if err != nil {
		log.Printf("database error: %v", err)
		errorJson(w, "database error", http.StatusInternalServerError)
		return
	}

	// write response
	apiParticipant, err := toAPIType(participant)
	if err != nil {
		log.Printf("error converting participant to API type: %v", err)
		errorJson(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(apiParticipant)
	if err != nil {
		log.Println("ERROR: encoding response:", err)
	}
}

func deleteParticipant(w http.ResponseWriter, r *http.Request, user *gen.User, event_id, participant_id int64) {
	ctx := r.Context()
	err := database.Default().RemoveParticipant(ctx, gen.RemoveParticipantParams{
		EventID: event_id,
		ID:      participant_id,
	})
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

func EventTemplates(w http.ResponseWriter, r *http.Request) {
	user, event_id, errorWritten := userAndEvent(r, w)
	if errorWritten {
		return
	}

	if r.Method == http.MethodGet {
		getTemplates(w, r, user, event_id)
	} else if r.Method == http.MethodPost {
		postTemplates(w, r, user, event_id)
	} else {
		errorJson(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func postTemplates(w http.ResponseWriter, r *http.Request, user *gen.User, event_id int64) {
	ctx := r.Context()
	var req types.Template
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errorJson(w, "invalid request body", http.StatusBadRequest)
		return
	}

	template, err := database.Default().AddTemplate(ctx, gen.AddTemplateParams{
		EventID: event_id,
		Name:    req.Name,
		Body:    req.Body,
	})
	if err != nil {
		log.Printf("database error: %v", err)
		errorJson(w, "database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(types.Template{
		ID:   fmt.Sprint(template.ID),
		Name: template.Name,
		Body: template.Body,
	})
	if err != nil {
		log.Println("ERROR: encoding response:", err)
	}
}

func getTemplates(w http.ResponseWriter, r *http.Request, user *gen.User, event_id int64) {
	ctx := r.Context()
	templates, err := database.Default().Templates(ctx, event_id)
	if err != nil {
		log.Printf("database error: %v", err)
		errorJson(w, "database error", http.StatusInternalServerError)
		return
	}

	// transform database templates to API templates
	tpls := make([]types.Template, 0)
	for _, tpl := range templates {
		tpls = append(tpls, types.Template{
			ID:   fmt.Sprint(tpl.ID),
			Name: tpl.Name,
			Body: tpl.Body,
		})
	}

	// write response
	err = json.NewEncoder(w).Encode(tpls)
	if err != nil {
		log.Println("ERROR: encoding response:", err)
	}
}

func EventTemplate(w http.ResponseWriter, r *http.Request) {
	user, event_id, errorWritten := userAndEvent(r, w)
	if errorWritten {
		return
	}

	template_id_str := r.PathValue("template_id")
	if template_id_str == "" {
		errorJson(w, "missing template_id", http.StatusBadRequest)
		return
	}
	template_id, err := strconv.ParseInt(template_id_str, 10, 64)
	if err != nil {
		errorJson(w, "invalid template_id", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodGet {
		getTemplate(w, r, user, event_id, template_id)
	} else if r.Method == http.MethodPost {
		postTemplate(w, r, user, event_id, template_id)
	} else if r.Method == http.MethodDelete {
		deleteTemplate(w, r, user, event_id, template_id)
	} else {
		errorJson(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func getTemplate(w http.ResponseWriter, r *http.Request, user *gen.User, event_id, template_id int64) {
	ctx := r.Context()
	template, err := database.Default().Template(ctx, template_id)
	if err != nil {
		log.Printf("database error: %v", err)
		errorJson(w, "database error", http.StatusInternalServerError)
		return
	}

	// write response
	err = json.NewEncoder(w).Encode(types.Template{
		ID:   fmt.Sprint(template.ID),
		Name: template.Name,
		Body: template.Body,
	})
	if err != nil {
		log.Println("ERROR: encoding response:", err)
	}
}

func postTemplate(w http.ResponseWriter, r *http.Request, user *gen.User, event_id, template_id int64) {
	// parse request
	ctx := r.Context()
	var req types.Template
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errorJson(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// database query
	err = database.Default().UpdateTemplate(ctx, gen.UpdateTemplateParams{
		ID:   template_id,
		Name: req.Name,
		Body: req.Body,
	})
	if err != nil {
		log.Printf("database error: %v", err)
		errorJson(w, "database error", http.StatusInternalServerError)
		return
	}

	// write response
	w.WriteHeader(http.StatusNoContent)
}

func deleteTemplate(w http.ResponseWriter, r *http.Request, user *gen.User, event_id, template_id int64) {
	ctx := r.Context()
	err := database.Default().RemoveTemplate(ctx, template_id)
	if err != nil {
		log.Printf("database error: %v", err)
		errorJson(w, "database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func EventTemplateRender(w http.ResponseWriter, r *http.Request) {
	user, event_id, errorWritten := userAndEvent(r, w)
	if errorWritten {
		return
	}

	template_id_str := r.PathValue("template_id")
	if template_id_str == "" {
		errorJson(w, "missing template_id", http.StatusBadRequest)
		return
	}
	template_id, err := strconv.ParseInt(template_id_str, 10, 64)
	if err != nil {
		errorJson(w, "invalid template_id", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodGet {
		getTemplateRender(w, r, user, event_id, template_id)
	} else {
		errorJson(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func getTemplateRender(w http.ResponseWriter, r *http.Request, user *gen.User, event_id, template_id int64) {
	ctx := r.Context()

	// Get first participant's data
	participants, err := database.Default().Participants(ctx, event_id)
	if err != nil {
		log.Printf("database error: %v", err)
		errorJson(w, "database error", http.StatusInternalServerError)
		return
	}
	if len(participants) == 0 {
		errorJson(w, "no participants found", http.StatusNotFound)
		return
	}

	message, err := email.NewEmailFromDB(ctx, template_id, participants[0].ID)
	if err != nil {
		log.Printf("email error: %v", err)
		errorJson(w, "email error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message.BodyHTML()))
}

func EventTemplateSend(w http.ResponseWriter, r *http.Request) {
	user, event_id, errorWritten := userAndEvent(r, w)
	if errorWritten {
		return
	}

	template_id_str := r.PathValue("template_id")
	if template_id_str == "" {
		errorJson(w, "missing template_id", http.StatusBadRequest)
		return
	}
	template_id, err := strconv.ParseInt(template_id_str, 10, 64)
	if err != nil {
		errorJson(w, "invalid template_id", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodPost {
		sendTemplate(w, r, user, event_id, template_id)
	} else {
		errorJson(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func sendTemplate(w http.ResponseWriter, r *http.Request, user *gen.User, event_id, template_id int64) {
	ctx := r.Context()

	// Get first participant's data
	participants, err := database.Default().Participants(ctx, event_id)
	if err != nil {
		log.Printf("database error: %v", err)
		errorJson(w, "database error", http.StatusInternalServerError)
		return
	}
	if len(participants) == 0 {
		errorJson(w, "no participants found", http.StatusNotFound)
		return
	}

	message, err := email.NewEmailFromDB(ctx, template_id, participants[0].ID)
	if err != nil {
		log.Printf("email error: %v", err)
		errorJson(w, "email error", http.StatusInternalServerError)
		return
	}

	cookie, err := r.Cookie(auth.GoogleOauthTokenCookieName)
	if err != nil {
		errorJson(w, "missing google oauth token", http.StatusUnauthorized)
		return
	}

	err = message.SendGoogle(ctx, cookie.Value)
	if err != nil {
		log.Printf("email error: %v", err)
		errorJson(w, "email error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
