const fs = require("fs");
const { addLeaf, getTree, initialize, hash, getSiblings } = require("../../circom/txs/fmt");
const { readJSONFilesInFolder, getAddresFromAsciiString, saveJsonData } = require("../../circom/helper");
const DepositInfor = require("../models/DepositInfor");
const DepositTree = require("../models/DepositTree");
const { execSync } = require("child_process");

const range = (start, stop, step) =>
    Array.from({ length: (stop - start) / step }, (_, i) => start + i * step);

exports.generateProofUpdateRoot = async () => {
    await initialize();
    const tree = getTree();

    const data = await DepositInfor.find().sort({"key": 1});
    const depositTree = await DepositTree.findOne();
    if (depositTree == null) {
        throw("Invalid tree!");
    }


    let n_leafs = depositTree.n_leafs;
    let nqueue_leafs = depositTree.nqueue_leafs
    let nmax_leafs_update = 5;

    const newValue = Array.from(data, (data) => hash([data.eth_bridge_address, data.eth_receiver, data.amount, data.eth_token_address]));
    console.log("----new---\n", newValue, "\n-----")
    for (i = nqueue_leafs; i < nmax_leafs_update; i++) {
        newValue.push(hash([0]));
    }

    for (let i = 0; i < n_leafs; i++) {
        tree.update(i, newValue[i]);
    }

    if (n_leafs == 0) {
        tree.update(0, hash([0]));
    }

    const oldRoot = tree.root();

    const siblings = [];
    console.log(newValue);
    for (let i = n_leafs; i < n_leafs + nmax_leafs_update; i++) {
        tree.update(i, newValue[i]);
        siblings.push(getSiblings(i));
    }
    const newRoot = tree.root();

    const input = {
        key: range(n_leafs, n_leafs + nmax_leafs_update, 1),
        newValue: newValue.map(e => e.toString()).slice(n_leafs, n_leafs + nmax_leafs_update),
        oldRoot: oldRoot,
        siblings: siblings.map(sib => sib.map(e => e.toString())),
        newRoot: newRoot,
    };

    // generate proof
    const path = "./circom/circuit/verifyRootBatchTxsCosmos";
    saveJsonData(`${path}/input.json`, input)
    execSync(`node ${path}/verifyrootbatchtxscosmos_js/generate_witness.js ${path}/verifyrootbatchtxscosmos_js/*.wasm ${path}/input.json ${path}/witness.wtns`);
    execSync(`./rapidsnark/build/prover ${path}/circuit_final.zkey ${path}/witness.wtns ./resources/updateRootDepositToCosmosBridge/proof.json ./resources/updateRootDepositToCosmosBridge/public.json`);
};