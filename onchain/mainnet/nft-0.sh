
export PATH=/extra/iohk/bin:$PATH
export CARDANO_NODE_SOCKET_PATH=/extra/iohk/networks/mainnet/node.socket
export CARDANO_NODE_NETWORK_ID=mainnet
NETWORK=mainnet

cardano-cli conway query utxo --address $(cat nacatgunma.$NETWORK.address)

cardano-cli conway transaction build \
  --tx-in 069bc8f72c074c9c4dc2b25a0fcdc0158a88260a37b63a9f16beca04ff6e227f#1 \
  --tx-in ad4f9d12f8f9ffd80e5352c1d40137052ca0ee6e608c1eea645346b1e616d12f#0 \
  --tx-out "$(cat nacatgunma.$NETWORK.address)+1500000+1 $(cat controllers.hash).4e6163617467756e6d61" \
  --change-address $(cat nacatgunma.$NETWORK.address) \
  --mint "1 $(cat controllers.hash).4e6163617467756e6d61" \
    --mint-script-file controllers.script \
  --invalid-before 58312 \
  --json-metadata-no-schema \
  --metadata-json-file nft-0.json \
  --out-file nft-0.unsigned

cardano-cli conway transaction sign \
  --tx-body-file nft-0.unsigned \
  --out-file nft-0.signed \
  --signing-key-file nacatgunma.skey

cardano-cli conway transaction submit \
  --tx-file nft-0.signed

cardano-cli conway transaction txid \
  --tx-file nft-0.signed

cardano-cli conway query utxo --address $(cat nacatgunma.$NETWORK.address)
