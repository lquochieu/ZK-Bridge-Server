const DepositInfor = require("../models/DepositInfor");
const DepositTree = require("../models/DepositTree");
const { generateProofUpdateRoot } = require("../services/GenerateProofUpdateRoot");
const { generateUserProof } = require("../services/GenerateUserProof");
const { queryDepositQueue } = require("../services/QueryDepositQueue");
const { updateDepositRootToCosmosBridge, bridgeBlockHeader, generateProofEth, updateDepositRootOnEth } = require("../services/UpdateDepositRoots");

exports.queryDepositQueue = async (req, res) => {
    try {
        const response = await queryDepositQueue();
        res.status(200).send(response);
    } catch (err) {
        console.log(err);
        res.status(500).send({ err: err.toString() });
    }
}

exports.updateRootDeposit = async (req, res) => {
    try {
        await queryDepositQueue();
        await generateProofUpdateRoot();
        await updateDepositRootToCosmosBridge();
        await bridgeBlockHeader();
        await generateProofEth();
        await updateDepositRootOnEth();
        res.status(200).send({});
    } catch (err) {
        console.log(err);
        res.status(500).send({ err: err.toString() });
    }
}

exports.genereateUserProof = async (req, res) => {
    try {
        const { key } = req.params;
        if (key == null) {
            throw ("Invalid key");
        }
        const response = await generateUserProof(Number(key));
        res.status(200).send(response);
    } catch (err) {
        console.log(err);
        res.status(500).send({ err: err.toString() });
    }
}

exports.queryDepositInforByUserAddress = async (req, res) => {
    try {
        let { sender, receiver } = req.query;
        if (sender == null) {
            sender = "";
        }
        if (receiver == null) {
            receiver = "";
        }
        if (sender != "") {
            sender = sender.toLowerCase();
        }
        if (receiver != "") {
            receiver = receiver.toLowerCase();
            if (receiver[1] != 'x')
                receiver = "0x" + receiver;
            receiver = BigInt(receiver).toString(10);
        }
        let response = [];
        if (receiver == "") {
            response = await DepositInfor.find({ sender: sender });
        } else {
            if (sender == "") {
                response = await DepositInfor.find({ eth_receiver: receiver });
            } else {
                response = await DepositInfor.find({ sender: sender, eth_receiver: receiver });
            }
        }
        res.status(200).send(response);
    } catch (err) {
        console.log(err);
        res.status(500).send({ err: err.toString() });
    }
}

exports.cronjobUpdate = async (req, res) => {
    try {
        console.log("------------cronjob auto update------------")
        await queryDepositQueue();
        await generateProofUpdateRoot();
        await updateDepositRootToCosmosBridge();
        await bridgeBlockHeader();
        await generateProofEth();
        await updateDepositRootOnEth();
        console.log("------------cronjob end update-------------")
        res.status(200).send({});

    } catch (err) {
        console.log(err);
        res.status(500).send({ err: err.toString() });

    }
}

exports.deleteDb = async (req, res) => {
    try {
        await DepositInfor.remove({});
        await DepositTree.remove({});
        const resposne = await DepositInfor.find();
        res.status(200).send({ resposne });

    } catch (err) {
        console.log(err);
        res.status(500).send({ err: err.toString() });

    }
}