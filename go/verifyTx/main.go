package main

import (
	//"github.com/tendermint/tendermint/crypto/tmhash"
	"encoding/json"
	//"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"encoding/base64"
	types "github.com/tendermint/tendermint/types"
)



func GetUnstructBytesFileJson(path string) map[string]interface{} {
	byteValue := GetBytesFileJson(path)
    var result map[string]interface{}
    json.Unmarshal([]byte(byteValue), &result)
	return result
}

func GetBytesFileJson(path string) []byte {
	// Open our jsonFile
	jsonFile, err := os.Open(path + ".json")
	// if we os.Open returns an error then handle it
	if err != nil {
		
	}
	
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	return byteValue
}

func TxHashToBytes(txHashs []string) types.Txs {
	var txs types.Txs
	for _, txHash := range txHashs {
		txBytes, err := base64.StdEncoding.DecodeString(txHash)
		// fmt.Printf("%x\n", txBytes)
		// fmt.Println()
		// fmt.Println("txByte", hex.EncodeToString(txBytes))
		if err != nil {
			panic(err)
		}
		txs = append(txs, txBytes)
		// fmt.Println(txs)
	}
	return txs
}

func First[T, U any](val T, _ U) T {
	return val
}

func main() {

	validatorsSet := GetValidatorsSet("../../resources/updateRootDepositToCosmosBridge/validators1")

	block_header_commit := GetBlockHeader("../../resources/updateRootDepositToCosmosBridge/block_header_commit")

	verifySignaturesInput := GenValidatorSignaturesInput(validatorsSet, block_header_commit)

	for i := 0; i < len(verifySignaturesInput); i++ {
		SaveVerifySignatureInputToJsonFile(verifySignaturesInput[i], "./verifySignatures/input" + strconv.Itoa(i) )
	}

	// txs := GetTx("../../resources/updateRootDepositToCosmosBridge/tx_data")
	// // fmt.Println("sender", tmhash.Sum([]byte(txs.Body.Messages[0].Sender))[:20])

	// input := GetDepositTreeInput("../../resources/updateRootDepositToCosmosBridge/block_header", txs)

	// SaveDepositRootToJsonFile(input, "../../resources/updateRootDepositToEthBridge/input_go")
	// fmt.Println("txs", input);
	// fmt.Printf("decode string  %x\n", tmhash.Sum(First(txs.Marshal())))
}