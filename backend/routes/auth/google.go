package auth

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gOauth "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

var (
	googleOauthConfig = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "http://127.0.0.1:8080/auth/google-callback",
		Scopes: []string{
			gOauth.UserinfoEmailScope,
			gOauth.UserinfoProfileScope,
		},
		Endpoint: google.Endpoint,
	}
)

func initializeGoogleOAuth() {
	googleOauthConfig.ClientID = viper.GetString("google.client_id")
	googleOauthConfig.ClientSecret = viper.GetString("google.client_secret")
}

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	handleOauth2Login(w, r, googleOauthConfig, oauth2StateString)
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauth2StateString {
		log.Println("invalid oauth state, expected " + oauth2StateString + ", got " + state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	if code == "" {
		_, _ = w.Write([]byte("Code Not Found to provide AccessToken..\n"))
		reason := r.FormValue("error_reason")
		if reason == "user_denied" {
			_, _ = w.Write([]byte("User has denied permission.."))
		}
		// User has denied access..
		// http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	token, err := googleOauthConfig.Exchange(r.Context(), code)
	if err != nil {
		log.Println("ERROR: oauthConfGl.Exchange() failed with " + err.Error())
		return
	}

	userInfoService, err := gOauth.NewService(r.Context(), option.WithTokenSource(oauth2.StaticTokenSource(token)))
	if err != nil {
		log.Printf("ERROR: google NewService: %v", err)
		return
	}
	userInfo, err := userInfoService.Userinfo.Get().Do()
	if err != nil {
		log.Printf("ERROR: UserInfo.Get: %v", err)
		return
	}

	// Once we have all the data, store it in the database
	userID := fmt.Sprintf("google:%s", userInfo.Email)
	user, err := upsertUser(r.Context(), userID, userInfo.Email, userInfo.Name, userInfo.Picture)
	if err != nil {
		log.Printf("ERROR: upsertUser: %v", err)
		return
	}

	// And redirect to the main page
	err = AddSessionCookie(r.Context(), w, user)
	if err != nil {
		log.Printf("ERROR: AddSessionCookie: %v", err)
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
