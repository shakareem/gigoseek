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
	server  *http.Server
	storage storage.Storage
}

type userInfo struct {
	chatID int64
	token  *oauth2.Token
}

func NewAuthServer(storage storage.Storage) *AuthServer {
	return &AuthServer{storage: storage}
}

func (s *AuthServer) Run() error {
	handler := http.NewServeMux()
	handler.HandleFunc("/callback", s.completeAuth)

	s.server = &http.Server{Addr: ":8080", Handler: handler}

	go func() {
		for userInfo := range userInfoChan {
			s.storage.SaveToken(userInfo.chatID, *userInfo.token)
			client := spotify.New(auth.Client(context.Background(), userInfo.token)) // подумать про контексты

			user, err := client.CurrentUser(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("New user logged in:", user.DisplayName, "\nchatID:", userInfo.chatID, "\ntoken:", userInfo.token.AccessToken)
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

	userInfoChan <- userInfo{chatID: chatID, token: token}

	// тут мб посылать вебхук боту

	http.Redirect(w, r, "https://web.telegram.org/k/#@gigoseek_bot", http.StatusSeeOther)
}
