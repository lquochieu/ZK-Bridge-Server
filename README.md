
# ZK-Bridge-Server
## Pre inquire
```bash
npm i
```

Download [zkey file](https://drive.google.com/file/d/1uZuEaCcZmhyCOMOnYJh5GZ7_pgR4SMlv/view?usp=sharing) and copy them to corressponding folder

Install [rapid sanrk](https://github.com/iden3/rapidsnark) here

## Activity Flow
### 1. Deposit asset
User deposit assets from Cosmos to ETH
```bash
npx ts-node rust/send-token.ts
```
Front-end will send server txHash of txs. Server store sender infos and txs info
### 2. Query deposit in queue
After a certain time, server query Cosmos Bridge Contract to get these transactions in DepositQueue which is deposited but don't include in Deposit Tree. To Query transactions:
```bash
npx ts-node rust/query-deposit-tree.ts
```
After that, server will save deposit info to a merkle tree. You can see example in file ./circom/txs/genDepositTreeProof.js

Besides server also store 2 attributes is: ***isInTree = false** and **isDeposited = false**. Note that, we only add 5 txs in a row

### 3. Gen proof for updating a new root deposit to Cosmos bridge
 Then, server creates a proof for updating new root DepositTree on Cosmos Bridge. To do that, first genarate input for verifyRooBatchTxsCosmos. Then generate proof with zkey file, wasm file and input file
```bash
# Gen input.json file for generating proof
node circom/txs/genDepositTreeProof.js

cd circom/circuit/verifyRootBatchTxsCosmos/verifyrootbatchtxscosmos_js/

# Gen witness
node generate_witness.js ./*.wasm ../input.json ../witness.wtns

# Gen proof with rapid sanrk
../../../../rapidsnark/build/prover ../circuit_final.zkey ../witness.wtns ../../../../resources/updateRootDepositToCosmosBridge/proof.json ../../../../resources/updateRootDepositToCosmosBridge/public.json

cd ....
```

### 4. Update root deposit to Cosmos bridge
After having proof, server executes function to update root deposit tree in Cosmos Bridge. After that, it receive a txHash and query block_header and tx_data
```bash
# In this command, we also query block_header, tx_data and validator_set after updating root
npx ts-node rust/update-deposit-root.ts
```
Then, update the isInTree of this txs is **True**

### 5. Bridge block_header
When server have block_header. It will create proof for updating new block_header to eth_bridge. 
```
cd go/verifyBlock
go run .
cd ../python
python ed25519.py
cd ..
node circom/validator/genBlockHeaderInput.js
cd circom/circuit/validators/validators-testnet/verifyblockheader_js
node generate_witness.js ./*.wasm ../input.json ../witness.wtns
../../../../rapidsnark/build/prover ../circuit_final.zkey ../witness.wtns ../../../../resources/cosmosHeader/proof.json ../../../../resources/cosmosHeader/public.json
```
After we have proof, pulish it to ETH Bridge
```
node solidity/scripts/sdk/examples/updateBlockHeaderTestnet.js
```
### 6. Generate proof for briding a new deposit root on ETH
Next, server will update deposit root to eth bridge. To do that, it need a proof for updating.
```bash
cd go/verifyTx/
#gen input with golang for updating deposit root to ETH Bridge
go run .

cd ...

#gen procees golang input to genarate input for updating deposit root to ETH Bridge
node circom/txs/genInputUpdateDepositRootToEthBridge.js

cd circom/circuit/verifyDepositRoot/rootdepositverifier_js/

# Gen witness
node generate_witness.js ./*.wasm ../input.json ../witness.wtns

# Gen proof with rapid sanrk
../../../../rapidsnark/build/prover ../circuit_final.zkey ../witness.wtns ../../../../resources/updateRootDepositToEthBridge/proof.json ../../../../resources/updateRootDepositToEthBridge/public.json

cd ....
```
### 7. Update a new deposit root on ETH
Then, server will interact with ETH Bridge contract to update new deposit root.
```bash
node solidity/scripts/sdk/examples/updateRootDeposit.js
```

### 8. Generate proof for user
When new root is updated, server update status of txs were bridged: **isBridged =  True**.  Server create for  user who has assset wasn't withdrawed a proof to receive their asset. Server bases on key of each DepositTree tx to gen proof. Example

```bash
#gen input for claim asset
node circom/txs/genClaimDepositTransaction.js

cd circom/circuit/verifyClaimTransaction/verifyclaimtransaction_js/

# Gen witness
node generate_witness.js ./*.wasm ../input.json ../witness.wtns

# Gen proof with rapid sanrk
../../../../rapidsnark/build/prover ../circuit_final.zkey ../witness.wtns ../../../../resources/verifyClaimTransaction/proof.json ../../../../resources/verifyClaimTransaction/public.json

cd ....
```

### 9. Now, we can test withdraw asset by command below:
```bash
node solidity/scripts/sdk/examples/claimTransaction.js
```





