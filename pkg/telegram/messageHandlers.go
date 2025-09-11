package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shakareem/gigoseek/pkg/config"
)

const (
	startCommand = "start"
	authCommand  = "auth"
	helpCommand  = "help"
)

var responses = config.Get().Responses

// тут надо будет по chat id понимать, от кого сообщение,
// авторизирован ли пользователь и чота делать
func (b *Bot) handleMessage(msg *tgbotapi.Message) error {
	response := tgbotapi.NewMessage(msg.Chat.ID, "")

	if !msg.IsCommand() {
		response.Text = "пока поддерживаются только команды" //TODO: help responses
	}

	switch msg.Command() {
	case startCommand:
		// проверка, авторизован ли пользователь
		// ещё тут нужно сообщение с информацией о боте
		return b.handleAuth(msg.Chat.ID)
	case authCommand:
		return b.handleAuth(msg.Chat.ID)
	case helpCommand:
		response.Text = responses.Help
	default:
		response.Text = responses.UnknownCommand
	}

	_, err := b.bot.Send(response)
	return err
}

// эта функция перенаправляет пользователя на сервер авторизации
func (b *Bot) handleAuth(chatID int64) error {
	state := generateState()

	b.storage.SaveState(state, chatID)

	url := auth.AuthURL(state) // тут передаётся state, по которому потом можно понять, кто авторизовался
	response := tgbotapi.NewMessage(chatID, responses.AuthPrompt+url)

	_, err := b.bot.Send(response)
	return err
}
