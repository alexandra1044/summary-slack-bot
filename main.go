package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"github.com/shomali11/slacker"
	"google.golang.org/api/option"
)

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Command Events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
		fmt.Println()
	}
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	slackBotToken := os.Getenv("SLACK_BOT_TOKEN")
	slackAppToken := os.Getenv("SLACK_APP_TOKEN")
	googleGeminiToken := os.Getenv("GOOGLE_API_KEY")

	bot := slacker.NewClient(slackBotToken, slackAppToken)

	go printCommandEvents(bot.CommandEvents())

	bot.Command("summarize the last <number> messages and return text", &slacker.CommandDefinition{
		Description: "Uses Gemini AI to summarize the last X number of messages recieved",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			message_number := request.Param("number")
			fmt.Print(message_number)

			ctx := context.Background()
			client, err := genai.NewClient(ctx, option.WithAPIKey(googleGeminiToken))
			if err != nil {
				log.Panic("Error when generating gemini client")
			}
			defer client.Close()

			model := client.GenerativeModel("gemini-1.5-flash")

			prompt := "summarize the last" + message_number + "messages"
			gen_response, err := model.GenerateContent(ctx, genai.Text(prompt))

			if err != nil {
				log.Fatal(err)
			}

			response.Reply(fmt.Sprintf(gen_response))
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
