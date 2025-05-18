
export PATH=/extra/iohk/bin:$PATH
export CARDANO_NODE_SOCKET_PATH=/extra/iohk/networks/preprod/node.socket
export CARDANO_NODE_NETWORK_ID=1
NETWORK=testnet

cardano-cli conway query utxo --address $(cat nacatgunma.$NETWORK.address)

cardano-cli conway query utxo --address $(cat controllers.$NETWORK.address)

cardano-cli conway transaction build \
  --tx-in 0b863a6ca567b3eed6f703a6546eaade7e7d9aeea8ccc182067300f92348c95b#1 \
    --tx-in-script-file controllers.script \
  --tx-out "$(cat script-0.$NETWORK.address)+2000000+1 $(cat controllers.hash).4e6163617467756e6d61" \
  --tx-out-datum-embed-file block-0.datum \
  --change-address $(cat controllers.$NETWORK.address) \
  --mint "1 $(cat controllers.hash).4e6163617467756e6d61" \
  --mint-script-file controllers.script \
  --invalid-before 58312 \
  --json-metadata-no-schema \
  --metadata-json-file block-0.json \
  --required-signer-hash $(cat nacatgunma.hash) \
  --out-file tx-0.unsigned

cardano-cli conway transaction sign \
  --tx-body-file tx-0.unsigned \
  --out-file tx-0.signed \
  --signing-key-file nacatgunma.skey

cardano-cli conway transaction submit \
  --tx-file tx-0.signed

cardano-cli conway transaction txid \
  --tx-file tx-0.signed

cardano-cli conway query utxo --address $(cat script-0.$NETWORK.address)
