package main

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	//"fmt"
	"io/ioutil"
	//"os"
	"strconv"
	"bytes"
	//"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/tendermint/tendermint/crypto/ed25519"
	//"github.com/tendermint/tendermint/light"
	"github.com/tendermint/tendermint/types"
	//"github.com/tendermint/tendermint/crypto/merkle"
)

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


func SaveVerifySignatureInputToJsonFile(data VerifySignatures, path string)  {
	file, _ := json.MarshalIndent(data, "", " ")	
	err := ioutil.WriteFile(path + ".json", file, 0644)

	if err != nil {
		panic(err)
	}
}

func StringToInt64(str string) (int64) {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

func ElementEncodeTobytes(elements []string) [][]byte {
	var txs [][]byte
	for _, e := range elements {
		txBytes, _ := hex.DecodeString(e)
		// // fmt.Printf("%x\n", txBytes)
		// // fmt.Println()
		// // fmt.Println("txByte", hex.EncodeToString(txBytes))
		// if err != nil {
		// 	panic(err)
		// }
		txs = append(txs, txBytes)
		// fmt.Println(txs)
	}
	return txs
}


func GetValidatorsSet(path string) types.ValidatorSet {
	var current_validators CurrentValidator
	json.Unmarshal(GetBytesFileJson(path), &current_validators)

	var validators []types.Validator
	validators_pointer :=  make([]*types.Validator, len(current_validators.Validators))

	for i := 0; i < len(current_validators.Validators); i++ {
		validators = append(validators, types.Validator{
			Address:  First(hex.DecodeString(current_validators.Validators[i].Address)),
			PubKey:  ed25519.PubKey(First(base64.StdEncoding.DecodeString(current_validators.Validators[i].PubKey.Value))),
			// PubKey:           ed25519.PubKey{},
			VotingPower:      StringToInt64(current_validators.Validators[i].VotingPower),
			ProposerPriority: StringToInt64(current_validators.Validators[i].ProposerPriority),
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

type VerifySignatures struct {
	Height int64 `json:"height"`
	BlockHash []byte `json:"blockHash"`
	BlockTime int64 `json:"blockTime"`
	PartsTotal uint32 `json:"partsTotal"`
	PartsHash []byte `json:"partsHash"`
	SigTimeSeconds int `json:"sigTimeSeconds"`
	SigTimeNanos int `json:"sigTimeNanos"`
	PubKeys []byte `json:"pubKeys"`
	R8 []byte `json:"R8"`
	S []byte `json:"S"`
}

func FindPubkey(validatorsSet types.ValidatorSet, index int, address  types.Address) []byte {
	var pubkey = validatorsSet.Validators[index].PubKey.Bytes()
	if !bytes.Equal(address.Bytes(), validatorsSet.Validators[index].Address.Bytes()) {
		for i := 0; i < validatorsSet.Size(); i++ {
			if bytes.Equal(address.Bytes(), validatorsSet.Validators[i].Address.Bytes()) {
				pubkey = validatorsSet.Validators[i].PubKey.Bytes()
				return pubkey
			}
		}
	} 
	return pubkey
}

func FindIndexPubkeyInCommit(signedHeader types.SignedHeader, index int, address  types.Address) int {
	var addr = signedHeader.Commit.Signatures[index].ValidatorAddress
	if !bytes.Equal(address.Bytes(), addr.Bytes()) {
		for i := 0; i < len(signedHeader.Commit.Signatures); i++ {
			if bytes.Equal(address.Bytes(), signedHeader.Commit.Signatures[i].ValidatorAddress) {
				return i
			}
		}
	} else {
		return index
	}
	return -1
}

func GenValidatorSignaturesInput (validatorsSet types.ValidatorSet, signedHeader types.SignedHeader) []VerifySignatures {
	var verifySignatures []VerifySignatures
	for  i := 0; i < validatorsSet.Size(); i++ {
		verifySignatures = append(verifySignatures, VerifySignatures{
			Height: signedHeader.Height,
			BlockHash: signedHeader.Commit.BlockID.Hash.Bytes(),
			BlockTime: signedHeader.Time.Unix(),
			PartsTotal: signedHeader.LastBlockID.PartSetHeader.Total,
			PartsHash: signedHeader.LastBlockID.PartSetHeader.Hash.Bytes(),
			SigTimeSeconds: signedHeader.Commit.Signatures[i].Timestamp.Second(),
			SigTimeNanos: signedHeader.Commit.Signatures[i].Timestamp.Nanosecond(),
			PubKeys: FindPubkey(validatorsSet, i, signedHeader.Commit.Signatures[i].ValidatorAddress),
			R8: signedHeader.Commit.Signatures[i].Signature[0:32],
			S: signedHeader.Commit.Signatures[i].Signature[32:64],
		})
	}
	return verifySignatures
}

type ValidatosLeft struct {
	Pubkeys [][]byte `json:"pubkeys"`
	VotingPower []int64 `json:"votingPowers"`
	ChildRoot [][]byte `json:"childRoot"`
	ValidatorHash []byte `json:"chilRoot"`
	Signed int64 `json:"signed"`
}

type ValidatosRight struct {
	Pubkeys [][]byte `json:"pubkeys"`
	VotingPower []int64 `json:"votingPowers"`
	ChildRoot [][]byte `json:"childRoot"`
	ValidatorHash []byte `json:"chilRoot"`
	DataHash []byte `json:"dataHash"`
	ParrentSiblings [][]byte `json:"parrentSiblings"`
	BlockHash [][]byte `json:"blockHash"`
	Signed int64 `json:"signed"`
} 

type VerifyValidatorsSet struct {
	ValidatorsLeft ValidatosLeft
	ValidatosRight ValidatosRight
}
// func GenValidatorsInput (validatorsSet types.ValidatorSet, signedHeader types.SignedHeader) []VerifySignatures {
// 	var validatorEncode []string
// 	var pubKey [][]byte
// 	var votingPower []int64
// 	var isSigned []int
	
// 	// txBzs := make([][]byte, len(txs))
// 	// for i := 0; i < len(txs); i++ {
// 	// 	txBzs[i] = txs[i].Hash()
// 	// }
// 	for i := 0; i < validatorsSet.Size(); i++ {
// 		validatorEncode = append(validatorEncode, base64.StdEncoding.EncodeToString(validatorsSet.Validators[i].Bytes()))
// 		pubKey = append(pubKey, validatorsSet.Validators[i].PubKey.Bytes())
// 		votingPower = append(votingPower, validatorsSet.Validators[i].VotingPower)
// 		index := FindIndexPubkeyInCommit(signedHeader, i, validatorsSet.Validators[i].Address)
// 		if index >= 0  {
// 			isSigned = append(isSigned, 1)
// 		} else {
// 			isSigned = append(isSigned, 0)
// 		}
// 	}

// 	var validatorTxs = TxHashToBytes(validatorEncode)
// 	validatorBz := make([][]byte, len(validatorTxs))
// 	for i := 0; i < len(validatorTxs); i++ {
// 		validatorBz[i] = validatorTxs[i].Hash()
// 	}
// 	validatorHash := merkle.HashFromByteSlices(validatorBz)
	
// 	var header [][]byte

// 	version := signedHeader.Header.Version.ToProto()
// 	chainId :=  signedHeader.Header.ChainID
// 	header = append(header, First(version.Marshal()))

// 	var verifySignatures []VerifySignatures
// 	for  i := 0; i < validatorsSet.Size(); i++ {
// 		verifySignatures = append(verifySignatures, VerifySignatures{
// 			Height: signedHeader.Height,
// 			BlockHash: signedHeader.Commit.BlockID.Hash.Bytes(),
// 			BlockTime: signedHeader.Time.Unix(),
// 			PartsTotal: signedHeader.LastBlockID.PartSetHeader.Total,
// 			PartsHash: signedHeader.LastBlockID.PartSetHeader.Hash.Bytes(),
// 			SigTimeSeconds: signedHeader.Commit.Signatures[i].Timestamp.Second(),
// 			SigTimeNanos: signedHeader.Commit.Signatures[i].Timestamp.Nanosecond(),
// 			PubKeys: FindPubkey(validatorsSet, i, signedHeader.Commit.Signatures[i].ValidatorAddress),
// 			R8: signedHeader.Commit.Signatures[i].Signature[0:32],
// 			S: signedHeader.Commit.Signatures[i].Signature[32:64],
// 		})
// 	}
// }
// func main() {
// 	// Open our jsonFile
// 	last_signed_header_path := "resources/last_signed_header"
// 	new_signed_header_path := "resources/new_signed_header"
// 	validator_path := "resources/validators"
// 	// transactionsDataPath := "resources/transaction_data"

// 	var (
// 		header = GetBlockHeader(last_signed_header_path)

// 		newHeader = GetBlockHeader(new_signed_header_path)
// 		// // 20, 30, 40, 50 - the first 3 don't have 2/3, the last 3 do!
// 		vals = GetValidatorsSet(validator_path)
// 		// txs = getTransactions(transactionsDataPath)
// 	)

// 	// for i := 0; i < len(vals.Validators); i++ {
// 	// 	fmt.Print("pub key vals ", i , vals.Validators[i].PubKey, "\n")
// 	// }
// 	//fmt.Println(header)
// 	// fmt.Println(newHeader)

// 	err := light.VerifyAdjacent(&header, &newHeader, &vals, 1000*time.Hour, First(time.Parse(time.RFC3339, "2023-02-18T17:07:42Z")), maxClockDrift)

// 	fmt.Println(err)
// 	// fmt.Println(txs)
// 	// fmt.Printf("%x\n", txs.Hash())
	
	
// }
