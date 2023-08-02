const mongoose = require("mongoose");

const TokenPair = mongoose.Schema({
    cosmosTokenAddress: String,
    ethTokenAddress: String,
    cosmosTokenAddressInEth: String,
    symbol: String
});

module.exports = mongoose.model("TokenPair", TokenPair);