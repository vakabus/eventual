package api

import (
	"context"
	"database/sql"
	"errors"
	"events/backend/api/types"
	"events/backend/database"
	"fmt"
	"log"
	"net/http"
)

func Server() http.Handler {
	return MagicRouter(&apiService{})
}

type apiService struct {
}

func (s *apiService) Auth(ctx context.Context, in types.AuthRequest) (*types.AuthResponse, error) {
	user, err := database.Default().UserByName(ctx, in.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &types.AuthResponse{
				ErrorResponse: types.ErrorResponse{
					ErrorMessage: "no such user found",
				},
			}, nil
		} else {
			log.Printf("database error: %v", err)
			return nil, fmt.Errorf("database error")
		}
	} else {
		response := types.AuthResponse{
			Token: fmt.Sprint(user.ID),
		}
		return &response, nil
	}
}
