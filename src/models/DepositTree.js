const mongoose = require("mongoose");

const DepositTree = mongoose.Schema({
    root: String,
    n_leafs: Number,
    nqueue_leafs: Number
});

module.exports = mongoose.model("DepositTree", DepositTree);