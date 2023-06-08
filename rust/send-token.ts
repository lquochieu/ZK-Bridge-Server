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


async function sendToken(amount: String) {
    // const query = await client.getTx("2D925C0F81EF1E26662B0A2A9277180CE853F9F07C60CA2F3E64E7F565A19F78")
    const wallet = await getWallet();
    const client = await getClient();

    const senderAddress = (await wallet.getAccounts())[0].address;
    const contractAddress = process.env.COSMOS_TOKEN || "";
    const msg_bridge = {
        destination_chainid: 97,
        eth_bridge_address: hexToDecimal("0x66EeCaf1D881F7D828224422387b6fe0359AA69f"),
        eth_receiver: hexToDecimal("0x72e03B6E7AB9DdFe1699B65B8A46A3Cf30092020")
    }
    console.log("msg", msg_bridge);
    console.log("binary", toBinary(msg_bridge));
    const msg = {
        send: {
            contract: process.env.COSMOS_BRIDGE || "",
            amount: amount,
            msg: toBinary(msg_bridge)
        }
    }
    const fee = "auto"
    const memo: any = null
    const res = await client.execute(senderAddress, contractAddress, msg, fee, memo)
    console.log(res)
    return res;

}

async function main() {
    const resSendToken = await sendToken("10");
    console.log("sendtoken 0", resSendToken)
    const resSendToken1 = await sendToken("10");
    console.log("sendtoken 1", resSendToken1)
    const resSendToken2 = await sendToken("10");
    console.log("sendtoken 2", resSendToken2)

}
main();