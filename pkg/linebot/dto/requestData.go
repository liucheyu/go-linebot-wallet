package data

import(
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type requestData struct {
	Text string
	ReplyToken string
	Source *linebot.EventSource
}