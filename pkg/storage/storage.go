package storage

import (
	"errors"
	"log"

	"golang.org/x/oauth2"
)

type Storage interface {
	SaveState(state string, chatID int64) error
	GetChatIDbyState(state string) (int64, error)
	DeleteState(state string) error

	SaveToken(chatID int64, token oauth2.Token) error
	GetToken(chatID int64) (oauth2.Token, error)
	DeleteToken(chatID int64) error

	SaveCity(chatID int64, city string) error
	GetCity(chatID int64) (string, error)
	DeleteCity(chatID int64) error
}

// tmp
type InMemoryStorage struct {
	states map[string]int64
	tokens map[int64]oauth2.Token
	cities map[int64]string
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		states: make(map[string]int64),
		tokens: make(map[int64]oauth2.Token),
		cities: make(map[int64]string),
	}
}

func (s *InMemoryStorage) SaveState(state string, chatID int64) error {
	s.states[state] = chatID
	return nil
}

func (s *InMemoryStorage) GetChatIDbyState(state string) (int64, error) {
	chatID, ok := s.states[state]
	if !ok {
		return 0, errors.New("state not found in storage")
	}
	return chatID, nil
}

func (s *InMemoryStorage) DeleteState(state string) error {
	delete(s.states, state)
	return nil
}

func (s *InMemoryStorage) SaveToken(chatID int64, token oauth2.Token) error {
	s.tokens[chatID] = token
	return nil
}

func (s *InMemoryStorage) GetToken(chatID int64) (oauth2.Token, error) {
	token, ok := s.tokens[chatID]
	if !ok {
		return oauth2.Token{}, errors.New("client not found in storage")
	}
	log.Printf("Token for chat %d retrieved successfully:\n%v", chatID, token)
	return token, nil
}

func (s *InMemoryStorage) DeleteToken(chatID int64) error {
	delete(s.tokens, chatID)
	return nil
}

func (s *InMemoryStorage) SaveCity(chatID int64, city string) error {
	s.cities[chatID] = city
	return nil
}

func (s *InMemoryStorage) GetCity(chatID int64) (string, error) {
	city, ok := s.cities[chatID]
	if !ok {
		return "", errors.New("city not found in storage")
	}
	return city, nil
}

func (s *InMemoryStorage) DeleteCity(chatID int64) error {
	delete(s.cities, chatID)
	return nil
}
