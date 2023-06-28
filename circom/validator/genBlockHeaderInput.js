const utils = require('../utils');
const fs = require("fs");
const input_python = require("../../resources/cosmosHeader/input_python.json");
const { saveJsonData } = require('../helper');

function hexStringToArray(str) {
    const bytesArray = [];

    for (let i = 0; i < str.length; i += 2) {
        bytesArray.push(parseInt(str.slice(i, i + 2), 16));
    }
    return bytesArray;
}

const main = async () => {
    const pointA = input_python.A.map(e => BigInt(e));
    const pointR = input_python.R.map(e => BigInt(e));
    // console.log("pA", utils.point_compress(BigInt("0x" + input_python.pubKeys)))
    // console.log("pR", utils.point_compress(BigInt("0x" + input_python.R8)))

    const chunkA = [];
    const chunkR = [];

    for (let i = 0; i < 4; i++) {
        chunkA.push(utils.chunkBigInt(pointA[i], BigInt(2 ** 85)));
        chunkR.push(utils.chunkBigInt(pointR[i], BigInt(2 ** 85)));
    }

    for (let i = 0; i < 4; i++) {
        utils.pad(chunkA[i], 3);
        utils.pad(chunkR[i], 3);
    }

    const input = {
        height: input_python.height,
        dataHash: hexStringToArray(input_python.dataHash),
        parrentSiblings: input_python.parrentSiblings.map(e => hexStringToArray(e)),
        blockHash: hexStringToArray(input_python.blockHash),
        blockTime: input_python.blockTime,
        partsTotal: input_python.partsTotal,
        partsHash: hexStringToArray(input_python.partsHash),
        sigTimeSeconds: input_python.sigTimeSeconds,
        sigTimeNanos: input_python.sigTimeNanos,
        pubKeys: hexStringToArray(input_python.pubKeys),
        votingPowers: input_python.votingPowers,
        R8: hexStringToArray(input_python.R8),
        S: hexStringToArray(input_python.S),
        PointA: chunkA,
        PointR: chunkR,
    }

    console.log(input)
    saveJsonData('./circom/circuit/validators/validators-testnet/input.json', input)

}
main()
    .then(() => { })
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });