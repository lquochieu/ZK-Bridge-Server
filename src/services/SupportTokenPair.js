const { SigningCosmWasmClient, Secp256k1HdWallet, toBinary } = require("cosmwasm");
require("dotenv").config();
const { Decimal } = require("@cosmjs/math");


const rpcEndpoint = "https://testnet-rpc.orai.io:443/";
const mnemonic = process.env.MNEMONIC_COSMOS ?? "";

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

const hexToDecimal = (hex) => {
    if (hex.startsWith('0x')) {
        hex = hex.slice(2);
    }
    const decimalString = BigInt(`0x${hex}`).toString();
    return decimalString;
}

exports.supportTokenPair = async(cosmos_token_address, eth_token_address) => {
    const query = await client.getTx("2D925C0F81EF1E26662B0A2A9277180CE853F9F07C60CA2F3E64E7F565A19F78")
    const wallet = await getWallet();
    const client = await getClient();

    const senderAddress = (await wallet.getAccounts())[0].address;
    const contractAddress = process.env.COSMOS_BRIDGE || "";
    const msg = {
        support_token_pair: {
            destination_chainid: Number(process.env.ETH_CHAIN_ID),
            cosmos_token_address: cosmos_token_address,
            eth_token_address: hexToDecimal(eth_token_address)
        }
    };
    console.log("msg", msg);
    const fee = "auto"
    const memo = null;
    const res = await client.execute(senderAddress, contractAddress, msg, fee, memo)
    console.log(res)
    return res;
}