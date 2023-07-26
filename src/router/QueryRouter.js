const express = require("express");
const router = express.Router();
const Controller = require("../controller/Controller");
const TokenController = require("../controller/TokenController");

router.get("/query", Controller.queryDepositQueue);
router.put("/update", Controller.updateRootDeposit);
router.get("/proof/:key", Controller.genereateUserProof);
router.get("/infor", Controller.queryDepositInforByUserAddress);
router.delete("/db", Controller.deleteDb);

router.post("/token-pair", TokenController.addToken);
router.get("/token-pair", TokenController.getListTokenPair);
router.delete("/token-pair", TokenController.removeTokenPair);

module.exports = router;