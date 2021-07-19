package slack

import (
	"fmt"
	"github.com/filecoin-project/lotus/api/v1api"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/slack-go/slack"
	"log"
	"net/http"
	"twinQuasarAppV2/backend/customTypes"
)

func Connect() (*slack.Client, customTypes.Config) {
	var config customTypes.Config

	errReadConfig := cleanenv.ReadConfig("config/config.yaml", &config)

	if errReadConfig != nil {
		log.Fatal(errReadConfig)
	}

	return slack.New(config.Slack.BotToken), config
}

func InitSlackBot(filecoinApi v1api.FullNodeStruct) {
	var slackApi, config = Connect()

	http.HandleFunc("/api/slack", func(w http.ResponseWriter, r *http.Request) {
		handleSlash(w, r, slackApi, config, filecoinApi)
	})

	http.HandleFunc("/api/slack/modals", func(w http.ResponseWriter, r *http.Request) {
		handleModal(w, r, slackApi, config, filecoinApi)
	})

	http.HandleFunc("/api/slack/events", func(w http.ResponseWriter, r *http.Request) {
		handleEvents(w, r, slackApi, config)
	})

	fmt.Println("[INFO] Slack Server listening on port 3001")
	log.Fatal(http.ListenAndServe(":3001", nil))
}

func SendMessageToUser(slackUserId string, messageToSend string) {
	var slackApi, _ = Connect()

	_, _, err := slackApi.PostMessage(
		slackUserId,
		slack.MsgOptionText(messageToSend, false),
		slack.MsgOptionAttachments(),
		slack.MsgOptionAsUser(false),
	)

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
}
