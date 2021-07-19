package customTypes

import (
	"database/sql"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/big"
)

type UniqueBlockKey struct {
	Nonce 			uint64  			`json:"nonce"`
	From 			address.Address  	`json:"from"`
}

type UniqueBlockMessage struct {
	Cid				sql.NullString			`json:"messageCid"`
	To				sql.NullString			`json:"to"`
	From			sql.NullString			`json:"from"`
	Nonce			uint64					`json:"nonce"`
	MethodName		sql.NullString			`json:"method"`
	Actor			sql.NullString			`json:"actor"`
	DecodedParams 	sql.NullString			`json:"params"`
	GasFeeCap		int64					`json:"gasFeeCap"`
	GasLimit		int64					`json:"gasLimit"`
	GasPremium		int64					`json:"gasPremium"`
	Value			int64					`json:"value"`
	BlockCid		sql.NullString			`json:"blockCid"`
	Status			bool					`json:"status"`
	Duplicate		bool					`json:"duplicate"`
}

type UniqueBlockHeader struct {
	Height 				sql.NullString 		`json:"height"`
	Timestamp			sql.NullString		`json:"timestamp"`
	Cid					sql.NullString		`json:"cid"`
	Miner 				sql.NullString		`json:"miner"`
	Validated 			bool				`json:"validated"`
	Reward				big.Int				`json:"reward"`
	WinCount			int64				`json:"winCount"`
}

type UniqueTipSetContent struct {
	Header 		[]UniqueBlockHeader  	`json:"headers"`
	Messages 	[]UniqueBlockMessage 	`json:"messages"`
}

type UniqueBlockSingle struct {
	Header 		UniqueBlockHeader  	`json:"header"`
	Message 	UniqueBlockMessage 	`json:"message"`
	User		User				`json:"user"`
}