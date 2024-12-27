package api

import (
	"context"
	"database/sql"
	"errors"
	"events/backend/database"
	"events/backend/routes/api/types"
	"fmt"
	"log"
	"net/http"
)

func Server() http.Handler {
	mux := http.NewServeMux()

	AddMagicRoutes(mux, &publicApiService{}, public)
	AddMagicRoutes(mux, &privateApiService{}, private)

	return mux
}

type publicApiService struct{}

type privateApiService struct{}

func (s *privateApiService) Dashboard(ctx context.Context, in types.DashboardRequest) (*types.DashboardResponse, error) {
	user, err := database.Default().UserByTextId(ctx, in.Token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &types.DashboardResponse{
				ErrorResponse: types.ErrorResponse{
					ErrorMessage: "invalid token",
				},
			}, nil
		} else {
			log.Printf("database error: %v", err)
			return nil, fmt.Errorf("database error")
		}
	}

	dbEvents, err := database.Default().ListEvents(ctx, user.ID)
	if err != nil {
		log.Printf("database error: %v", err)
		return nil, fmt.Errorf("database error")
	}

	var events []types.Event
	for _, event := range dbEvents {
		events = append(events, types.Event{
			ID:          string(event.ID),
			Name:        event.Name,
			Description: event.Description,
		})
	}

	response := types.DashboardResponse{
		Events: events,
	}
	return &response, nil
}
