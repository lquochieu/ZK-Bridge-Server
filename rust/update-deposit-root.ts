//@ts-nocheck
import { SigningCosmWasmClient, Secp256k1HdWallet  } from "cosmwasm";
import * as dotenv from "dotenv";
import { Decimal } from "@cosmjs/math";
import axios from 'axios';
import * as fs from "fs";
// This is your rpc endpoint
const getTxAPI = "https://testnet-lcd.orai.io/cosmos/"

const rpcEndpoint = "https://testnet-rpc.orai.io:443/";
dotenv.config();

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


function ReadJsonFile(path: string): Record<string, any> {

    const jsonData = fs.readFileSync(path, 'utf-8');
    const parsedData: Record<string, any> = JSON.parse(jsonData);
    return parsedData;
}

function getInputUpdateDepositTree() {
    let proofFile: Record<string, any> = ReadJsonFile("./resources/updateRootDepositToCosmosBridge/proof.json");
    let publicFile: Record<string, any> = ReadJsonFile("./resources/updateRootDepositToCosmosBridge/public.json");
    let proof: any[] = [];
    for (let i = 0; i < 2; i++) {
        proof.push(proofFile.pi_a[i]);
    }
    for (let i = 0; i < 2; i++) {
        for (let j = 0; j < 2; j++) {
            proof.push(proofFile.pi_b[i][j]);
        }
    }
    for (let i = 0; i < 2; i++) {
        proof.push(proofFile.pi_c[i]);
    }

    let msg = {
        update_deposit_tree: {
            root: publicFile[publicFile.length - 1],
            proof: proof,
        }
    }
    return msg
}

function saveJsonData(path: string, data: any) {
    const jsonData = JSON.stringify(data, null, 2);
    fs.writeFileSync(path, jsonData, 'utf-8');
    console.log('Data has been saved to file:', path);
}

async function QueryTxByHash(txHash: string): Promise<any> {
    // console.log(resInitiate)
    // console.log(wasmCode)
    let res = await axios.get(getTxAPI + "tx/v1beta1/txs/" + txHash).then(function (response: any) {
        // handle success
        return { tx: response.data.tx, height: response.data.tx_response.height }
        // console.dir(response.data.tx, { depth: null });
    })
        .catch(function (error) {
            // handle error
            console.log(error);
        })
        .finally(function () {
            // always executed
        });

    return res
}

async function QueryBlockHeaderByHeight(height: string): Promise<any> {
    let res = await axios.get(getTxAPI + "tx/v1beta1/txs/block/" + height).then(function (response: any) {
        // handle success
        return { header: response.data.block.header, txs: response.data.block.data.txs }
        // console.dir(response.data.tx, { depth: null });
    })
        .catch(function (error) {
            // handle error
            console.log(error);
        })
        .finally(function () {
            // always executed
        });

    return res
}

async function QueryValidatorsSetByHeight(height: string): Promise<any> {
    let res = await axios.get("https://testnet-rpc.orai.io/validators?height=" + height).then(function (response: any) {
        // handle success
        return response.data.result
        // console.dir(response.data.tx, { depth: null });
    })
        .catch(function (error) {
            // handle error
            console.log(error);
        })
        .finally(function () {
            // always executed
        });

    return res
}

async function updateDepositTree() {
    // const query = await client.getTx("2D925C0F81EF1E26662B0A2A9277180CE853F9F07C60CA2F3E64E7F565A19F78")
    const wallet = await getWallet();
    const client = await getClient();

    const senderAddress = (await wallet.getAccounts())[0].address;
    const contractAddress = process.env.COSMOS_BRIDGE || "";  
    const msg = getInputUpdateDepositTree();
    console.log("msg", msg);
    const fee = "auto"
    const memo: any = null;

    const res = await client.execute(senderAddress, contractAddress, msg, fee, memo)
    console.log(res)
    return res;
}

async function saveBlockHeader(height: string) {
    const response = await axios.request({
        method: "get",
        url: `https://testnet-rpc.orai.io/commit?heigh=${height}`
    });

    console.log(response.data);
    saveJsonData("./resources/cosmosHeader/cosmosHeader.json", response.data);
}

async function main() {
    const resUpdate = await updateDepositTree();
    console.log(resUpdate);
    writeToEnvFile("TX_HASH", resUpdate.transactionHash)

    const resQueryDepositRootTx = await QueryTxByHash(resUpdate.transactionHash);
    console.log(resQueryDepositRootTx);
    saveJsonData("./resources/updateRootDepositToCosmosBridge/tx_data.json", resQueryDepositRootTx);
    const resQueryBlock = await QueryBlockHeaderByHeight(resQueryDepositRootTx.height);
    saveJsonData("./resources/updateRootDepositToCosmosBridge/block.json" , resQueryBlock);
    console.log(resQueryBlock)
    const resQueryValidatorsSet = await QueryValidatorsSetByHeight(resQueryDepositRootTx.height);
    console.log(resQueryValidatorsSet)
    saveJsonData("./resources/updateRootDepositToCosmosBridge/validators.json", resQueryValidatorsSet);

    await saveBlockHeader(resQueryDepositRootTx.height);
}

main();