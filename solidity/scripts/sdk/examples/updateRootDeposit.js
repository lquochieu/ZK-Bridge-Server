const {updateRootDepositTree} = require("../oraisan-bridge")

require("dotenv").config();

const main = async () => {
  
   let root = await updateRootDepositTree("./resources/updateRootDepositToEthBridge/public.json", "./resources/updateRootDepositToEthBridge/proof.json");
   console.log("RootDeposit", root);
};

main()
  .then(() => {})
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
