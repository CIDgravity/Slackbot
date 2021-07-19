package database

import (
	"context"
	"fmt"
	"time"
	"twinQuasarAppV2/backend/customTypes"

	"github.com/jackc/pgx/v4"
)

func CheckIfUserAlreadyExist(slackUserId string, slackTeamId string) bool {

	countExistingUsers := 0

	dBClient := Connect()
	err := dBClient.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM users u WHERE u.slack_user_id = $1 AND u.slack_team_id = $2",
		slackUserId, slackTeamId).Scan(&countExistingUsers)

	if err != nil || countExistingUsers <= 0 {
		return false
	} else {
		return true
	}
}

func CheckNotAlreadyRegisteredOrMaxTriesReached(slackUserId string, minerAddr string) bool {

	countValues := 0

	dBClient := Connect()
	err := dBClient.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM users u INNER JOIN miners m ON m.user_id = u.id WHERE u.slack_user_id = $1 AND "+
			"m.miner_address = $2 AND m.is_active = true",
		slackUserId, minerAddr).Scan(&countValues)

	if err != nil {
		return false
	}

	return countValues > 0
}

func CheckUsageLimitForSubscriptionNotReached(slackUserId string, slackTeamId string) bool {

	usageLimitForUser := 0
	countIterationForUser := 0

	dBClient := Connect()

	// First get usage limit for the specific user
	errFindLimit := dBClient.QueryRow(context.Background(),
		"SELECT value FROM users_has_usage_limits uhl INNER JOIN users u ON u.id = uhl.user_id "+
			"INNER JOIN usage_limits ul ON ul.id = uhl.usage_limit_id WHERE u.slack_user_id = $1 AND "+
			"u.slack_team_id = $2 AND ul.name = 'SUBSCRIBE_TO_MINER'",
		slackUserId, slackTeamId).Scan(&usageLimitForUser)

	if errFindLimit != nil {
		return false
	}

	// After count iteration and check if limit is reached
	errFindCount := dBClient.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM miners m INNER JOIN users u ON u.id = m.user_id "+
			"WHERE u.slack_user_id = $1 AND u.slack_team_id = $2 AND m.is_active = true",
		slackUserId, slackTeamId).Scan(&countIterationForUser)

	if errFindCount != nil {
		return false
	}

	return countIterationForUser >= usageLimitForUser
}

func AddUserBeginSigningProcess(slackUserId string, slackTeamId string, username string,
	minerAddr string, publicOwnerKey string, signProcessKey string) {
	dBClient := Connect()

	addedUserId := 0
	currentDateTime := time.Now().Format("01-02-2006 15:04:05")

	// First select user with same slack_user_id, slack_team_id and same username
	// If not found, add it into database, otherwise, get the entry id
	// With all where clauses, normally only one can be found if it's the case
	rows, err := dBClient.Query(context.Background(),
		"SELECT id FROM users WHERE slack_user_id = $1 AND slack_team_id = $2 AND username = $3",
		slackUserId, slackTeamId, username)

	if err != nil {
		panic(err)
	}

	defer rows.Close()
	var existingUserIds []int

	for rows.Next() {
		var currentUserId int
		errReadExistingUser := rows.Scan(&currentUserId)

		if errReadExistingUser != nil {
			panic(errReadExistingUser)
		}

		existingUserIds = append(existingUserIds, currentUserId)
	}

	// If not user found, add it
	// With all filters, we use only the first value if user already exist
	if len(existingUserIds) == 0 {
		errAddUser := dBClient.QueryRow(context.Background(),
			"INSERT INTO users(slack_user_id, slack_team_id, username, last_update) values($1, $2, $3, $4) RETURNING id",
			slackUserId, slackTeamId, username, currentDateTime).Scan(&addedUserId)

		if errAddUser != nil {
			fmt.Println("Error while adding user in database and returning ID")
			fmt.Println(errAddUser)
		} else {
			AddMinerToUser(addedUserId, minerAddr, publicOwnerKey, signProcessKey)
			InitUserLimitations(addedUserId)
		}
	} else {
		AddMinerToUser(existingUserIds[0], minerAddr, publicOwnerKey, signProcessKey)
	}
}

func InitUserLimitations(userId int) {
	dBClient := Connect()
	batch := &pgx.Batch{}

	rows, errReadUsageLimits := dBClient.Query(context.Background(), "SELECT * FROM usage_limits")

	if errReadUsageLimits != nil {
		panic(errReadUsageLimits)
	}

	defer rows.Close()
	var usageLimits []customTypes.UsageLimits

	for rows.Next() {
		var currentUsageLimit customTypes.UsageLimits
		errReadUsageLimits = rows.Scan(&currentUsageLimit.Id, &currentUsageLimit.Name, &currentUsageLimit.DefaultValue)

		if errReadUsageLimits != nil {
			panic(errReadUsageLimits)
		}

		usageLimits = append(usageLimits, currentUsageLimit)
	}

	// When all usage limits are retrieved, insert in junction table for specific user
	// By default, we use default value, but this can be overwrite (in the future can be in UI application)
	for _, usageLimit := range usageLimits {
		batch.Queue("INSERT INTO users_has_usage_limits(user_id, usage_limit_id, value) values($1, $2, $3)",
			userId, usageLimit.Id, usageLimit.DefaultValue)
	}

	br := dBClient.SendBatch(context.Background(), batch)
	_, errAssignLimitsToUser := br.Exec()

	if errAssignLimitsToUser != nil {
		fmt.Println("Error while assign usage limits to specific user")
		fmt.Println(errAssignLimitsToUser)
	}
}

func CheckToInitSignatureProcess(userId string, teamId string, minerId string, userName string, ownerKey string, messageToSign string) string {
	userAlreadyExist := CheckIfUserAlreadyExist(userId, teamId)

	if userAlreadyExist {
		usageLimitReached := CheckUsageLimitForSubscriptionNotReached(userId, teamId)
		alreadyRegistered := CheckNotAlreadyRegisteredOrMaxTriesReached(userId, minerId)

		if alreadyRegistered {
			return "already_subscribe"
		}

		if usageLimitReached {
			return "usage_limit"
		}
	}

	// All checks passed, proceed to signature verification
	AddUserBeginSigningProcess(userId, teamId, userName, minerId, ownerKey, messageToSign)
	return "proceed"
}
