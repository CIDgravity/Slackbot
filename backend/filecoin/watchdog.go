package filecoin

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/api/v1api"
	"github.com/filecoin-project/lotus/chain/stmgr"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/ipfs/go-cid"
	"strings"
	"time"
	"twinQuasarAppV2/backend/customTypes"
	"twinQuasarAppV2/backend/database"
	"twinQuasarAppV2/backend/utils"
)

func DecodeBlocks(api v1api.FullNodeStruct, chainHeadResult *types.TipSet, cacheMiners map[address.Address]types.Actor) {

	start := time.Now()

	var currentTipSetHeaders []customTypes.UniqueBlockHeader
	var currentTipSetMessages []customTypes.UniqueBlockMessage

	for _, blockCid := range chainHeadResult.Parents().Cids() {
		header, errGetHeader := api.ChainGetBlock(context.Background(), blockCid)
		messages, errGetMessages := api.ChainGetBlockMessages(context.Background(), blockCid)

		if errGetHeader != nil || header == nil {
			fmt.Printf("Error while retrieving header with message : %s\n", errGetHeader)
			return
		}

		if errGetMessages != nil || messages == nil {
			fmt.Printf("Error while retrieving messages with message : %s\n", errGetMessages)
			return
		}

		// Threat the blockHeader to get useful information about the miner and reward
		currentBlockHeader := customTypes.UniqueBlockHeader {
			Height: sql.NullString{ String : header.Height.String() },
			Miner: sql.NullString{ String : header.Miner.String() },
			Timestamp: sql.NullString{ String : time.Now().String() },
			Cid: sql.NullString{ String : blockCid.String() },
			Validated: header.IsValidated(),
			WinCount: header.ElectionProof.WinCount,
		}

		// Get the reward value for the current block CID
		rewardAddress, _ := address.NewFromString("f02")
		rewardData, errorGetReward := api.StateReadState(context.Background(), rewardAddress, chainHeadResult.Parents())

		if errorGetReward != nil {
			fmt.Println("Error when retrieving the reward with message" + errorGetReward.Error())
			return
		}

		var Reward customTypes.Rewards
		b, _ := json.Marshal(rewardData.State)

		if err := json.Unmarshal(b, &Reward); err != nil {
			fmt.Println(err)
			return
		}

		currentBlockHeader.Reward = Reward.ThisEpochReward
		currentTipSetHeaders = append(currentTipSetHeaders, currentBlockHeader)

		// Iterate over all message and analyse them separately
		for _, chainMessage := range messages.BlsMessages {
			messageCid := chainMessage.Cid()
			currentTipSetMessages = AnalyseBlockMessage(api, cacheMiners, currentTipSetMessages, chainHeadResult, *chainMessage, messageCid, blockCid.String())
		}

		for _, chainMessage := range messages.SecpkMessages {
			messageCid := chainMessage.Cid()
			currentTipSetMessages = AnalyseBlockMessage(api, cacheMiners, currentTipSetMessages, chainHeadResult, chainMessage.Message, messageCid, blockCid.String())
		}
	}

	// After all analyse process is over, call DB transaction to add all data inside tables
	// A separate micro-service will threat changes to send notification over Slack or e-mail
	database.InsertChainMessagesBatch(currentTipSetHeaders, currentTipSetMessages)

	elapsed := time.Since(start)
	fmt.Printf("Message decoded in %s\n", elapsed)
}

func WaitForChange(api v1api.FullNodeStruct, cacheMiners map[address.Address]types.Actor) {
	firstChainHead, _ := api.ChainHead(context.Background())

	for range time.Tick(1*time.Millisecond) {
		chainHead, errGetChainHead := api.ChainHead(context.Background())

		if errGetChainHead != nil {
			fmt.Printf("Error while checking the ChainHead with message : %s\n", errGetChainHead)
			return
		}

		if chainHead.Height() > firstChainHead.Height() {
			fmt.Printf("Height changed to %s !\n", chainHead.Height())
			time.Sleep(7 * time.Second)
			DecodeBlocks(api, chainHead, cacheMiners)
			firstChainHead = chainHead
		}
	}
}

func AnalyseBlockMessage(api v1api.FullNodeStruct, cacheMiners map[address.Address]types.Actor,
	currentTipSetMessages []customTypes.UniqueBlockMessage,
	chainHeadResult *types.TipSet, chainMessage types.Message, messageCid cid.Cid, blockCid string) []customTypes.UniqueBlockMessage {

	// Every messages will be add to the final list
	// But every messages that needs to be "Cancelled", status will be changed
	// For example if message is duplicate in the same TipSet or if old gasPremium < new gasPremium
	var actorFromMemory types.Actor
	actorFromMemory, actorFound := cacheMiners[chainMessage.To]

	if !actorFound {
		actorFromApi, _ := api.StateGetActor(context.Background(), chainMessage.To, chainHeadResult.Key())

		if actorFromApi != nil {
			actorFromMemory = *actorFromApi
			cacheMiners[chainMessage.To] = actorFromMemory
		}
	}

	if &actorFromMemory != nil {
		method := stmgr.MethodsMap[actorFromMemory.Code][chainMessage.Method]
		decodedMessage, _ := JsonParams(actorFromMemory.Code, chainMessage.Method, chainMessage.Params)
		actorSplit := strings.Split(string(actorFromMemory.Code.Hash()), "/")
		actorValue := "N/A"

		if len(actorSplit) >= 2 {
			actorValue = actorSplit[2]
		}

		currentBlockMessage := customTypes.UniqueBlockMessage{
			Cid: 		   sql.NullString{ String: messageCid.String() },
			DecodedParams: sql.NullString{ String: decodedMessage },
			MethodName:    sql.NullString{ String: method.Name },
			Actor:         sql.NullString{ String: actorValue },
			From:          sql.NullString{ String: chainMessage.From.String() },
			Nonce: 		   chainMessage.Nonce,
			To:            sql.NullString{ String: chainMessage.To.String() },
			GasFeeCap:     chainMessage.GasFeeCap.Int64(),
			GasLimit:      chainMessage.GasLimit,
			GasPremium:    chainMessage.GasPremium.Int64(),
			Value:         chainMessage.Value.Int64(),
			BlockCid:      sql.NullString{ String: blockCid },
			Duplicate:     false,
			Status:        true,
		}

		// Test the message to check if duplicate or need to update his status
		// If messageExist, the status of the new one need to set to false
		// If messageExist but gasPremium < gasPremium of new one, update old message status to false
		// If doesn't exist, add the new one with status true
		messageExist, messageIndex := utils.CheckIfExist(currentBlockMessage, currentTipSetMessages)

		if messageExist {
			if chainMessage.GasPremium.Int64() > currentTipSetMessages[messageIndex].GasPremium {
				currentBlockMessage.Status = false
			} else {
				currentBlockMessage.Duplicate = true
			}
		}

		currentTipSetMessages = append(currentTipSetMessages, currentBlockMessage)
	}

	return currentTipSetMessages
}