exports.registerCosmosBridge = exports.registerTokenPair = exports.updateRootDepositTree = exports.claimTransaction = void 0;
const { rdOwnerOraisanBridge } = require("./rdOwner");
const { getProofUpdateRootDeposit, getProofClaimTransaction } = require("../utils/getProof");
require("dotenv").config();

const isSentProof = async (
    eth_bridge_address,
    eth_receiver,
    amount,
    eth_token_address,
    key
) => {
    console.log("eth_bridge_address", eth_bridge_address)
    console.log("eth_receiver", eth_receiver)
    console.log("amount", amount)
    console.log("eth_token_address", eth_token_address)
    console.log("key", key)
    const rdOwner = await rdOwnerOraisanBridge();
    const msg = await rdOwner.encodeProof(
        eth_bridge_address,
        eth_receiver,
        amount,
        eth_token_address,
        key
    );
    console.log("msg", msg)
    return await rdOwner.sentProof(msg);
}
exports.isSentProof = isSentProof;

const registerCosmosBridge = async (_cosmosTokenAddress, _ethTokenAddress) => {
    const rdOwner = await rdOwnerOraisanBridge();
    const res = await rdOwner.registerCosmosBridge(_cosmosTokenAddress, _ethTokenAddress);
    await res.wait();

    return res;
}
exports.registerCosmosBridge = registerCosmosBridge;

const registerTokenPair = async (_cosmosTokenAddress, _ethTokenAddress) => {
    const rdOwner = await rdOwnerOraisanBridge();
    const res = await rdOwner.registerTokenPair(_cosmosTokenAddress, _ethTokenAddress);
    await res.wait();
    const ethTokenAddress = await rdOwner.getTokenPairOfCosmosToken(_cosmosTokenAddress);
    return ethTokenAddress;
}
exports.registerTokenPair = registerTokenPair;

const deleteTokenPair = async (_cosmosTokenAddress) => {
    const rdOwner = await rdOwnerOraisanBridge();
    const res = await rdOwner.deleteTokenPair(_cosmosTokenAddress);
    await res.wait();
    const ethTokenAddress = await rdOwner.getTokenPairOfCosmosToken(_cosmosTokenAddress);
    return ethTokenAddress;
}
exports.deleteTokenPair = deleteTokenPair;

const updateRootDepositTree = async (pathInput, pathProof) => {
    const input = getProofUpdateRootDeposit(pathInput, pathProof);
    console.log(input)
    const rdOwner = await rdOwnerOraisanBridge();

    const res = await rdOwner.updateRootDepositTree([
        input.optionName,
        input.pi_a,
        input.pi_b,
        input.pi_c,
        input.cosmosSender,
        input.cosmosBridge,
        input.depositRoot,
        input.dataHash
    ], {gasLimit: 2e6});
    await res.wait();
    return (await rdOwner.getLastDepositRoot());
}
exports.updateRootDepositTree = updateRootDepositTree;

const claimTransaction = async (pathInput, pathProof) => {
    const input = getProofClaimTransaction(pathInput, pathProof);
    console.log(input)
    const rdOwner = await rdOwnerOraisanBridge();
    const res = await rdOwner.claimTransaction([
        input.optionName,
        input.pi_a,
        input.pi_b,
        input.pi_c,
        input.eth_bridge_address,
        input.eth_receiver,
        input.amount,
        input.eth_token_address,
        input.key,
        input.depositRoot
    ], {gasLimit: 2e6});
    await res.wait();
}
exports.claimTransaction = claimTransaction;