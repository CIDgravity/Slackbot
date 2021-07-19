package slack

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"twinQuasarAppV2/backend/customTypes"
	"twinQuasarAppV2/backend/database"
	"twinQuasarAppV2/backend/filecoin"
	"twinQuasarAppV2/backend/utils"

	"github.com/filecoin-project/lotus/api/v1api"
	"github.com/slack-go/slack"
)

type CustomElements struct {
	Emoji bool   `json:"emoji"`
	Text  string `json:"text"`
	Type  string `json:"type"`
}

type CustomContextBlock struct {
	Type     string           `json:"type"`
	BlockID  string           `json:"block_id"`
	Elements []CustomElements `json:"elements"`
}

func generateModalRequestForSignature(minerAddress string, ownerPublicKey string, messageToSign string, isError bool, errMessage string) slack.ModalViewRequest {
	var blocks slack.Blocks

	titleText := slack.NewTextBlockObject("plain_text", "Verify "+minerAddress+" ownership", false, false)
	closeText := slack.NewTextBlockObject("plain_text", "Abort", false, false)
	submitText := slack.NewTextBlockObject("plain_text", "Verify", false, false)

	headerText := slack.NewTextBlockObject("mrkdwn", "To confirm you are the owner of the miner, follow this steps", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	divider := slack.NewDividerBlock()

	thirdStepTitleText := slack.NewTextBlockObject("mrkdwn", "*1. Execute this command to sign a message*", false, false)
	thirdStepTitle := slack.NewSectionBlock(thirdStepTitleText, nil, nil)

	commandHelpText := slack.NewTextBlockObject("plain_text", "lotus wallet sign "+ownerPublicKey+" "+messageToSign, true, false)
	commandHelp := slack.NewContextBlock("sign command", []slack.MixedElement{commandHelpText}...)

	errorText := slack.NewTextBlockObject("mrkdwn", ":exclamation: *Something wrong :* "+errMessage, false, false)
	errorSection := slack.NewSectionBlock(errorText, nil, nil)

	signedMessageText := slack.NewTextBlockObject(slack.PlainTextType, "2. Write the result of the above command here", false, false)
	signedMessagePlaceholder := slack.NewTextBlockObject(slack.PlainTextType, "Enter the message signed to launch the verification process", false, false)
	signedMessageElement := slack.NewPlainTextInputBlockElement(signedMessagePlaceholder, "signedMessage")
	signedMessageElement.MaxLength = 300
	signedMessageElement.Multiline = true
	signedMessage := slack.NewInputBlock("signed-message", signedMessageText, signedMessageElement)

	if isError {

		blocks = slack.Blocks{
			BlockSet: []slack.Block{
				headerSection, divider,
				thirdStepTitle, commandHelp, divider,
				signedMessage, errorSection,
			},
		}

	} else {

		blocks = slack.Blocks{
			BlockSet: []slack.Block{
				headerSection, divider,
				thirdStepTitle, commandHelp, divider,
				signedMessage,
			},
		}

	}

	var modalRequest slack.ModalViewRequest
	modalRequest.Type = "modal"
	modalRequest.Title = titleText
	modalRequest.ExternalID = "checkSignatureModal_" + utils.GenerateRandomHash(5)
	modalRequest.Close = closeText
	modalRequest.Submit = submitText
	modalRequest.Blocks = blocks
	return modalRequest
}

func generateModalRequestForMinerId(defaultValue string, isError bool, errMessage string) slack.ModalViewRequest {
	var blocks slack.Blocks

	titleText := slack.NewTextBlockObject("plain_text", "Add miner", false, false)
	closeText := slack.NewTextBlockObject("plain_text", "Abort", false, false)
	submitText := slack.NewTextBlockObject("plain_text", "Next", false, false)

	headerText := slack.NewTextBlockObject("mrkdwn", "Fill in your miner id. On the next step you will be asked to sign a message", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	errorText := slack.NewTextBlockObject("mrkdwn", ":exclamation: *Something wrong :* "+errMessage, false, false)
	errorSection := slack.NewSectionBlock(errorText, nil, nil)

	divider := slack.NewDividerBlock()

	minerIdText := slack.NewTextBlockObject(slack.PlainTextType, " ", false, false)
	minerIdPlaceholder := slack.NewTextBlockObject(slack.PlainTextType, "Your miner ID", false, false)
	minerIdElement := slack.NewPlainTextInputBlockElement(minerIdPlaceholder, "minerId")
	minerIdElement.MaxLength = 20
	minerIdElement.Multiline = false

	if defaultValue != "" {
		minerIdElement.InitialValue = defaultValue
	}

	minerId := slack.NewInputBlock("miner-id", minerIdText, minerIdElement)

	if isError {

		blocks = slack.Blocks{
			BlockSet: []slack.Block{headerSection, divider, minerId, errorSection},
		}

	} else {

		blocks = slack.Blocks{
			BlockSet: []slack.Block{headerSection, divider, minerId},
		}

	}

	var modalRequest slack.ModalViewRequest
	modalRequest.Type = "modal"
	modalRequest.Title = titleText
	modalRequest.ExternalID = "fillMinerIdModal_" + utils.GenerateRandomHash(5)
	modalRequest.Close = closeText
	modalRequest.Submit = submitText
	modalRequest.Blocks = blocks
	return modalRequest
}

func generateSuccessModalSignature(minerId string) slack.ModalViewRequest {
	var blocks slack.Blocks

	titleText := slack.NewTextBlockObject("plain_text", "Signature verified", false, false)
	closeText := slack.NewTextBlockObject("plain_text", "Finish", false, false)

	approvalText := slack.NewTextBlockObject("mrkdwn", "Your miner *"+minerId+"* has been successfully verified !\nNotifications for this miner are now activated.", false, false)
	approvalImage := slack.NewImageBlockElement("https://tqslackbot.s3.eu-west-3.amazonaws.com/approvedIcon.png", "approved")
	fieldsSection := slack.NewSectionBlock(approvalText, nil, slack.NewAccessory(approvalImage))

	blocks = slack.Blocks{
		BlockSet: []slack.Block{fieldsSection},
	}

	var modalRequest slack.ModalViewRequest
	modalRequest.Type = "modal"
	modalRequest.Title = titleText
	modalRequest.ExternalID = "succesSignatureProcess_" + utils.GenerateRandomHash(5)
	modalRequest.Close = closeText
	modalRequest.Blocks = blocks
	return modalRequest
}

func generateModalToConfirmRemoval(minerAddress string, minerId string) slack.ModalViewRequest {
	var blocks slack.Blocks

	titleText := slack.NewTextBlockObject("plain_text", "Remove miner "+minerId, false, false)
	closeText := slack.NewTextBlockObject("plain_text", "Abort", false, false)
	submitText := slack.NewTextBlockObject("plain_text", "Confirm", false, false)

	approvalText := slack.NewTextBlockObject("mrkdwn", "You are about to remove your miner from Filecoin Gravity. *This action can't be undone !*", false, false)
	approvalImage := slack.NewImageBlockElement("https://tqslackbot.s3.eu-west-3.amazonaws.com/delete.png", "delete")
	fieldsSection := slack.NewSectionBlock(approvalText, nil, slack.NewAccessory(approvalImage))

	blocks = slack.Blocks{
		BlockSet: []slack.Block{fieldsSection},
	}

	var modalRequest slack.ModalViewRequest
	modalRequest.Type = "modal"
	modalRequest.Title = titleText
	modalRequest.ExternalID = "confirmRemoveMiner_" + minerAddress + "_" + minerId + "_" + utils.GenerateRandomHash(5)
	modalRequest.Close = closeText
	modalRequest.Submit = submitText
	modalRequest.Blocks = blocks
	return modalRequest
}

func handleModal(w http.ResponseWriter, r *http.Request, slackApi *slack.Client, config customTypes.Config, filecoinApi v1api.FullNodeStruct) {

	err := verifySigningSecret(r, config.Slack.SigningSecret)

	if err != nil {
		fmt.Printf(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var i slack.InteractionCallback
	err = json.Unmarshal([]byte(r.FormValue("payload")), &i)

	if err != nil {
		fmt.Printf(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Handle clic on button on Home page (add miner, disable notification, remove miner for example)
	if len(i.ActionCallback.BlockActions) > 0 {

		if i.ActionCallback.BlockActions[0].ActionID == "add-new-miner" {

			fmt.Println("add miner btn detected")

			modalRequest := generateModalRequestForMinerId("", false, "")
			_, err := slackApi.OpenView(i.TriggerID, modalRequest)

			if err != nil {
				fmt.Printf("Error opening view: %s", err)
			}

		} else if strings.HasPrefix(i.ActionCallback.BlockActions[0].SelectedOption.Value, "disableNotifications_") {

			fmt.Println("DISABLE NOTIFICATIONS")

		} else if strings.HasPrefix(i.ActionCallback.BlockActions[0].SelectedOption.Value, "removeMiner_") {

			clickedButtonParams := strings.Split(i.ActionCallback.BlockActions[0].SelectedOption.Value, "_")
			modalRequestToConfirm := generateModalToConfirmRemoval(clickedButtonParams[1], clickedButtonParams[2])
			_, err := slackApi.OpenView(i.TriggerID, modalRequestToConfirm)

			if err != nil {
				fmt.Printf("Error opening view: %s", err)
			}
		}
	}

	// Here we also handle modals buttons and modal opened
	if strings.HasPrefix(i.View.ExternalID, "checkSignatureModal_") {

		commandToSignMessage, _ := json.Marshal(i.View.Blocks.BlockSet[3])
		var commandToSign customTypes.CustomContextBlock
		_ = json.Unmarshal(commandToSignMessage, &commandToSign)

		// Split the command to get every values inside
		commandToSignSplitted := strings.Fields(commandToSign.Elements[0].Text)
		ownerKey := commandToSignSplitted[3]
		messageToSign := commandToSignSplitted[4]

		// Get all data from the modal
		signedMessage := i.View.State.Values["signed-message"]["signedMessage"].Value

		// Get miner and all details from database
		minerId, minerAddress := database.GetMinerWithSignatureKey(messageToSign)

		if minerId == -1 {

			data := []byte(`{
				"response_action":"errors",
				"errors":{ "signed-message": "something wrong, try again in few minutes" }
			 }`)

			w.Header().Set("Content-Type", "application/json")
			w.Write(data)

		} else {

			result := filecoin.VerifySignature(filecoinApi, ownerKey, messageToSign, signedMessage)

			// Display an error in the modal or display the success modal if signature is valid
			if result == false {

				data := []byte(`{
					"response_action":"errors",
					"errors":{ "signed-message": "signed message is invalid, check it again to continue" }
				 }`)

				w.Header().Set("Content-Type", "application/json")
				w.Write(data)

			} else {

				database.UpdateMinerStatus(minerId)
				UpdateHomePage(slackApi, i.User.ID)

				response := customTypes.SlackEventCallback{
					ResponseAction: "update",
					View:           generateSuccessModalSignature(minerAddress),
				}

				js, err := json.Marshal(response)

				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.Write(js)

			}
		}

	} else if strings.HasPrefix(i.View.ExternalID, "fillMinerIdModal") {

		filledMinerId := i.View.State.Values["miner-id"]["minerId"].Value
		errGetOwnerKey, publicOwnerKey := filecoin.GetOwnerKey(filecoinApi, filledMinerId)

		if !errGetOwnerKey {

			messageToSign := hex.EncodeToString([]byte(utils.GenerateRandomHash(20)))
			errorInit := database.CheckToInitSignatureProcess(i.User.ID, i.Team.ID, filledMinerId, i.User.Name, publicOwnerKey, messageToSign)

			if errorInit == "already_subscribe" {

				data := []byte(`{
					"response_action":"errors",
					"errors":{ "miner-id": "You have already subscribe to this miner or you have reached your daily tries !" }
				 }`)

				w.Header().Set("Content-Type", "application/json")
				w.Write(data)

			} else if errorInit == "already_subscribe" {

				data := []byte(`{
					"response_action":"errors",
					"errors":{ "miner-id": "You have reached you usage limit for the miner subscribe !" }
				 }`)

				w.Header().Set("Content-Type", "application/json")
				w.Write(data)

			} else if errorInit == "proceed" {

				response := customTypes.SlackEventCallback{
					ResponseAction: "update",
					View:           generateModalRequestForSignature(filledMinerId, publicOwnerKey, messageToSign, false, ""),
				}

				js, err := json.Marshal(response)

				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.Write(js)

			}

		} else {

			data := []byte(`{
				"response_action":"errors",
				"errors":{ "miner-id": "miner ID not found on the network" }
			 }`)

			w.Header().Set("Content-Type", "application/json")
			w.Write(data)

		}

	} else if strings.HasPrefix(i.View.ExternalID, "confirmRemoveMiner") {

		validateModalParams := strings.Split(i.View.ExternalID, "_")
		isRemoved := database.RemoveMinerWithAddressAndSlackUser(validateModalParams[2])

		if isRemoved {

			UpdateHomePage(slackApi, i.User.ID)

		} else {

			fmt.Println("Error when removing the miner")

		}
	}
}

func UpdateHomePage(slackApi *slack.Client, userID string) {
	miners := database.GetAllMinersFromSlackUserId(userID)

	appHomeContent := generateAppHomeContent(miners)
	_, err := slackApi.PublishView(userID, appHomeContent, "")

	if err != nil {
		fmt.Printf("Error publishing app home view: %s", err)
		return
	}
}
