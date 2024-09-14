package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"

	"github.com/slack-go/slack"
	"google.golang.org/api/option"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// fetch all tokens from env file
	slackBotToken := os.Getenv("SLACK_BOT_TOKEN")
	googleGeminiToken := os.Getenv("GOOGLE_API_KEY")
	slackChannelID := os.Getenv("SLACK_CHANNEL_ID")

	api := slack.New(slackBotToken)

	geminiSummary(googleGeminiToken, slackChannelID, slackBotToken, api)

}

func geminiSummary(googleGeminiToken string, slackChannelID string, slackBotToken string, api *slack.Client) {

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(googleGeminiToken))
	if err != nil {
		log.Panic("Error when generating gemini client")
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	// get appropriate amount of prior messages here

	params := slack.GetConversationHistoryParameters{
		ChannelID: slackChannelID,
	}

	messages, err := api.GetConversationHistoryContext(context.Background(), &params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	var totalMessages string
	for _, message := range messages.Messages {
		//fmt.Printf(message.Text)
		totalMessages = totalMessages + message.Text
	}

	prompt := "summarize these messages in a short paragraph" + totalMessages
	gen_response, err := model.GenerateContent(ctx, genai.Text(prompt))

	if err != nil {
		log.Fatal(err)
	}

	byte_val := getResponse(gen_response)

	// remove null terminator
	string_val := string(byte_val)
	trimmed_val := strings.Replace(string_val, `\n`, "", -1)

	printMessageToChat(trimmed_val, slackBotToken, slackChannelID)

}

func printMessageToChat(message string, token string, channelID string) {

	api := slack.New(token)
	attachment := slack.Attachment{
		Pretext: "",
		Text:    "",
	}

	channelID, timestamp, err := api.PostMessage(
		channelID,
		slack.MsgOptionText(message, false),
		slack.MsgOptionAttachments(attachment),
	)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)

}

func getResponse(resp *genai.GenerateContentResponse) []byte {
	var new_json []byte

	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				new_json, _ = json.MarshalIndent(part, "", "    ")
			}
		}
	}
	return new_json
}
