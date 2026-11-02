package telegram

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const apiBase = "https://api.telegram.org/bot"

type apiResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description"`
}

// SendPhoto отправляет картинку по ссылке с подписью.
func SendPhoto(token, chatID, photoURL, caption string) error {
	form := url.Values{}
	form.Set("chat_id", chatID)
	form.Set("photo", photoURL)
	form.Set("caption", caption)
	return call(token, "sendPhoto", form)
}

// SendMessage отправляет обычное текстовое сообщение.
func SendMessage(token, chatID, text string) error {
	form := url.Values{}
	form.Set("chat_id", chatID)
	form.Set("text", text)
	return call(token, "sendMessage", form)
}

// call — общий код для любого метода Telegram Bot API.
func call(token, method string, form url.Values) error {
	client := &http.Client{Timeout: 60 * time.Second}

	resp, err := client.PostForm(apiBase+token+"/"+method, form)
	if err != nil {
		return fmt.Errorf("telegram %s: %w", method, errors.Unwrap(err))
	}
	defer resp.Body.Close()

	var r apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return fmt.Errorf("telegram %s decode (status %d): %w", method, resp.StatusCode, err)
	}

	if !r.OK {
		return fmt.Errorf("telegram %s: %s", method, r.Description)
	}

	return nil
}
