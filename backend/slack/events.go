package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"twinQuasarAppV2/backend/customTypes"
	"twinQuasarAppV2/backend/database"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func handleEvents(w http.ResponseWriter, r *http.Request, slackApi *slack.Client, config customTypes.Config) {

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sv, err := slack.NewSecretsVerifier(r.Header, config.Slack.SigningSecret)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := sv.Write(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := sv.Ensure(); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	eventsAPIEvent, err := slackevents.ParseEvent(body, slackevents.OptionNoVerifyToken())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Event : Challenge received for url verification
	if eventsAPIEvent.InnerEvent.Type == "app_home_opened" {

		// Get all the informations from the JSON event (including slack_user_id)
		eventData, _ := eventsAPIEvent.InnerEvent.Data.(*slackevents.AppHomeOpenedEvent)

		// Get all the miner associated with the user who opened the Home tab
		miners := database.GetAllMinersFromSlackUserId(eventData.User)

		// Generate the app home and send it to the user
		appHomeContent := generateAppHomeContent(miners)
		_, err = slackApi.PublishView(eventData.User, appHomeContent, "")

		if err != nil {
			fmt.Printf("Error publishing app home view: %s", err)
			return
		}

	} else if eventsAPIEvent.InnerEvent.Type == "url_verification" {

		var r *slackevents.ChallengeResponse
		err := json.Unmarshal(body, &r)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text")
		w.Write([]byte(r.Challenge))

	}
}
