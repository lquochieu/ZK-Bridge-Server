package main

import (
	//"github.com/tendermint/tendermint/crypto/tmhash"
	"os"
	"io/ioutil"
	"encoding/json"
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

func First[T, U any](val T, _ U) T {
	return val
}

func main() {
	txs := GetTx("../../resources/updateRootDepositToCosmosBridge/tx_data")
	// fmt.Println("sender", tmhash.Sum([]byte(txs.Body.Messages[0].Sender))[:20])

	input := GetDepositTreeInput("../../resources/updateRootDepositToCosmosBridge/block_header", txs)

	SaveBlockDepositRoot(input, "../../resources/updateRootDepositToEthBridge/input_go")
	// fmt.Println("txs", input);
	// fmt.Printf("decode string  %x\n", tmhash.Sum(First(txs.Marshal())))
}