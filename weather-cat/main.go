package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	bot, err := linebot.New(env.LINE_SCRET_TOKEN, env.LINE_ACCESS_TOKEN)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "LineBot SDK Init Failed.",
			StatusCode: 500,
		}, nil
	}
	var r http.Request
	r.Body = io.NopCloser(strings.NewReader(request.Body))
	r.Header = make(http.Header)
	r.Header.Add("X-Line-Signature", request.Headers["X-Line-Signature"])

	message_events, err := bot.ParseRequest(&r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			fmt.Println("ErrInvalidSignature")
		} else {
			fmt.Println("Bad Requst")
			fmt.Println(err)
		}
	}
	for _, event := range message_events {
		if event.Type == linebot.EventTypeFollow { //加好友
			if err := handleFollow(bot, event); err != nil {
				fmt.Printf("Handle Follow Failed , User %s", event.Source.UserID)
			}
		} else if event.Type == linebot.EventTypeMessage { //訊息
			if err := handleMessage(bot, event); err != nil {
				fmt.Printf("Handle Message Failed , User %s", event.Source.UserID)
			}
		}
	}
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Hello, %v", string("Sheep")),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}