const fs = require("fs");
const { saveJsonData, base64ToBytes } = require("./helper");
const input_go = require("../../resources/updateRootDepositToEthBridge/input_go.json")

const main = async () => {
    let input = {
        txBody: base64ToBytes(input_go.txBody),
        txAuthInfos: base64ToBytes(input_go.txAuthInfos),
        signatures: base64ToBytes(input_go.signatures),
        key: input_go.key,
        dataHash: base64ToBytes(input_go.dataHash),
        siblings: input_go.siblings.map(e => base64ToBytes(e))
    }
    console.log(input)
    saveJsonData("./circom/circuit/verifyDepositRoot/input.json", input)
};

main()
    .then(() => { })
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });