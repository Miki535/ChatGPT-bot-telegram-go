package main

import (
	"context"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
	"os"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func main() {

	botToken := ""

	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	updates, _ := bot.UpdatesViaLongPolling(nil)

	bh, _ := th.NewBotHandler(bot, updates)

	defer bh.Stop()
	defer bot.StopLongPolling()

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		chatId := tu.ID(update.Message.Chat.ID)

		keyboard := tu.Keyboard(
			tu.KeyboardRow(),
		)
		message := tu.Message(
			chatId,
			"Hello!",
		).WithReplyMarkup(keyboard)

		bot.SendMessage(message)

	}, th.CommandEqual("start"))

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		chatID := tu.ID(update.Message.Chat.ID)
		content := update.Message.Text
		reqToGPT(content, bot, chatID)
	}, th.AnyMessageWithText())

	bh.Start()
}

func reqToGPT(content string, bot *telego.Bot, chatId telego.ChatID) {
	client := openai.NewClient("")
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo,
			MaxTokens:   150,
			Temperature: 0,
			N:           1,
			Stop:        nil,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}
	bot.SendMessage(tu.Message(chatId, resp.Choices[0].Message.Content))

}
