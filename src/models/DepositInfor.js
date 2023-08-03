const mongoose = require("mongoose");

const DepositInfor = mongoose.Schema({
    is_deposit: Boolean,
    sender: String,
    destination_chainid: Number,
    eth_bridge_address: String,
    eth_receiver: String,
    amount: String,
    eth_token_address: String,
    cosmos_token_address: String,
    key: Number,
    value: String
});

module.exports = mongoose.model("DepositInfor", DepositInfor);