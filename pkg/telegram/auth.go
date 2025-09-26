package telegram

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/shakareem/gigoseek/pkg/config"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

var (
	redirectURL = config.Get().AuthServerURL
	auth        = spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURL),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserLibraryRead,
			spotifyauth.ScopeUserFollowRead,
			spotifyauth.ScopeUserTopRead,
		),
		spotifyauth.WithClientID(config.Get().SpotifyClientID),
		spotifyauth.WithClientSecret(config.Get().SpotifyClientSecret),
	)
	userInfoChan = make(chan userInfo)
)

func generateState() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return base64.URLEncoding.EncodeToString(b)
}

type AuthServer struct {
	server      *http.Server
	storage     Storage
	authUpdates chan<- int64
}

type userInfo struct {
	chatID int64
	token  *oauth2.Token
}

func NewAuthServer(storage Storage, authUpdates chan<- int64) *AuthServer {
	return &AuthServer{
		storage:     storage,
		authUpdates: authUpdates,
	}
}

func (s *AuthServer) Run() error {
	handler := http.NewServeMux()
	handler.HandleFunc("/", s.completeAuth)

	s.server = &http.Server{Addr: ":8080", Handler: handler}

	go func() {
		for userInfo := range userInfoChan {
			client := spotify.New(auth.Client(context.Background(), userInfo.token))

			user, err := client.CurrentUser(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("New user logged in:", user.DisplayName, "\nchatID:", userInfo.chatID, "\ntoken:", userInfo.token.AccessToken)
		}
	}()

	log.Println("Starting auth server")
	return s.server.ListenAndServeTLS("tmp/server.crt", "tmp/server.key") // TODO: move to config
}

func (s *AuthServer) completeAuth(w http.ResponseWriter, r *http.Request) {
	receivedState := r.FormValue("state")

	chatID, err := s.storage.GetChatIDbyAuthState(receivedState)
	if err != nil {
		log.Printf("HTTP request to auth server with invalid state: %v", receivedState)
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	token, err := auth.Token(r.Context(), receivedState, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Println("Failed to get token from state during auth")
	}

	s.storage.SaveToken(chatID, *token)
	s.storage.DeleteAuthState(receivedState)
	userInfoChan <- userInfo{chatID: chatID, token: token}

	s.authUpdates <- chatID

	http.Redirect(w, r, "https://web.telegram.org/k/#@gigoseek_bot", http.StatusSeeOther)
}
