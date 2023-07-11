//@ts-nocheck
import { SigningCosmWasmClient, Secp256k1HdWallet, setupWebKeplr, coin, UploadResult, InstantiateResult, toBinary } from "cosmwasm";
import { CosmWasmClient } from "cosmwasm";
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

function writeToEnvFile(key: String, value: String) {
    const envFilePath = '.env';
    const envString = `${key}=${value}`;

    try {
        if (fs.existsSync(envFilePath)) {
            let data = fs.readFileSync(envFilePath, 'utf8');
            const lines = data.trim().split('\n');
            let keyExists = false;
            const updatedLines = lines.map(line => {
                const [existingKey] = line.split('=');
                if (existingKey === key) {
                    keyExists = true;
                    return envString;
                }
                return line;
            });
            if (!keyExists) {
                updatedLines.push(envString);
            }
            const updatedData = updatedLines.join('\n');
            fs.writeFileSync(envFilePath, updatedData + '\n');
        } else {
            fs.writeFileSync(envFilePath, envString + '\n');
        }
        console.log('Successfully wrote to .env file.');
    } catch (err) {
        console.error('Error writing to .env file:', err);
    }
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

async function Upload(path): Promise<UploadResult> {
    // const query = await client.getTx("2D925C0F81EF1E26662B0A2A9277180CE853F9F07C60CA2F3E64E7F565A19F78")
    const wallet = await getWallet();
    const client = await getClient();

    const senderAddress = (await wallet.getAccounts())[0].address;
    const wasmCode = ReadFile(path)
    const fee = "auto"
    const memo: any = null
    // const fund = [coin(2, "orai")]
    // const res = await client.execute(senderAddress, contractAddress, msg, fee, memo, fund)
    const res = await client.upload(senderAddress, wasmCode, fee, memo)
    console.log(res)
    return res;
}

async function instantiate(codeID: number): Promise<InstantiateResult> {
    // const query = await client.getTx("2D925C0F81EF1E26662B0A2A9277180CE853F9F07C60CA2F3E64E7F565A19F78")
    const wallet = await getWallet();
    const client = await getClient();

    const senderAddress = (await wallet.getAccounts())[0].address;

    const msg = {
        token_address: process.env.COSMOS_TOKEN,
        root: "19650381120764303012524736994073988891758455017340428694151271035030351959738"

    }
    const label = "Test ew20"

    const fee = "auto"
    // const option = {
    //     fund: {
    //         denom: "orai",
    //         amount: Decimal.fromUserInput("0.001", 6)
    //     },
    //     admin: senderAddress
    // }
    const res = await client.instantiate(senderAddress, codeID, msg, label, fee)
    console.log(res)
    return res;
}




async function main() {
    const resUpload = await Upload("./artifacts/oraisan_cosmos_contract_demo.wasm");
    const resInitiate = await instantiate(resUpload.codeId);
    writeToEnvFile("COSMOS_BRIDGE", resInitiate.contractAddress)
}
main();