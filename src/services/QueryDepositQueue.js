//@ts-nocheck
const { SigningCosmWasmClient, Secp256k1HdWallet } = require("cosmwasm");
require("dotenv").config();

const { Decimal } = require("@cosmjs/math");

const DepositTree = require("../models/DepositTree");
const DepositInfor = require("../models/DepositInfor");
// import MerkleTree from "fixed-merkle-tree";
// import { buildMimc7 } from "circomlibjs";

// This is your rpc endpoint
const getTxAPI = "https://testnet-lcd.orai.io/cosmos/"

const rpcEndpoint = "https://testnet-rpc.orai.io:443/";

const mnemonic = process.env.MNEMONIC_COSMOS ?? "";

function saveJsonData(path, data) {
    const jsonData = JSON.stringify(data, null, 2);
    fs.writeFileSync(path, jsonData, 'utf-8');
    console.log('Data has been saved to file:', path);
}

const getWallet = async() => {
    const wallet = await Secp256k1HdWallet.fromMnemonic(mnemonic, { prefix: "orai" });
    return wallet;
}
const getClient = async() => {
    // Create a wallet
    const wallet = await getWallet();

    // Using
    const client = await SigningCosmWasmClient.connectWithSigner(
        rpcEndpoint,
        wallet,
        {
            gasPrice: {
                denom: "orai",
                amount: Decimal.fromUserInput("0.001", 6)
            }
        }
    );

    return client;
}

const QueryDepositTree = async() => {

    const client = await getClient();

    const contract_address = process.env.COSMOS_BRIDGE || "";
    const query_message = {
        deposit_tree: {}
    }
    const res = await client.queryContractSmart(contract_address, query_message)
    return {
        root: res.root ?? "",
        n_leafs: res.n_leafs ?? 0,
        nqueue_leafs: res.nqueue_leafs ?? 0
    }
}
const QueryDepositQueue = async() => {
    const client = await getClient();

    const contract_address = process.env.COSMOS_BRIDGE || "";
    const query_message = {
        deposit_queue: {}
    }
    console.log("deposit_queue mssage", query_message);
    const res = await client.queryContractSmart(contract_address, query_message)
    return res
}

exports.queryDepositQueue = async() => {
    const resDepositTree = await QueryDepositTree();
    console.log("depositTree", resDepositTree);
    
    const depositTree = await DepositTree.findOne();
    if (!depositTree) {
        const newDepositTree = new DepositTree({
            root: resDepositTree.root,
            n_leafs: resDepositTree.n_leafs,
            nqueue_leafs: resDepositTree.nqueue_leafs
        });

        await newDepositTree.save();
    } else {
        depositTree.root = resDepositTree.root;
        depositTree.n_leafs = resDepositTree.n_leafs;
        depositTree.nqueue_leafs = resDepositTree.nqueue_leafs;
        
        await depositTree.save();
    }

    const resDepositQueue = await QueryDepositQueue();
    console.log(resDepositQueue);

    for(let i = 0; i < resDepositQueue.length; i++) {
        const key = resDepositQueue[i].key;
        const depositInfor = await DepositInfor.findOne({key: key});
        if (depositInfor != null) {
            depositInfor.is_deposit = resDepositQueue[i].is_deposit;
            await depositInfor.save();
            continue;
        }

        const depositQueue = new DepositInfor({
            is_deposit: resDepositQueue[i].is_deposit,
            sender: resDepositQueue[i].sender,
            destination_chainid: resDepositQueue[i].destination_chainid,
            eth_bridge_address: resDepositQueue[i].eth_bridge_address,
            eth_receiver: resDepositQueue[i].eth_receiver,
            amount: resDepositQueue[i].amount,
            eth_token_address: resDepositQueue[i].eth_token_address,
            cosmos_token_address: resDepositQueue[i].cosmos_token_address,
            key: resDepositQueue[i].key,
            value: resDepositQueue[i].value
        });

        await depositQueue.save();
    }

    return {
        depositTree: depositTree,
        depositQueue: resDepositQueue
    }
}