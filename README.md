# summary-slack-bot

## Features
- Summarizes recent slack messages using Google Gemini AI.
## Tech Stack

## Getting Started

1. Clone the repository:
```
https://github.com/alexandra1044/summary-slack-bot.git
```
```
ls summary-slack-bot
```
2. Download Dependencies

```
go mod download
```


### Prerequisites 

To run the slack-bot you need:

- Go installed
- Add appropriate API keys to a .env file within the program as shown below:

```
touch .env

echo "SLACK_CHANNEL_ID="<YOUR_SLACK_CHANNEL_ID_HERE> > .env

echo "SLACK_BOT_TOKEN="<YOUR_TOKEN_HERE>"" > .env

echo "GOOGLE_API_KEY="<YOUR_TOKEN_HERE>"" > .env
```