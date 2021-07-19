package database

import (
	"context"
	"fmt"
	"twinQuasarAppV2/backend/customTypes"
)

func AddMinerToUser(userId int, minerAddr string, publicOwnerKey string, signProcessKey string) {
	dBClient := Connect()

	// For the signProcessKey we use the 10 first char of the messageToSign
	// It will use during the signature check process on fileCoin to identify both miner and user
	dBClient.QueryRow(context.Background(),
		"INSERT INTO miners(miner_address, owner_key, user_id, is_active, sign_process_key) values($1, $2, $3, false, $4)",
		minerAddr, publicOwnerKey, userId, signProcessKey)
}

func GetMinerWithSignatureKey(signatureKey string) (int, string) {
	dBClient := Connect()

	fmt.Println(signatureKey)

	minerId := -1
	minerAddress := ""

	err := dBClient.QueryRow(context.Background(), "SELECT miner_address, id FROM miners WHERE sign_process_key = $1", signatureKey).Scan(&minerAddress, &minerId)

	if err != nil {
		fmt.Println(err)
		return -1, ""
	} else {
		return minerId, minerAddress
	}
}

func UpdateMinerStatus(minerId int) {
	dBClient := Connect()
	dBClient.QueryRow(context.Background(), "UPDATE miners SET is_active = true WHERE id = $1", minerId)
}

func GetAllMinersFromSlackUserId(slackUserId string) []customTypes.Miner {

	dBClient := Connect()

	rows, err := dBClient.Query(context.Background(), "SELECT m.id, m.miner_address, m.owner_key, m.is_active, m.sign_process_key "+
		"FROM users u INNER JOIN miners m ON m.user_id = u.id WHERE u.slack_user_id = $1 AND m.is_active = true", slackUserId)

	if err != nil {
		panic(err)
	}

	defer rows.Close()
	var miners []customTypes.Miner

	for rows.Next() {
		var currentMiner customTypes.Miner

		err = rows.Scan(&currentMiner.Id, &currentMiner.Address, &currentMiner.OwnerKey,
			&currentMiner.IsActive, &currentMiner.SignProcessKey)

		if err != nil {
			panic(err)
		}

		miners = append(miners, currentMiner)
	}

	return miners
}

func RemoveMinerWithAddressAndSlackUser(minerID string) bool {

	dBClient := Connect()

	if _, err := dBClient.Exec(context.Background(), "DELETE FROM miners WHERE id = $1", minerID); err == nil {
		return true
	} else {
		return false
	}

}
