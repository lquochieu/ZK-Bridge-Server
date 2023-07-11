//@ts-nocheck
import { SigningCosmWasmClient, Secp256k1HdWallet } from "cosmwasm";
import * as dotenv from "dotenv";
import { Decimal } from "@cosmjs/math";
import axios from 'axios';
import * as fs from "fs";
// import MerkleTree from "fixed-merkle-tree";
// import { buildMimc7 } from "circomlibjs";

dotenv.config();
// This is your rpc endpoint
const getTxAPI = "https://testnet-lcd.orai.io/cosmos/"

const rpcEndpoint = "https://testnet-rpc.orai.io:443/";
const mnemonic = process.env.MNEMONIC_COSMOS!;


function saveJsonData(path: string, data: any) {
    const jsonData = JSON.stringify(data, null, 2);
    fs.writeFileSync(path, jsonData, 'utf-8');
    console.log('Data has been saved to file:', path);
}

async function getWallet(): Promise<Secp256k1HdWallet> {
    const wallet = await Secp256k1HdWallet.fromMnemonic(mnemonic, { prefix: "orai" });
    return wallet;
}
async function getClient(): Promise<SigningCosmWasmClient> {
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



async function QueryDepositTree() {
    // const query = await client.getTx("2D925C0F81EF1E26662B0A2A9277180CE853F9F07C60CA2F3E64E7F565A19F78")
    const client = await getClient();

    const contract_address = process.env.COSMOS_BRIDGE || "";
    const query_message = {
        deposit_tree: {}
    }
    const res = await client.queryContractSmart(contract_address, query_message)
    return res
}
async function QueryDepositQueue() {
    const client = await getClient();

    const contract_address = process.env.COSMOS_BRIDGE || "";
    const query_message = {
        deposit_queue: {}
    }
    console.log("deposit_queue mssage", query_message);
    const res = await client.queryContractSmart(contract_address, query_message)
    return res
}



async function main() {

    const resDepositTree = await QueryDepositTree();
    console.log("depositTree", resDepositTree);
    saveJsonData("./circom/txs/depositTree/deposit_tree.json", resDepositTree);
    const resDepositQueue = await QueryDepositQueue();
    console.log(resDepositQueue);

    for(let i = 0; i < resDepositQueue.length; i++) {
        saveJsonData(`./circom/txs/depositInfo/deposit${i + resDepositTree.n_leafs}.json`, resDepositQueue[i])
    }

}
main();