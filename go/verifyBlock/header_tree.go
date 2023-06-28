package main

import (
	// "bytes"
	//"encoding/base64"
	//"fmt"
	//"encoding/hex"
	//"encoding/json"
	// "github.com/tendermint/tendermint/crypto/tmhash"
	//"bytes"

	"github.com/tendermint/tendermint/crypto/merkle"

	types "github.com/tendermint/tendermint/types"
	//"io/ioutil"
	gogotypes "github.com/gogo/protobuf/types"
	tdmbytes "github.com/tendermint/tendermint/libs/bytes"
	"reflect"
)

// func SaveDepositRootToJsonFile(data DepositRootCosmos, path string)  {
// 	file, _ := json.MarshalIndent(data, "", " ")	
// 	err := ioutil.WriteFile(path + ".json", file, 0644)
// 	fmt.Println("Save file " + path + ".json successed")
// 	if err != nil {
// 		panic(err)
// 	}
// }

// type BlockDepositRootCosmos struct {
// 	Header Header `json:"header"`
// 	Txs []string `json:"txs"`
// }

// type DepositRootCosmos struct {
// 	TxBody []byte `protobuf:"byte,1,opt,name=tx_body,proto3" json:"txBody"`
// 	TxAuthInfos []byte `protobuf:"byte,1,opt,name=tx_ath_info,proto3" json:"txAuthInfos"`
// 	Signatures []byte `protobuf:"byte,1,opt,name=signatures,proto3" json:"signatures"`
// 	Key int `protobuf:"int,1,opt,name=key,proto3" json:"key"`
// 	DataHash []byte `protobuf:"byte,1,opt,name=data_hash,proto3" json:"dataHash"`
// 	Siblings [][]byte `protobuf:"byte,1,opt,name=siblings_root,proto3" json:"siblings"`
// 	// Sender string `protobuf:"byte,1,opt,name=sender,proto3" json:"sender"`
// 	// Contract string `protobuf:"byte,1,opt,name=contract,proto3" json:"contract"`
// 	// DepositRoot string `protobuf:"byte,1,opt,name=deposit_root,proto3" json:"deposit_root"`
// }

func isTypedNil(o interface{}) bool {
	rv := reflect.ValueOf(o)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Slice:
		return rv.IsNil()
	default:
		return false
	}
}

// Returns true if it has zero length.
func isEmpty(o interface{}) bool {
	rv := reflect.ValueOf(o)
	switch rv.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return rv.Len() == 0
	default:
		return false
	}
}

func cdcEncode(item interface{}) []byte {
	if item != nil && !isTypedNil(item) && !isEmpty(item) {
		switch item := item.(type) {
		case string:
			i := gogotypes.StringValue{
				Value: item,
			}
			bz, err := i.Marshal()
			if err != nil {
				return nil
			}
			return bz
		case int64:
			i := gogotypes.Int64Value{
				Value: item,
			}
			bz, err := i.Marshal()
			if err != nil {
				return nil
			}
			return bz
		case tdmbytes.HexBytes:
			i := gogotypes.BytesValue{
				Value: item,
			}
			bz, err := i.Marshal()
			if err != nil {
				return nil
			}
			return bz
		default:
			return nil
		}
	}

	return nil
}

func TxHashToBytes(bytesArray [][]byte) types.Txs {
	var txs types.Txs
	for _, element := range bytesArray {
		txs = append(txs, element)
		// fmt.Println(txs)
	}
	return txs
}


func GetDataAndValHashSiblings(header types.SignedHeader) [][]byte {
	var header_field [][]byte
	
	hpb := header.Version.ToProto()
	hbz := First(hpb.Marshal())
	header_field = append(header_field, hbz)

	chainIDMarshal := cdcEncode(header.ChainID)
	header_field = append(header_field, chainIDMarshal)

	pbt := First(gogotypes.StdTimeMarshal(header.Time))
	header_field = append(header_field, pbt)

	pbbi := header.LastBlockID.ToProto()
	bzbi := First(pbbi.Marshal())
	header_field = append(header_field, bzbi)

	HeightMarshal := cdcEncode(header.Height)
	header_field = append(header_field, HeightMarshal)

	LastCommitHashMarshal := cdcEncode(header.LastCommitHash)
	header_field = append(header_field, LastCommitHashMarshal)

	DataHashMarshal := cdcEncode(header.DataHash)
	header_field = append(header_field, DataHashMarshal)

	ValidatorsHashMarshal := cdcEncode(header.ValidatorsHash)
	header_field = append(header_field, ValidatorsHashMarshal)

	NextValidatorsHashMarshal := cdcEncode(header.NextValidatorsHash)
	header_field = append(header_field, NextValidatorsHashMarshal)

	ConsensusHashMarshal := cdcEncode(header.ConsensusHash)
	header_field = append(header_field, ConsensusHashMarshal)

	AppHashMarshal := cdcEncode(header.AppHash)
	header_field = append(header_field, AppHashMarshal)

	LastResultsHashMarshal := cdcEncode(header.LastResultsHash)
	header_field = append(header_field, LastResultsHashMarshal)

	EvidenceHashMarshal := cdcEncode(header.EvidenceHash)
	header_field = append(header_field, EvidenceHashMarshal)

	ProposerAddressMarshal := cdcEncode(header.ProposerAddress)
	header_field = append(header_field, ProposerAddressMarshal)
	
	txs := TxHashToBytes(header_field)
	txBzs := make([][]byte, len(txs))
	for i := 0; i < len(txs); i++ {
		txBzs[i] = txs[i].Hash()
	}

	merkle.HashFromByteSlices(txBzs)
	
	proof := txs.Proof(6)
	siblings := proof.Proof.Aunts[1:]

	return siblings
}
