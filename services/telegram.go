package services

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type TelegramClient struct {
	botToken string
	chatID   string
}

func NewTelegramClient(botToken, chatID string) *TelegramClient {
	return &TelegramClient{botToken: botToken, chatID: chatID}
}

func (t *TelegramClient) SendMessage(message string) error {
	endpoint := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.botToken)
	data := url.Values{
		"chat_id": {t.chatID},
		"text":    {message},
	}

	resp, err := http.PostForm(endpoint, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error enviando mensaje a Telegram: %s", resp.Status)
	}
	return nil
}
