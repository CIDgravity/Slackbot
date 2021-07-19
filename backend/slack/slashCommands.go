package slack

import (
	"fmt"
	"net/http"
	"twinQuasarAppV2/backend/customTypes"

	"github.com/filecoin-project/lotus/api/v1api"
	"github.com/slack-go/slack"
)

func handleSlash(w http.ResponseWriter, r *http.Request, slackApi *slack.Client, config customTypes.Config, filecoinApi v1api.FullNodeStruct) {

	err := verifySigningSecret(r, config.Slack.SigningSecret)

	if err != nil {
		fmt.Printf(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	s, err := slack.SlashCommandParse(r)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}

	if s.ChannelName != "directmessage" {
		fmt.Println("Only available in private")
		return
	}

	switch s.Command {

	case "/add-miner":
		modalRequest := generateModalRequestForMinerId(s.Text, false, "")
		_, err := slackApi.OpenView(s.TriggerID, modalRequest)

		if err != nil {
			fmt.Printf("Error opening view: %s", err)
		}

	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
