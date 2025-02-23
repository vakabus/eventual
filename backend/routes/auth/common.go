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
	"time"

	"golang.org/x/oauth2"
)

// Random state string for oauth2
var oauth2StateString = newToken()

func newToken() string {
	var bytes = make([]byte, 24)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

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

func upsertUser(ctx context.Context, userID string, email string, name string, pictureURL string) (*gen.User, error) {
	_, err := database.Default().UserByTextId(ctx, userID)
	if errors.Is(err, sql.ErrNoRows) {
		user, err := database.Default().NewUser(ctx, gen.NewUserParams{
			TextID:     userID,
			Email:      email,
			Name:       name,
			PictureUrl: pictureURL,
		})
		if err != nil {
			return nil, fmt.Errorf("add user: %v", err)
		}
		return &user, nil
	} else {
		user, err := database.Default().UpdateUser(ctx, gen.UpdateUserParams{
			TextID:     userID,
			Email:      email,
			Name:       name,
			PictureUrl: pictureURL,
		})
		if err != nil {
			return nil, fmt.Errorf("update user: %v", err)
		}
		return &user, nil
	}
}

func NewSession(ctx context.Context, user *gen.User) (string, error) {
	token := newToken()
	err := database.Default().AddSession(ctx, gen.AddSessionParams{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour).Format(time.DateTime),
	})
	if err != nil {
		return "", fmt.Errorf("add session: %v", err)
	}
	return token, nil
}

func AddSessionCookie(ctx context.Context, w http.ResponseWriter, user *gen.User) error {
	token, err := NewSession(ctx, user)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
	return nil
}

func GetUserFromCookies(r *http.Request) (*gen.User, error) {
	cookie, err := r.Cookie(SessionCookieName)
	if errors.Is(err, http.ErrNoCookie) {
		return nil, fmt.Errorf("no cookie found")
	}

	user, err := database.Default().UserBySession(r.Context(), cookie.Value)
	return &user, err
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
