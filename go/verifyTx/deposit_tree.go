package main

import (
	// "bytes"
	"encoding/base64"
	"fmt"
	//"fmt"
	//"encoding/hex"
	"encoding/json"
	// "github.com/tendermint/tendermint/crypto/tmhash"
	"bytes"

	"github.com/tendermint/tendermint/crypto/merkle"

	//types "github.com/tendermint/tendermint/types"
	"io/ioutil"
)

func SaveDepositRootToJsonFile(data DepositRootCosmos, path string)  {
	file, _ := json.MarshalIndent(data, "", " ")	
	err := ioutil.WriteFile(path + ".json", file, 0644)
	fmt.Println("Save file " + path + ".json successed")
	if err != nil {
		panic(err)
	}
}

type BlockDepositRootCosmos struct {
	Header Header `json:"header"`
	Txs []string `json:"txs"`
}

type DepositRootCosmos struct {
	TxBody []byte `protobuf:"byte,1,opt,name=tx_body,proto3" json:"txBody"`
	TxAuthInfos []byte `protobuf:"byte,1,opt,name=tx_ath_info,proto3" json:"txAuthInfos"`
	Signatures []byte `protobuf:"byte,1,opt,name=signatures,proto3" json:"signatures"`
	Key int `protobuf:"int,1,opt,name=key,proto3" json:"key"`
	DataHash []byte `protobuf:"byte,1,opt,name=data_hash,proto3" json:"dataHash"`
	Siblings [][]byte `protobuf:"byte,1,opt,name=siblings_root,proto3" json:"siblings"`
	// Sender string `protobuf:"byte,1,opt,name=sender,proto3" json:"sender"`
	// Contract string `protobuf:"byte,1,opt,name=contract,proto3" json:"contract"`
	// DepositRoot string `protobuf:"byte,1,opt,name=deposit_root,proto3" json:"deposit_root"`
}


func GetDepositTreeInput(path string, transaction Transaction) DepositRootCosmos {
	var blockDepositRootCosmos BlockDepositRootCosmos
	json.Unmarshal(GetBytesFileJson(path), &blockDepositRootCosmos)

	// dataHash := First(base64.StdEncoding.DecodeString(blockDepositRootCosmos.Header.DataHash))

	var key int
	_ = key

	key = -1
	for i := 0; i < len(blockDepositRootCosmos.Txs); i++ {	
		var tx = First(base64.StdEncoding.DecodeString(blockDepositRootCosmos.Txs[i]))
		if(bytes.Equal(First(transaction.Marshal()), tx)) {
			key = i;
		}
	}


	txs := TxHashToBytes(blockDepositRootCosmos.Txs)
	fmt.Println("Txs ???", txs)
	txBzs := make([][]byte, len(txs))
	for i := 0; i < len(txs); i++ {
		txBzs[i] = txs[i].Hash()
	}

	root := merkle.HashFromByteSlices(txBzs)
	
	proof := txs.Proof(key)
	siblings := proof.Proof.Aunts

	if(len(siblings) < 2) {
		for i := len(siblings); i < 2; i++ {
			siblings = append(siblings, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
		}
	}
	dataHash := root

	txBody := First(transaction.Body.Marshal())
	if len(txBody) < 922 {
		for i := len(txBody); i < 922; i++ {
			txBody = append(txBody, 0)
		}
	}
	
	txAuthInfos := First(transaction.AuthInfo.Marshal())
	if len(txAuthInfos) < 109 {
		for i := len(txAuthInfos); i < 109; i++ {
			txAuthInfos = append(txAuthInfos, 0)
		}
	}

	
	var depositRootCosmos = DepositRootCosmos{
		TxBody: txBody,
		TxAuthInfos: txAuthInfos,
		Signatures: First(SignatureMarshal(transaction)),
		Key: key,
		DataHash: dataHash,
		Siblings: siblings,
	}
	return depositRootCosmos
}
