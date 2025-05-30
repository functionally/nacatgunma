export PATH=/extra/iohk/bin:$PATH
export CARDANO_NODE_SOCKET_PATH=/extra/iohk/networks/preprod/node.socket
export CARDANO_NODE_NETWORK_ID=1
NETWORK=testnet

cardano-cli conway query utxo --address $(cat nacatgunma.$NETWORK.address)

cardano-cli conway query utxo --address $(cat controllers.$NETWORK.address)

cardano-cli conway query utxo --address $(cat script-0.$NETWORK.address)

cardano-cli conway transaction build \
  --tx-in-collateral 7d5a14c0e4ee7ba23ac6c1430b97e6f18c6f4f9a2e397a118fb089b17f8ba7f4#1 \
  --tx-in c4fb42bdca24628f96d3268c71253ebb9565453bf29de9b4ace70521a9bb9959#1 \
    --tx-in-script-file controllers.script \
  --tx-in c4fb42bdca24628f96d3268c71253ebb9565453bf29de9b4ace70521a9bb9959#0 \
    --tx-in-script-file script-0.plutus \
    --tx-in-inline-datum-present \
    --tx-in-redeemer-file script-0.redeemer \
  --tx-out "$(cat nacatgunma.$NETWORK.address)+2000000" \
  --mint "-1 $(cat controllers.hash).4e6163617467756e6d61" \
  --mint-script-file controllers.script \
  --change-address $(cat controllers.$NETWORK.address) \
  --invalid-before 58312 \
  --out-file tx-1.unsigned

cardano-cli conway transaction sign \
  --tx-body-file tx-1.unsigned \
  --out-file tx-1.signed \
  --signing-key-file nacatgunma.skey

cardano-cli conway transaction submit \
  --tx-file tx-1.signed

cardano-cli conway transaction txid \
  --tx-file tx-1.signed

cardano-cli conway query utxo --address $(cat script-0.$NETWORK.address)
