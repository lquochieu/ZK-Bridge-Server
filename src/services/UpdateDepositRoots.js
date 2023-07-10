const { SigningCosmWasmClient, Secp256k1HdWallet  } = require("cosmwasm")
const dotenv = require("dotenv")
const { Decimal } = require("@cosmjs/math");
const axios = require('axios');
const fs = require("fs");
const { execSync } = require("child_process");
// This is your rpc endpoint
const getTxAPI = "https://testnet-lcd.orai.io/cosmos/"

const rpcEndpoint = "https://testnet-rpc.orai.io:443/";
dotenv.config();

const mnemonic = process.env.MNEMONIC_COSMOS ?? "";

const writeToEnvFile = (key, value) => {
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


const ReadJsonFile = (path) => {

    const jsonData = fs.readFileSync(path, 'utf-8');
    const parsedData = JSON.parse(jsonData);
    return parsedData;
}

const getInputUpdateDepositTree = () => {
    let proofFile = ReadJsonFile("./resources/updateRootDepositToCosmosBridge/proof.json");
    let publicFile = ReadJsonFile("./resources/updateRootDepositToCosmosBridge/public.json");
    let proof = [];
    for (let i = 0; i < 2; i++) {
        proof.push(proofFile.pi_a[i]);
    }
    for (let i = 0; i < 2; i++) {
        for (let j = 1; j >= 0; j--) {
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

const saveJsonData = (path, data) => {
    const jsonData = JSON.stringify(data, null, 2);
    fs.writeFileSync(path, jsonData, 'utf-8');
    console.log('Data has been saved to file:', path);
}

const QueryTxByHash = async(txHash) => {
    // console.log(resInitiate)
    // console.log(wasmCode)
    let res = await axios.get(geupdateDepositRootToCosmosBridgetTxAPI + "tx/v1beta1/txs/" + txHash).then( (response) => {
        // handle success
        return { tx: response.data.tx, height: response.data.tx_response.height }
        // console.dir(response.data.tx, { depth: null });
    })
        .catch( (error) => {
            // handle error
            console.log(error);
        })
        .finally( () => {
            // always executed
        });

    return res
}

const QueryBlockHeaderByHeight = async(height) => {
    let res = await axios.get(getTxAPI + "tx/v1beta1/txs/block/" + height).then( (response) => {
        // handle success
        return { header: response.data.block.header, txs: response.data.block.data.txs }
        // console.dir(response.data.tx, { depth: null });
    })
        .catch( (error) => {
            // handle error
            console.log(error);
        })
        .finally( () => {
            // always executed
        });

    return res
}

const QueryValidatorsSetByHeight = async(height) => {
    let res = await axios.get("https://testnet-rpc.orai.io/validators?height=" + height).then( (response) => {
        // handle success
        return response.data.result
        // console.dir(response.data.tx, { depth: null });
    })
        .catch( (error) => {
            // handle error
            console.log(error);
        })
        .finally( () => {
            // always executed
        });

    return res
}

const QueryBlockHeaderCommitByHeight = async(height) => {
    let res = await axios.get("https://testnet-rpc.orai.io/commit?height=" + height).then( (response) => {
        // handle success
        return response.data.result
        // console.dir(response.data.tx, { depth: null });
    })
        .catch( (error) => {
            // handle error
            console.log(error);
        })
        .finally( () => {
            // always executed
        });

    return res
}

const updateDepositTree = async() => {
    // const query = await client.getTx("2D925C0F81EF1E26662B0A2A9277180CE853F9F07C60CA2F3E64E7F565A19F78")
    const wallet = await getWallet();
    const client = await getClient();

    const senderAddress = (await wallet.getAccounts())[0].address;
    const contractAddress = process.env.COSMOS_BRIDGE || "";
    const msg = getInputUpdateDepositTree();
    console.log("msg", msg);
    const fee = "auto"
    const memo = null
    const res = await client.execute(senderAddress, contractAddress, msg, fee, memo)
    console.log(res)
    return res;
}

exports.updateDepositRootToCosmosBridge = async() => {
    const resUpdate = await updateDepositTree();
    console.log(resUpdate);
    writeToEnvFile("TX_HASH", resUpdate.transactionHash)

    const path = "./resources/updateRootDepositToCosmosBridge";
    const resQueryDepositRootTx = await QueryTxByHash(resUpdate.transactionHash);
    console.log(resQueryDepositRootTx);
    saveJsonData(`${path}/tx_data.json`, resQueryDepositRootTx);
    const resQueryBlock = await QueryBlockHeaderByHeight(resQueryDepositRootTx.height);
    saveJsonData(`${path}/block.json` , resQueryBlock);
    console.log(resQueryBlock)
    const resQueryValidatorsSet = await QueryValidatorsSetByHeight(resQueryDepositRootTx.height);
    console.log(resQueryValidatorsSet)
    saveJsonData(`${path}/validators.json`, resQueryValidatorsSet);
    const resQueryBlockHeaderCommit = await QueryBlockHeaderCommitByHeight(resQueryDepositRootTx.height);
    saveJsonData(`${path}/block_header_commit.json`);
}

exports.bridgeBlockHeader = async() => {
    execSync(`cd go/verifyBlock/ && go run .`);
    execSync(`cd python/ && python ed25519.py`);
    execSync(`node circom/validator/genBlockHeaderInput.js`);
    execSync(`cd circom/circuit/validators/validators-testnet/verifyblockheader_js && node generate_witness.js ./*.wasm ../input.json ../witness.wtns`);
    execSync(`cd circom/circuit/validators/validators-testnet/verifyblockheader_js && ../../../../rapidsnark/build/prover ../circuit_final.zkey ../witness.wtns ../../../../resources/cosmosHeader/proof.json ../../../../resources/cosmosHeader/public.json`);
    execSync(`node solidity/scripts/sdk/examples/updateBlockHeaderTestnet.js`);
}

exports.generateProofEth = async() => {
    execSync(`cd go/verifyTx/ && go run .`);
    execSync(`node circom/txs/genInputUpdateDepositRootToEthBridge.js`);
    execSync(`cd circom/circuit/verifyDepositRoot/rootdepositverifier_js/ && node generate_witness.js ./*.wasm ../input.json ../witness.wtns`);
    execSync(`cd circom/circuit/verifyDepositRoot/rootdepositverifier_js/ && ../../../../rapidsnark/build/prover ../circuit_final.zkey ../witness.wtns ../../../../resources/updateRootDepositToEthBridge/proof.json ../../../../resources/updateRootDepositToEthBridge/public.json`);
}

exports.updateDepositRootOnEth = async () => {
    execSync(`node solidity/scripts/sdk/examples/updateRootDeposit.js`);
}