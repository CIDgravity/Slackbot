package customTypes

import (
	"encoding/json"
	"github.com/filecoin-project/go-state-types/big"
)

type Rewards struct {
	BaselineTotal 				big.Int			`json:"BaselineTotal"`
	CumsumBaseline 				big.Int			`json:"CumsumBaseline"`
	CumsumRealized 				big.Int			`json:"CumsumRealized"`
	EffectiveBaselinePower 		big.Int			`json:"EffectiveBaselinePower"`
	EffectiveNetworkTime 		json.Number		`json:"EffectiveNetworkTime"`
	Epoch 						json.Number		`json:"Epoch"`
	SimpleTotal 				big.Int			`json:"SimpleTotal"`
	ThisEpochBaselinePower 		big.Int			`json:"ThisEpochBaselinePower"`
	ThisEpochReward 			big.Int			`json:"ThisEpochReward"`
}
