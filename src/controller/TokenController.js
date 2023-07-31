const {getAddresFromAsciiString} = require("../../solidity/scripts/utils/helper");
const {registerTokenPair} = require("../../solidity/scripts/sdk/oraisan-bridge");
const { supportTokenPair } = require("../services/SupportTokenPair");
const TokenPair = require("../models/TokenPair");
require("dotenv").config();

exports.addToken = async (req, res) => {
    try {
        const { cosmosTokenAddress, ethTokenAddress, symbol } = req.body;
        const cosmosTokenAddressInEth = getAddresFromAsciiString(cosmosTokenAddress);
        const ethToken = await registerTokenPair(cosmosTokenAddressInEth, ethTokenAddress);
        console.log(ethToken);
        await supportTokenPair(cosmosTokenAddress, ethTokenAddress);

        const tokenPair = new TokenPair({
            cosmosTokenAddress: cosmosTokenAddress,
            ethTokenAddress: ethTokenAddress,
            cosmosTokenAddressInEth: cosmosTokenAddressInEth,
            symbol: symbol
        });

        await tokenPair.save();
        res.status(200).send(tokenPair);
    } catch (err) {
        console.log(err);
        res.status(500).send({ err: err.toString() });
    }
}

exports.getListTokenPair = async (req, res) => {
    try {
        let {cosmosTokenAddress, ethTokenAddress, symbol} = req.query;
        let query = {};
        if (cosmosTokenAddress != undefined || ethTokenAddress != undefined || symbol != undefined) {
            query = {
                "$or": [
                    {
                        cosmosTokenAddress: cosmosTokenAddress
                    },
                    {
                        ethTokenAddress: ethTokenAddress
                    },
                    {
                        symbol: symbol
                    }
                ]
            };
            console.log(query);
        }

        const listToken = await TokenPair.find(query);

        res.status(200).send(listToken);

    } catch (err) {
        console.log(err);
        res.status(500).send({ err: err.toString() });
    }
}

exports.removeTokenPair = async (req, res) => {
    try {
        let {cosmosTokenAddress, ethTokenAddress, symbol} = req.query;
        let query = {};
        if (cosmosTokenAddress != null || ethTokenAddress != null || symbol == null) {
            query = {
                "$or": [
                    {
                        cosmosTokenAddress: cosmosTokenAddress
                    },
                    {
                        ethTokenAddress: ethTokenAddress
                    },
                    {
                        symbol: symbol
                    }
                ]
            };
        }

        const listToken = await TokenPair.find(query);

        if (listToken.length > 1) {
            throw("List tokenpair greater than 1");
        }
        for(let i = 0; i < listToken.length; i++) {
            listToken[i].remove();
        }
        res.status(200).send(listToken);

    } catch (err) {
        console.log(err);
        res.status(500).send({ err: err.toString() });
    }
}