package filecoin

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/lotus/api/v1api"
	"github.com/filecoin-project/lotus/chain/stmgr"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/ipfs/go-cid"
	"log"
	"net/http"
	"twinQuasarAppV2/backend/customTypes"
)

func Connect(config customTypes.Config) (v1api.FullNodeStruct, jsonrpc.ClientCloser) {
	var api v1api.FullNodeStruct
	var err error
	var closer jsonrpc.ClientCloser

	if config.Node.TokenRequired {
		closer, err = jsonrpc.NewMergeClient(
			context.Background(),
			config.Node.Endpoint,
			"Filecoin",
			[]interface{}{&api.Internal, &api.CommonStruct.Internal},
			http.Header{"Authorization": []string{config.Node.Token}})

		if err != nil {
			log.Fatalf("connecting with lotus failed: %s", err)
		}
	} else {
		closer, err = jsonrpc.NewMergeClient(
			context.Background(),
			config.Node.Endpoint,
			"Filecoin",
			[]interface{}{&api.Internal, &api.CommonStruct.Internal},
			http.Header{"test": []string{"application/json"}})

		if err != nil {
			log.Fatalf("connecting with lotus failed: %s", err)
		}
	}

	return api, closer
}

func Disconnect(closer jsonrpc.ClientCloser) {
	defer closer()
}

func VerifySignature(api v1api.FullNodeStruct, minerAddress string, originalMessage string, signedMessage string) bool {
	addr, err := address.NewFromString(minerAddress)

	if err != nil {
		fmt.Printf("Invalid address with error: %s\n", err)
		return false
	}

	message, err := hex.DecodeString(originalMessage)

	if err != nil {
		fmt.Printf("Invalid original message with error %s\n", err)
		return false
	}

	signature, err := hex.DecodeString(signedMessage)

	if err != nil {
		fmt.Printf("Invalid signature with error %s\n", err)
		return false
	}

	var sig crypto.Signature
	if err := sig.UnmarshalBinary(signature); err != nil {
		return false
	}

	result, err := api.WalletVerify(context.Background(), addr, message, &sig)

	if err != nil {
		fmt.Printf("calling wallet verify with error: %s\n", err)
		return false
	}

	return result
}

func GetOwnerKey(api v1api.FullNodeStruct, minerAddress string) (bool, string) {
	addr, err := address.NewFromString(minerAddress)

	if err != nil {
		fmt.Printf("Invalid address with error: %s\n", err)
		return true, err.Error()
	}

	minerInfo, errGetMinerInfo := api.StateMinerInfo(context.Background(), addr, types.EmptyTSK)

	if errGetMinerInfo != nil {
		fmt.Printf("Error when retrieving miner informations with message : %s\n", errGetMinerInfo)
		return true, errGetMinerInfo.Error()
	}

	result, err := api.StateAccountKey(context.Background(), minerInfo.Owner, types.EmptyTSK)

	if err != nil {
		fmt.Printf("Invalid address with error: %s\n", err)
		return true, err.Error()
	}

	return false, result.String()
}

// JsonParams Function to decode a message using specific parameters
func JsonParams(code cid.Cid, method abi.MethodNum, params []byte) (string, error) {
	p, err := stmgr.GetParamType(code, method)
	if err != nil {
		return "", err
	}

	if err := p.UnmarshalCBOR(bytes.NewReader(params)); err != nil {
		return "", err
	}

	b, err := json.MarshalIndent(p, "", "  ")
	return string(b), err
}
