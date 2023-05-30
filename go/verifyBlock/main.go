package main

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	//"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/light"
	"github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/version"
)

const (
	maxClockDrift = 600 * time.Second * 10000000
)

type PartSetHeader struct {
	Total uint32 `json:"total"`
	Hash string `json:"hash"`
}

type Version struct {
	Block string `json:"block"`
	Parts PartSetHeader `json:"parts"`
}

type Block_ID struct {
	Hash string `json:"hash"`
	Parts PartSetHeader `json:"parts"`

}
type Header struct {
	Version Version `json:"version"`
	ChainID string `json:"chain_id"`
	Height string `json:"height"`
	Time string `json:"time"`
	LastBlockID Block_ID `json:"last_block_id"`
	LastCommitHash string `json:"last_commit_hash"`
	DataHash string `json:"data_hash"`
	ValidatorsHash string `json:"validators_hash"`
	NextValidatorsHash string `json:"next_validators_hash"`
	ConsensusHash string `json:"consensus_hash"`
	AppHash string `json:"app_hash"`
	LastResultsHash string `json:"last_results_hash"`
	EvidenceHash string `json:"evidence_hash"`
	ProposerAddress string `json:"proposer_address"`
}

type Data struct {
	Txs []string `json:"txs"`
}

type Signature struct {
	BlockIDFlag int `json:"block_id_flag"`
	ValidatorAddress string `json:"validator_address"`
	Timestamp string `json:"timestamp"`
	Signature string `json:"signature"`
}

type Commit struct {
	Height string `json:"height"`
	Round int32 `json:"round"`
	BlockID Block_ID `json:"block_id"`
	Signatures []Signature `json:"signatures"`
}

type PubKey struct {
	Type string `json:"type"`
	Value string `json:"value"`
}
type Validator struct {
	Address string `json:"address"`
	PubKey PubKey `json:"pub_key"`
	VotingPower string `json:"voting_power"`
	ProposerPriority string `json:"proposer_priority"`
}

type CurrentValidator struct {
	BlockHeight string `json:"block_height"`
	Validators []Validator `json:"validators"`
}

type SignedHeader struct {
	Header Header `json:"header"`
	Commit Commit `json:"commit"`
}
func getBytesFileJson(path string) []byte {
	// Open our jsonFile
    jsonFile, err := os.Open(path + ".json")
    // if we os.Open returns an error then handle it
    if err != nil {
        fmt.Println(err)
    }
    // fmt.Println("Successfully Opened " + path + ".json")
    // defer the closing of our jsonFile so that we can parse it later on
    defer jsonFile.Close()
    byteValue, _ := ioutil.ReadAll(jsonFile)
	return byteValue
}

func First[T, U any](val T, _ U) T {
	return val
}


func stringToInt64(str string) (int64) {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

func getBlockHeader(signed_header_path string) types.SignedHeader {
	var signed_header SignedHeader
	json.Unmarshal(getBytesFileJson(signed_header_path), &signed_header)
	var commit_sigs []types.CommitSig
	
	// commit_sigs_pointer := make([]*types.CommitSig, len(signed_header.Commit.Signatures))
    for i := 0; i < len(signed_header.Commit.Signatures); i++ {
		commit_sigs = append(commit_sigs, types.CommitSig {
			BlockIDFlag: types.BlockIDFlag(signed_header.Commit.Signatures[i].BlockIDFlag),
			ValidatorAddress: First(hex.DecodeString(signed_header.Commit.Signatures[i].ValidatorAddress)),
			Timestamp:        First(time.Parse(time.RFC3339, signed_header.Commit.Signatures[i].Timestamp)),
			Signature:        First(base64.StdEncoding.DecodeString(signed_header.Commit.Signatures[i].Signature)),
		})
		// fmt.Println("address:", hexutil.Encode(First(hex.DecodeString(signed_header.Commit.Signatures[i].ValidatorAddress))))
		// fmt.Println("signature",First(base64.StdEncoding.DecodeString(signed_header.Commit.Signatures[i].Signature)))
		// fmt.Println(i, First(time.Parse(time.RFC3339, signed_header.Commit.Signatures[i].Timestamp)));
		// commit_sigs_pointer[i] =  &commit_sigs[i]
	}
	fmt.Println("")
	var header = types.SignedHeader{
		Header: &types.Header{
			Version: version.Consensus{
				Block: uint64(stringToInt64(signed_header.Header.Version.Block)),
			},
			ChainID: signed_header.Header.ChainID,
			Height:  stringToInt64(signed_header.Header.Height),
			Time:    First(time.Parse(time.RFC3339, signed_header.Header.Time)),
			LastBlockID: types.BlockID{
				Hash: First(hex.DecodeString(signed_header.Header.LastBlockID.Hash)),
				PartSetHeader: types.PartSetHeader{
					Total: signed_header.Header.LastBlockID.Parts.Total,
					Hash:  First(hex.DecodeString(signed_header.Header.LastBlockID.Parts.Hash)),
				},
			},
			LastCommitHash:     First(hex.DecodeString(signed_header.Header.LastCommitHash)),
			DataHash:           First(hex.DecodeString(signed_header.Header.DataHash)),
			ValidatorsHash:     First(hex.DecodeString(signed_header.Header.ValidatorsHash)),
			NextValidatorsHash: First(hex.DecodeString(signed_header.Header.NextValidatorsHash)),
			ConsensusHash:      First(hex.DecodeString(signed_header.Header.ConsensusHash)),
			AppHash:            First(hex.DecodeString(signed_header.Header.AppHash)),
			LastResultsHash:    First(hex.DecodeString(signed_header.Header.LastResultsHash)),
			EvidenceHash:       First(hex.DecodeString(signed_header.Header.EvidenceHash)),
			ProposerAddress:    First(hex.DecodeString(signed_header.Header.ProposerAddress)),
		},
		Commit: &types.Commit{
			Height: stringToInt64(signed_header.Commit.Height),
			Round:  signed_header.Commit.Round,
			BlockID: types.BlockID{
				Hash: First(hex.DecodeString(signed_header.Commit.BlockID.Hash)),
				PartSetHeader: types.PartSetHeader{
					Total: uint32(signed_header.Commit.BlockID.Parts.Total),
					Hash:  First(hex.DecodeString(signed_header.Commit.BlockID.Parts.Hash)),
				},
			},
			Signatures: commit_sigs,
		},
	}
	// fmt.Println(First(time.Parse(time.RFC3339, signed_header.Header.Time)));
	return header
}


func getValidatorsSet(path string) types.ValidatorSet {
	var current_validators CurrentValidator
	json.Unmarshal(getBytesFileJson(path), &current_validators)

	var validators []types.Validator
	validators_pointer :=  make([]*types.Validator, len(current_validators.Validators))

	for i := 0; i < len(current_validators.Validators); i++ {
		validators = append(validators, types.Validator{
			Address:  First(hex.DecodeString(current_validators.Validators[i].Address)),
			PubKey:  ed25519.PubKey(First(base64.StdEncoding.DecodeString(current_validators.Validators[i].PubKey.Value))),
			// PubKey:           ed25519.PubKey{},
			VotingPower:      stringToInt64(current_validators.Validators[i].VotingPower),
			ProposerPriority: stringToInt64(current_validators.Validators[i].ProposerPriority),
		})
		validators_pointer[i] =  &validators[i]
		// fmt.Println(ed25519.PubKey(First(base64.StdEncoding.DecodeString(current_validators.Validators[i].PubKey.Key))))
		
		// fmt.Println("address val", First(hex.DecodeString(current_validators.Validators[i].Address)))
		// fmt.Println(`["val",`, First(base64.StdEncoding.DecodeString(current_validators.Validators[i].PubKey.Key)))
	}

	var vals = types.ValidatorSet{
		Validators: validators_pointer,
		// Proposer: &validators[0],
	}
	return vals

}

func getTransactions(transactionsDataPath string) types.Txs {
	var data Data
	json.Unmarshal(getBytesFileJson(transactionsDataPath), &data)
	var tx []types.Tx
	
    for i := 0; i < len(data.Txs); i++ {
		tx = append(tx, 
			First(base64.StdEncoding.DecodeString(data.Txs[i])),
		)
	}
	return tx;
}

func main() {
	// Open our jsonFile
	last_signed_header_path := "resources/last_signed_header"
	new_signed_header_path := "resources/new_signed_header"
	validator_path := "resources/validators"
	// transactionsDataPath := "resources/transaction_data"

	var (
		header = getBlockHeader(last_signed_header_path)

		newHeader = getBlockHeader(new_signed_header_path)
		// // 20, 30, 40, 50 - the first 3 don't have 2/3, the last 3 do!
		vals = getValidatorsSet(validator_path)
		// txs = getTransactions(transactionsDataPath)
	)

	// for i := 0; i < len(vals.Validators); i++ {
	// 	fmt.Print("pub key vals ", i , vals.Validators[i].PubKey, "\n")
	// }
	//fmt.Println(header)
	// fmt.Println(newHeader)

	err := light.VerifyAdjacent(&header, &newHeader, &vals, 1000*time.Hour, First(time.Parse(time.RFC3339, "2023-02-18T17:07:42Z")), maxClockDrift)

	fmt.Println(err)
	// fmt.Println(txs)
	// fmt.Printf("%x\n", txs.Hash())
	
	
}
