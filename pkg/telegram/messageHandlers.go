package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const startCommand = "start"

// тут надо будет по chat id понимать, от кого сообщение,
// авторизирован ли пользователь и чота делать
func (b *Bot) HandleMessage(msg *tgbotapi.Message) error {
	response := tgbotapi.NewMessage(msg.Chat.ID, "")

	if !msg.IsCommand() {
		response.Text = "пока поддерживаются только команды" //TODO: help responses
	}

	switch msg.Command() {
	case startCommand:
		// init auth

		/* tmp */
		response.Text = b.responses.Start
	default:
		//TODO: responses
	}

	_, err := b.bot.Send(response)
	return err
}
