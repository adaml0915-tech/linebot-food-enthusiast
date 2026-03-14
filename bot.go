package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

func replyText(replyToken, text string) error {
	if _, err := bot.ReplyMessage(
		&messaging_api.ReplyMessageRequest{
			ReplyToken: replyToken,
			Messages: []messaging_api.MessageInterface{
				&messaging_api.TextMessage{Text: text},
			},
		},
	); err != nil {
		return err
	}
	return nil
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	cb, err := webhook.ParseRequest(os.Getenv("ChannelSecret"), r)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	for _, event := range cb.Events {
		switch e := event.(type) {
		case webhook.MessageEvent:
			switch message := e.Message.(type) {
			case webhook.TextMessageContent:
				// 直接呼叫 Gemini 對話
				answer := gemini.GeminiChatComplete(message.Text)
				if err := replyText(e.ReplyToken, answer); err != nil {
					log.Print(err)
				}
			case webhook.StickerMessageContent:
				replyText(e.ReplyToken, "好可愛的貼圖！")
			}
		}
	}
}
