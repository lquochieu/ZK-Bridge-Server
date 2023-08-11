const { addLeaf, getTree, initialize, hash, getSiblings } = require("../../circom/txs/fmt");
const { readJSONFilesInFolder, getAddresFromAsciiString, saveJsonData } = require("../../circom/helper");
const { execSync } = require("child_process");
const fs = require("fs");
const DepositInfor = require("../models/DepositInfor");
const DepositTree = require("../models/DepositTree");

exports.generateUserProof = async (index) => {
    await initialize();
    const treeInfor = await DepositTree.findOne();
    if ( treeInfor.n_leafs <= index ) {
        return [];
    }
    const data = await DepositInfor.find( { "key": { "$lt": treeInfor.n_leafs } } ).sort({"key": 1});
    const tree = getTree();
    
    if (data.length < index) {
        throw("Invalid key");
    }
    const depositTree = await DepositTree.findOne();
    if (depositTree == null) {
        throw("Invalid tree!");
    }

    const newValue = Array.from(data, (data) => hash([data.eth_bridge_address, data.eth_receiver, data.amount, data.eth_token_address]));

    for (let i = 0; i < newValue.length; i++) {
        tree.update(i, newValue[i]);
    }

    const siblings = getSiblings(index);

    const input = {
        eth_bridge_address: data[index].eth_bridge_address,
        eth_receiver: data[index].eth_receiver,
        amount: data[index].amount,
        eth_token_address: data[index].eth_token_address,
        key: data[index].key,
        siblings: siblings.map(e => e.toString()),
        root: tree.root()
    };
    // console.log(input);
    json = JSON.stringify(input, (key, value) =>
        typeof value === 'bigint'
            ? value.toString()
            : value
        , 2);
    console.log(json)
    saveJsonData("./circom/circuit/verifyClaimTransaction/input.json", input)

    // generate user proof
    execSync(`cd circom/circuit/verifyClaimTransaction/verifyclaimtransaction_js/ && node generate_witness.js ./*.wasm ../input.json ../witness.wtns`);
    execSync(`cd circom/circuit/verifyClaimTransaction/verifyclaimtransaction_js/ && ../../../../rapidsnark/build/prover ../circuit_final.zkey ../witness.wtns ../../../../resources/verifyClaimTransaction/proof.json ../../../../resources/verifyClaimTransaction/public.json`)

    const public = JSON.parse(fs.readFileSync("./resources/verifyClaimTransaction/public.json", "utf-8"));
    const proof = JSON.parse(fs.readFileSync("./resources/verifyClaimTransaction/proof.json", "utf-8"));

    return {
        public: public,
        proof: proof
    }
};