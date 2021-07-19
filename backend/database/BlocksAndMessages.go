package database

import (
	"context"
	"fmt"
	"twinQuasarAppV2/backend/customTypes"
)

func GetAllMessagesForTipSetAndCustomers(blockHeader customTypes.UniqueBlockHeader) []customTypes.UniqueBlockSingle {
	dBClient := Connect()

	rows, err := dBClient.Query(context.Background(),
		"SELECT bm.cid, bm.miner_to, bm.miner_from, bm.method_name, bm.actor, "+
			"bm.gas_fee_cap, bm.gas_limit, bm.gas_premium, bm.message_value, bm.status, bm.duplicate, u.id, "+
			"u.github_user, u.slack_user_id, u.slack_team_id, u.username FROM blocks b "+
			"INNER JOIN block_messages bm ON bm.block_cid = b.cid "+
			"INNER JOIN miners m ON m.miner_address = b.miner "+
			"INNER JOIN users u ON u.id = m.user_id "+
			"WHERE b.cid = $1 AND m.is_active = true", blockHeader.Cid.String)

	if err != nil {
		panic(err)
	}

	defer rows.Close()
	var messages []customTypes.UniqueBlockSingle

	for rows.Next() {
		var currentMessage customTypes.UniqueBlockSingle

		err = rows.Scan(&currentMessage.Message.Cid.String, &currentMessage.Message.To.String,
			&currentMessage.Message.From.String, &currentMessage.Message.MethodName.String,
			&currentMessage.Message.Actor.String, &currentMessage.Message.GasFeeCap,
			&currentMessage.Message.GasLimit, &currentMessage.Message.GasPremium, &currentMessage.Message.Value,
			&currentMessage.Message.Status, &currentMessage.Message.Duplicate, &currentMessage.User.Id,
			&currentMessage.User.GithubUser.String, &currentMessage.User.SlackUserId.String,
			&currentMessage.User.SlackTeamId.String, &currentMessage.User.SlackUsername.String)

		currentMessage.Header = blockHeader

		if err != nil {
			panic(err)
		}

		messages = append(messages, currentMessage)
	}

	return messages
}

func InsertChainMessagesBatch(blockHeaders []customTypes.UniqueBlockHeader, blockMessages []customTypes.UniqueBlockMessage) {

	dBClient := Connect()
	tx, errorTx := dBClient.Begin(context.Background())

	if errorTx != nil {
		fmt.Printf("Error when begin transaction with message : %s", errorTx)
		return
	}

	// Add all the headers, when sending the transactions, this will trigger function
	// to send every notifications
	for _, blockHeader := range blockHeaders {
		_, errorInsertHeader := tx.Exec(context.Background(),
			"INSERT INTO blocks(cid, height, block_timestamp, miner, is_validated, reward, win_count) "+
				"values($1, $2, $3, $4, $5, $6, $7)", blockHeader.Cid.String, blockHeader.Height.String,
			blockHeader.Timestamp.String, blockHeader.Miner.String, blockHeader.Validated, blockHeader.Reward.Uint64(),
			blockHeader.WinCount)

		if errorInsertHeader != nil {
			fmt.Printf("Error when insert header in transaction with message : %s", errorInsertHeader)
			return
		}
	}

	// Add all messages into a transaction, and after add the block header
	// PK is done on block CID so the order of add request doesn't matter
	for _, blockMessage := range blockMessages {
		_, errorInsertMessage := tx.Exec(context.Background(),
			"INSERT INTO block_messages(cid, miner_to, miner_from, method_name, actor, decoded_message, "+
				"gas_fee_cap, gas_limit, gas_premium, message_value, status, duplicate, block_cid, nonce) "+
				"values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)",
			blockMessage.Cid.String, blockMessage.To.String, blockMessage.From.String, blockMessage.MethodName.String,
			blockMessage.Actor.String, blockMessage.DecodedParams.String, blockMessage.GasFeeCap, blockMessage.GasLimit,
			blockMessage.GasPremium, blockMessage.Value, blockMessage.Status, blockMessage.Duplicate,
			blockMessage.BlockCid.String, blockMessage.Nonce)

		if errorInsertMessage != nil {
			fmt.Printf("Error when insert messages in transaction with message : %s", errorInsertMessage)
			return
		}
	}

	errCommit := tx.Commit(context.Background())

	if errCommit != nil {
		fmt.Printf("Error while commit the transaction with message : %s\n", errCommit)
		return
	}
}
