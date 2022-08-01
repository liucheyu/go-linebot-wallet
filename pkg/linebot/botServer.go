package linebot

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"

	"github.com/go-redis/redis/v8"
	"github.com/go-sql-driver/mysql"

	// "github.com/liucheyu/go-linebot-wallet/pkg/linebot/database"

	"github.com/liucheyu/go-linebot-wallet/pkg/linebot/handler"
	"github.com/liucheyu/go-linebot-wallet/pkg/linebot/service"
)

var (
	actionRe, _        = regexp.Compile("^#action:(\\d{2}),(\\w{1,10})$")
	numberic, _        = regexp.Compile("^\\d+")
	ctx                context.Context
	redisExpireSeconds time.Duration = 60
	linebotClient      *linebot.Client
)

func init() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	config := mysql.Config{
		User:                 "root",
		Passwd:               "123456789",
		Addr:                 "localhost:13306",
		Net:                  "tcp",
		DBName:               "linebot_wallet",
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		fmt.Println(err)
	}

	// 釋放連線
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}

	linebotClient = newLineBotClient()
	handler.Bot = linebotClient;
	ctx = context.Background()
	
	handler.UserService = &service.UserActionServiceImp{Context: ctx, Redis: rdb, DB:db}
	handler.BaseDataCache = service.BaseDataCache{Context: ctx, Redis: rdb, DB:db}
	
}

func Callback(c *gin.Context) {
	botServer(c.Request)
}

func newLineBotClient() *linebot.Client {
	bot, err := linebot.New("81c6a13588d5c0b4e5719fd7e34da9b9", "QMetcOLoLQzjvBfNBO0jxxHCAzrzhQu4lreovVEfPDUXxrol6m/a/PgAd9oXv9lHYn+r3nq9OgLWhNs6ZAtLhNN/7iRO1kro4hMJ/Drngv8Cb1iQciXM11vTFBPQVL4YyBlY+IilWyKD8FR5aqINqgdB04t89/1O/w1cDnyilFU=")
	if err != nil {
		fmt.Println(err)
	}
	return bot
}

func botServer(lineBotRquest *http.Request) {
	var requestText string	
	events, _ := linebotClient.ParseRequest(lineBotRquest)

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			textMsg := event.Message.(*linebot.TextMessage)
			requestText = textMsg.Text
		}

		if event.Type == linebot.EventTypePostback {
			requestText = event.Postback.Data
		}		

		context := handler.Context{HandlerName: "", UserID: event.Source.UserID, Token: event.ReplyToken, RequestText: requestText}

		if handlerName := handler.GetNextHandler(event.Source.UserID);handlerName != "" {
			context.HandlerName = handlerName
		} else {
			context.HandlerName = handler.MappingFirstStat(requestText)
		}

		handler.ExcuteHandler(&context)	
	}

}
