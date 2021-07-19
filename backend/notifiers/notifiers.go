package notifiers

import (
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/specs-actors/v5/actors/builtin"
	"twinQuasarAppV2/backend/customTypes"
	"twinQuasarAppV2/backend/slack"
)

func SlackNotifier(messages []customTypes.UniqueBlockSingle) {

	blockReward := big.Mul(messages[0].Header.Reward, big.NewInt(messages[0].Header.WinCount))
	blockReward = big.Div(blockReward, big.NewInt(builtin.ExpectedLeadersPerEpoch))
	divideBy1e18 := big.NewInt(1e18)
	gasFees := big.NewInt(0)

	for _, message := range messages {

		// For each non duplicate and with OK status messages, will add each gasFees
		// This allow us to calculate the final block reward for the miner
		if message.Message.Status && !message.Message.Duplicate {
			gasLimit := big.NewInt(message.Message.GasLimit)
			gasPremium := big.NewInt(message.Message.GasPremium)
			currentGasFees := big.Mul(gasLimit, gasPremium)
			gasFees = big.Add(gasFees, currentGasFees)
		}

		// For each message, treat them separately
		// v1: Proof of SpaceTime success or failed will be notified
		switch method := message.Message.MethodName.String; method {
		case "SubmitWindowedPoSt":
			slack.SendMessageToUser(message.User.SlackUserId.String, "Congrats! You just submit your PoSt successfully")
		}
	}

	// Calculate the final Reward
	blockReward = big.Add(blockReward, gasFees)
	blockRewardInFloat := big.Div(blockReward, divideBy1e18)

	// Print Reward original and Reward in float64 format
	// fmt.Printf("Reward in ATOFIL: %s\n", messages[0].Header.Reward.String())

	// Print the result of gasLimit * gasPremium in float64 for all messages
	// fmt.Printf("gasLimit + gasPremium for all messages : %s\n", gasFees.String())

	// Print the winCount value
	// fmt.Printf("winCount: %d\n", messages[0].Header.WinCount)

	// Print the result of the reward without division by 1e18
	// fmt.Printf("Total messages %d in block %s\n", len(messages), messages[0].Header.Cid.String)
	// fmt.Printf("Result reward with division by 1e18 : %f\n", big.Div(blockReward, divideBy1e18))

	// Send Slack notification for user that mine the current block
	slack.SendMessageToUser(messages[0].User.SlackUserId.String, "Congrats! You just mined a block and your reward was "+blockRewardInFloat.String())
}
