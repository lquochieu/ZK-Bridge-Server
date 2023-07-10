const express = require("express");
const router = express.Router();
const Controller = require("../controller/Controller");

router.get("/query", Controller.queryDepositQueue);
router.put("/update", Controller.updateRootDeposit);
router.get("/proof/:key", Controller.genereateUserProof);
router.get("/infor/:address", Controller.queryDepositInforByUserAddress);

module.exports = router;