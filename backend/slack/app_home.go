package slack

import (
	"strconv"
	"twinQuasarAppV2/backend/customTypes"
	"twinQuasarAppV2/backend/utils"

	"github.com/slack-go/slack"
)

func generateAppHomeContent(miners []customTypes.Miner) slack.HomeTabViewRequest {

	var blocks []slack.Block

	introText := slack.NewTextBlockObject("mrkdwn", "*Welcome to Filecoin Gravity !*\nThis bot notify you on chain events related to your miner like:\n- Mined block\n- Deadlines : passed, missed, late on\n- Sectors is reaching expiration", false, false)
	introImage := slack.NewImageBlockElement("https://tqslackbot.s3.eu-west-3.amazonaws.com/CIDgravity-Logo.png", "logo")
	introSection := slack.NewSectionBlock(introText, nil, slack.NewAccessory(introImage))

	divider := slack.NewDividerBlock()

	contactText := slack.NewTextBlockObject("mrkdwn", "```For any bugs or features request, feel free to reach us @JulienNoel @FlorianRuen```", false, false)
	contactSection := slack.NewSectionBlock(contactText, nil, nil)

	headerMiners := slack.NewTextBlockObject("plain_text", "My miners", false, false)
	headerMinersSection := slack.NewHeaderBlock(headerMiners, slack.HeaderBlockOptionBlockID("verified_miners_title"))

	limitText := slack.NewTextBlockObject("plain_text", "You are limited up to 3 three miners", true, false)
	limitSection := slack.NewContextBlock("", []slack.MixedElement{limitText}...)

	blocks = append(blocks, introSection, divider, contactSection, divider, headerMinersSection, limitSection)

	// For all the miners, push a specific block with specific values
	for _, miner := range miners {
		myMinerDescription := slack.NewTextBlockObject("mrkdwn", "<https://filfox.info/en/address/"+miner.Address.String+"|"+miner.Address.String+" (check on Filfox)>\nAdded on 02/02/20 10h42\nStatus: ✅ Verified", false, false)
		myMinerOverflowMenuTextOne := slack.NewTextBlockObject("plain_text", "Disable notifications", false, false)
		myMinerOverflowMenuTextTwo := slack.NewTextBlockObject("plain_text", "Remove", false, false)
		myMinerOverflowMenuOne := slack.NewOptionBlockObject("disableNotifications_"+miner.Address.String+"_"+strconv.Itoa(miner.Id)+"_"+utils.GenerateRandomHash(20), myMinerOverflowMenuTextOne, nil)
		myMinerOverflowMenuTwo := slack.NewOptionBlockObject("removeMiner_"+miner.Address.String+"_"+strconv.Itoa(miner.Id)+"_"+utils.GenerateRandomHash(20), myMinerOverflowMenuTextTwo, nil)
		myMinersOverflow := slack.NewOverflowBlockElement("miner-menu", myMinerOverflowMenuOne, myMinerOverflowMenuTwo)
		myMinersSection := slack.NewSectionBlock(myMinerDescription, nil, slack.NewAccessory(myMinersOverflow))
		blocks = append(blocks, myMinersSection, divider)
	}

	addNewMinerBtnText := slack.NewTextBlockObject("plain_text", "Add new miner", false, false)
	addNewMinerBtn := slack.NewButtonBlockElement("add-new-miner", "add-new-miner", addNewMinerBtnText)
	actionBlock := slack.NewActionBlock("action-block-btn-add", addNewMinerBtn)
	blocks = append(blocks, actionBlock)

	// Construct the blockSet with all the blocks before
	blockSet := slack.Blocks{
		BlockSet: blocks,
	}

	var homeTabRequest slack.HomeTabViewRequest
	homeTabRequest.Type = "home"
	homeTabRequest.Blocks = blockSet
	return homeTabRequest
}
