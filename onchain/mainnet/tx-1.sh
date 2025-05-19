export PATH=/extra/iohk/bin:$PATH
export CARDANO_NODE_SOCKET_PATH=/extra/iohk/networks/mainnet/node.socket
export CARDANO_NODE_NETWORK_ID=mainnet
NETWORK=mainnet

cardano-cli conway query utxo --address $(cat nacatgunma.$NETWORK.address)

cardano-cli conway query utxo --address $(cat controllers.$NETWORK.address)

cardano-cli conway query utxo --address $(cat script-0.$NETWORK.address)

cardano-cli conway transaction build \
  --tx-in-collateral 8042f83c9c52ce173750f79d0163b70ad127ee6c225cd2a57b1249689c1cf164#0 \
  --tx-in 8716dc99ad7051047a6fa8f872e3023ada692b71f2c71ffc94988448ed3d9341#1 \
    --tx-in-script-file controllers.script \
  --tx-in 8716dc99ad7051047a6fa8f872e3023ada692b71f2c71ffc94988448ed3d9341#0 \
    --tx-in-script-file script-0.plutus \
    --tx-in-inline-datum-present \
    --tx-in-redeemer-file script-0.redeemer \
  --tx-out "$(cat script-0.$NETWORK.address)+1500000+1 $(cat controllers.hash).4e6163617467756e6d61" \
    --tx-out-inline-datum-file block-1.datum \
  --change-address $(cat controllers.$NETWORK.address) \
  --invalid-before 58312 \
  --json-metadata-no-schema \
  --metadata-json-file block-1.json \
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
