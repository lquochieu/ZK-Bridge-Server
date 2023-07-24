const express = require("express");
const router = express.Router();
const Controller = require("../controller/Controller");

router.get("/query", Controller.queryDepositQueue);
router.put("/update", Controller.updateRootDeposit);
router.get("/proof/:key", Controller.genereateUserProof);
router.get("/infor", Controller.queryDepositInforByUserAddress);
router.delete("/db", Controller.deleteDb);
module.exports = router;