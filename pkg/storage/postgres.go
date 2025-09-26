package storage

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/shakareem/gigoseek/pkg/config"
	"golang.org/x/oauth2"
)

type ChatState int

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage() (*PostgresStorage, error) {
	cfg := config.Get().Database

	psqlInfo := fmt.Sprintf(
		"host=%s port= %d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBname,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) SaveAuthState(state string, chatID int64) error {
	_, err := s.db.Exec(`
		INSERT INTO auth_state (state, chat_id)
		VALUES ($1, $2)
		ON CONFLICT (state) DO UPDATE SET chat_id = EXCLUDED.chat_id
	`, state, chatID)
	return err
}

func (s *PostgresStorage) GetChatIDbyAuthState(state string) (int64, error) {
	var chatID int64
	err := s.db.QueryRow(`
		SELECT chat_id FROM auth_state WHERE state = $1
	`, state).Scan(&chatID)
	log.Printf("get chat id by state err: %v", err)
	return chatID, err
}

func (s *PostgresStorage) DeleteAuthState(state string) error {
	_, err := s.db.Exec(`
		DELETE FROM auth_state WHERE state = $1
	`, state)
	return err
}

func (s *PostgresStorage) SaveToken(chatID int64, token oauth2.Token) error {
	_, err := s.db.Exec(`
		INSERT INTO token (chat_id, access_token, token_type, refresh_token, expiry)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (chat_id) DO UPDATE
		SET access_token = EXCLUDED.access_token,
		    token_type = EXCLUDED.token_type,
		    refresh_token = EXCLUDED.refresh_token,
		    expiry = EXCLUDED.expiry
	`, chatID, token.AccessToken, token.TokenType, token.RefreshToken, token.Expiry)
	return err
}

func (s *PostgresStorage) GetToken(chatID int64) (oauth2.Token, error) {
	var token oauth2.Token

	err := s.db.QueryRow(`
		SELECT access_token, token_type, refresh_token, expiry
		FROM token WHERE chat_id = $1
	`, chatID).Scan(&token.AccessToken, &token.TokenType, &token.RefreshToken, &token.Expiry)

	if err != nil {
		return oauth2.Token{}, err
	}

	return token, nil
}

func (s *PostgresStorage) DeleteToken(chatID int64) error {
	_, err := s.db.Exec(`
		DELETE FROM token WHERE chat_id = $1
	`, chatID)
	return err
}

func (s *PostgresStorage) SaveCity(chatID int64, city string) error {
	var cityID int
	err := s.db.QueryRow(`
		SELECT id FROM city WHERE city_name = $1
	`, city).Scan(&cityID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("Город %q не найден в базе", city)
		}
		return err
	}

	_, err = s.db.Exec(`
		UPDATE chat SET city_id = $1 WHERE chat_id = $2
	`, cityID, chatID)
	return err
}

func (s *PostgresStorage) GetCity(chatID int64) (string, error) {
	var city string
	err := s.db.QueryRow(`
		SELECT c.city_name
		FROM chat ch
		JOIN city c ON ch.city_id = c.id
		WHERE ch.chat_id = $1
	`, chatID).Scan(&city)
	return city, err
}

func (s *PostgresStorage) DeleteCity(chatID int64) error {
	_, err := s.db.Exec(`
		UPDATE chat SET city_id = NULL WHERE chat_id = $1
	`, chatID)
	return err
}

func (s *PostgresStorage) SaveChatState(chatID int64, state ChatState) error {
	_, err := s.db.Exec(`
		INSERT INTO chat (chat_id, chat_state)
		VALUES ($1, $2)
		ON CONFLICT (chat_id) DO UPDATE SET chat_state = EXCLUDED.chat_state
	`, chatID, state)
	return err
}

func (s *PostgresStorage) GetChatState(chatID int64) (ChatState, error) {
	var state ChatState
	err := s.db.QueryRow(`
		SELECT chat_state FROM chat WHERE chat_id = $1
	`, chatID).Scan(&state)
	return state, err
}

func (s *PostgresStorage) DeleteChatState(chatID int64) error {
	_, err := s.db.Exec(`
		UPDATE chat SET chat_state = 0 WHERE chat_id = $1
	`, chatID)
	return err
}
