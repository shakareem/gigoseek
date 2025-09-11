package telegram

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/shakareem/gigoseek/pkg/config"
	"github.com/shakareem/gigoseek/pkg/storage"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

var (
	redirectURL = config.Get().AuthServerURL
	auth        = spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURL),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate),
		spotifyauth.WithClientID(config.Get().SpotifyClientID),
		spotifyauth.WithClientSecret(config.Get().SpotifyClientSecret),
	)
	tokenChan = make(chan *oauth2.Token)
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
	server  *http.Server
	storage storage.Storage
}

func NewAuthServer(storage storage.Storage) *AuthServer {
	return &AuthServer{storage: storage}
}

func (s *AuthServer) Run() error {
	handler := http.NewServeMux()
	handler.HandleFunc("/callback", s.completeAuth)

	s.server = &http.Server{Addr: ":8080", Handler: handler}

	go func() {
		for token := range tokenChan {
			// тут надо сохранять клиента в бд

			// use the token to get an authenticated client
			client := spotify.New(auth.Client(context.Background(), token)) // подумать про контексты

			user, err := client.CurrentUser(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("New user logged in:", user.ID, user.DisplayName)
		}
	}()

	return s.server.ListenAndServe()
}

func (s *AuthServer) completeAuth(w http.ResponseWriter, r *http.Request) {
	receivedState := r.FormValue("state")

	// Проверяем существование state и получаем chatID
	chatID, err := s.storage.GetChatIDbyState(receivedState)
	if err != nil {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	token, err := auth.Token(r.Context(), receivedState, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}

	s.storage.SaveToken(chatID, *token) // тут подумать про указатели

	tokenChan <- token
}
