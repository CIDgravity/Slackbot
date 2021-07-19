package triggers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/lib/pq"
	"time"
	"twinQuasarAppV2/backend/customTypes"
	"twinQuasarAppV2/backend/database"
	"twinQuasarAppV2/backend/notifiers"
)

func InitDatabaseTriggers(config customTypes.Config) {
	var databaseUrl = config.Database.EndpointWithCredentials + config.Database.DatabaseName

	_, err := sql.Open("postgres", databaseUrl)

	if err != nil {
		panic(err)
	}

	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	listener := pq.NewListener(databaseUrl, 10*time.Second, time.Minute, reportProblem)
	err = listener.Listen("events")

	if err != nil {
		panic(err)
	}

	go func() {
		fmt.Println("[INFO] Start monitoring PostgreSQL")
		waitForNotification(listener)
	}()
}

func waitForNotification(l *pq.Listener) {
	for {
		select {
		case n := <-l.Notify:
			eventData := customTypes.AddedBlockEvent{}
			err := json.Unmarshal([]byte(n.Extra), &eventData)

			if err != nil {
				fmt.Println(err)
				return
			}

			rewardAsBigInt, _ := types.BigFromString(eventData.Data.Reward.String())

			// Change the struct retrieved from database to the correct type
			addedBlockHeader := customTypes.UniqueBlockHeader{
				Cid:       sql.NullString{String: eventData.Data.Cid},
				Height:    sql.NullString{String: eventData.Data.Height},
				Miner:     sql.NullString{String: eventData.Data.Miner},
				Timestamp: sql.NullString{String: eventData.Data.Timestamp},
				Validated: eventData.Data.Validated,
				Reward:    rewardAsBigInt,
				WinCount:  eventData.Data.WinCount,
			}

			// Launch notifiers for every block events (according to methods)
			// We need to wait 2 seconds before get messages because triggers is on headers
			// All messages are added are headers, so to be sure all messages available on database, wait 2 sec
			time.Sleep(2 * time.Second)
			messages := database.GetAllMessagesForTipSetAndCustomers(addedBlockHeader)

			if len(messages) > 0 {
				fmt.Println("A SAT Customer just mined a block, will send a Slack notification !")
				notifiers.SlackNotifier(messages)
			} else {
				fmt.Println("There is no SAT customer concerned by the last mined blocks !")
			}

		case <-time.After(90 * time.Second):
			go func() {
				err := l.Ping()

				if err != nil {
					fmt.Println("Error while ping database to check trigger")
					return
				}
			}()

		}
	}
}
