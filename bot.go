package main

import (
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

// replyText: 封裝 LINE 回覆文字的邏輯
func replyText(replyToken, text string) error {
	if _, err := bot.ReplyMessage(
		&messaging_api.ReplyMessageRequest{
			ReplyToken: replyToken,
			Messages: []messaging_api.MessageInterface{
				&messaging_api.TextMessage{
					Text: text,
				},
			},
		},
	); err != nil {
		return err
	}
	return nil
}

// callbackHandler: 處理 LINE 傳來的 Webhook 事件
func callbackHandler(w http.ResponseWriter, r *http.Request) {
	cb, err := webhook.ParseRequest(os.Getenv("ChannelSecret"), r)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	for _, event := range cb.Events {
		switch e := event.(type) {
		case webhook.MessageEvent:
			// 注意：這行最後面的 { 必須存在
			switch message := e.Message.(type) {
			case webhook.TextMessageContent:
				log.Println("收到訊息:", message.Text)
				
				// 呼叫 Gemini
				answer := gemini.GeminiChatComplete(message.Text)
				
				// --- 強化防崩潰邏輯開始 ---
				if answer == "" {
					log.Println("警告: Gemini 回傳空值，可能是 API Key 或參數錯誤")
					answer = "AI 目前沒有回應，請檢查 Gemini API Key 設定。"
				}
				// --- 強化防崩潰邏輯結束 ---

				if err := replyText(e.ReplyToken, answer); err != nil {
					log.Print("LINE 回傳失敗:", err)
				}
				
			case webhook.StickerMessageContent:
				replyText(e.ReplyToken, "收到你的貼圖了！")
			}
		}
	}
}
