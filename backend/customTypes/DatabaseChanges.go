package customTypes

import "encoding/json"

type AddedBlockEvent struct {
	Table  string         `json:"table"`
	Action string         `json:"action"`
	Data   AddedBlockData `json:"data"`
}

type AddedBlockData struct {
	Cid       string      `json:"cid"`
	Height    string      `json:"height"`
	Timestamp string      `json:"timestamp"`
	Miner     string      `json:"miner"`
	Validated bool        `json:"validated"`
	Reward    json.Number `json:"reward"`
	WinCount  int64       `json:"win_count"`
}
