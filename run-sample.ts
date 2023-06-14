import { execSync } from "child_process";

export async function runSample() {
    console.log("--deposit assets---")
    execSync(`npx ts-node rust/send-token.ts`);
    console.log("--end deposit----")

    // query deposit queue
    console.log("---state query deposit queue---")
    execSync(`npx ts-node rust/query-deposit-tree.ts`);
    console.log("---end query deposit queue---")

    // // gen input
    console.log("---start gen deposit proof---")
    execSync(`node circom/txs/genDepositTreeProof.js`);
    console.log("---")

    // // gen witness
    execSync(`cd circom/circuit/verifyRootBatchTxsCosmos/verifyrootbatchtxscosmos_js/ && node generate_witness.js ./*.wasm ../input.json ../witness.wtns`);

    // // gen proof
    execSync(`cd circom/circuit/verifyRootBatchTxsCosmos/verifyrootbatchtxscosmos_js/ && ../../../../rapidsnark/build/prover ../circuit_final.zkey ../witness.wtns ../../../../resources/updateRootDepositToCosmosBridge/proof.json ../../../../resources/updateRootDepositToCosmosBridge/public.json`);

    // // update deposit
    execSync(`npx ts-node rust/update-deposit-root.ts`);

    execSync('cd go/verifyTx/ && go run .');

    execSync(`node circom/txs/genInputUpdateDepositRootToEthBridge.js`);

    execSync(`cd circom/circuit/verifyDepositRoot/rootdepositverifier_js/ && node generate_witness.js ./*.wasm ../input.json ../witness.wtns`);

    execSync(`cd circom/circuit/verifyDepositRoot/rootdepositverifier_js/ && ../../../../rapidsnark/build/prover ../circuit_final.zkey ../witness.wtns ../../../../resources/updateRootDepositToEthBridge/proof.json ../../../../resources/updateRootDepositToEthBridge/public.json`);
    
    // update new deposit to eth
    execSync('npx hardhat compile');

    execSync(`node solidity/scripts/deploy-upgrades/deploy-cosmos-block-header.js`);
    execSync(`node solidity/scripts/sdk/examples/setLibAddressManager.js`);
    execSync(`node solidity/scripts/sdk/examples/updateRootDeposit.js`);

    // gen proof for user
    execSync(`node circom/txs/genClaimDepositTransaction.js`);
    execSync(`cd circom/circuit/verifyClaimTransaction/verifyclaimtransaction_js/ && node generate_witness.js ./*.wasm ../input.json ../witness.wtns`);
    execSync(`cd circom/circuit/verifyClaimTransaction/verifyclaimtransaction_js/ && ../../../../rapidsnark/build/prover ../circuit_final.zkey ../witness.wtns ../../../../resources/verifyClaimTransaction/proof.json ../../../../resources/verifyClaimTransaction/public.json`);
    
    // test withdraw
    execSync(`node solidity/scripts/sdk/examples/claimTransaction.js`);
}

runSample()