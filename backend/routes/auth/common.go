package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"events/backend/database"
	"events/backend/database/gen"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
)

// Random state string for oauth2
var oauth2StateString = func() string {
	var bytes = make([]byte, 24)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}()

func handleOauth2Login(w http.ResponseWriter, r *http.Request, oauthConf *oauth2.Config, oauthStateString string) {
	URL, err := url.Parse(oauthConf.Endpoint.AuthURL)
	if err != nil {
		// this is a constant, so it should never fail
		panic(err)
	}

	parameters := url.Values{}
	parameters.Add("client_id", oauthConf.ClientID)
	parameters.Add("scope", strings.Join(oauthConf.Scopes, " "))
	parameters.Add("redirect_uri", oauthConf.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", oauthStateString)
	URL.RawQuery = parameters.Encode()
	url := URL.String()

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func upsertUser(ctx context.Context, userID string, email string, name string, pictureURL string) error {
	_, err := database.Default().UserByTextId(ctx, userID)
	if errors.Is(err, sql.ErrNoRows) {
		err = database.Default().NewUser(ctx, gen.NewUserParams{
			TextID:     userID,
			Email:      email,
			Name:       name,
			PictureUrl: pictureURL,
		})
		if err != nil {
			return fmt.Errorf("add user: %v", err)
		}
		return nil
	} else {
		err = database.Default().UpdateUser(ctx, gen.UpdateUserParams{
			TextID:     userID,
			Email:      email,
			Name:       name,
			PictureUrl: pictureURL,
		})
		if err != nil {
			return fmt.Errorf("update user: %v", err)
		}
		return nil
	}
}

// Hash-map to lookup users by their cookies
var session = make(map[string]string)

func registerSession(key string, userID string) {
	session[key] = userID
}

func GetUserFromCookies(r *http.Request) (*gen.User, error) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == googleCookieName {
			userID, ok := session[cookie.Value]
			if !ok {
				continue
			}

			user, err := database.Default().UserByTextId(r.Context(), userID)
			if err != nil {
				return nil, fmt.Errorf("database doesn't contain the stored user, is the data corrupted? %v", err)
			}
			return &user, nil
		}
	}

	return nil, fmt.Errorf("no matching user found")
}
