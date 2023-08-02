//@ts-nocheck
import { SigningCosmWasmClient, Secp256k1HdWallet, toBinary } from "cosmwasm";
import * as dotenv from "dotenv";
import { Decimal } from "@cosmjs/math";
dotenv.config();
// This is your rpc endpoint

const rpcEndpoint = "https://testnet-rpc.orai.io:443/";
const mnemonic = process.env.MNEMONIC_COSMOS!;

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


async function supportTokenPair() {
    // const query = await client.getTx("2D925C0F81EF1E26662B0A2A9277180CE853F9F07C60CA2F3E64E7F565A19F78")
    const wallet = await getWallet();
    const client = await getClient();

    const senderAddress = (await wallet.getAccounts())[0].address;
    const contractAddress = process.env.COSMOS_BRIDGE || "";
    const msg = {
        support_token_pair: {
            destination_chainid: Number(process.env.ETH_CHAIN_ID),
            cosmos_token_address: process.env.COSMOS_TOKEN,
            eth_token_address: hexToDecimal(process.env.ETH_TOKEN)
        }
    };
    console.log("msg", msg);
    const fee = "auto"
    const memo: any = null
    const res = await client.execute(senderAddress, contractAddress, msg, fee, memo)
    console.log(res)
    return res;
}

async function main() {
    const resSupportTokenPair = await supportTokenPair();
    console.log(resSupportTokenPair);

}

main();