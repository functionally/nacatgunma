
export PATH=/extra/iohk/bin:$PATH
export CARDANO_NODE_SOCKET_PATH=/extra/iohk/networks/preprod/node.socket
export CARDANO_NODE_NETWORK_ID=1
NETWORK=testnet

cardano-cli conway query utxo --address $(cat nacatgunma.$NETWORK.address)

cardano-cli conway transaction build \
  --tx-in 1e6a020537ffc34b9e10cb8c4eab9c8cc95f5bbf6eb94ae2648c01c0a32e58d8#0 \
  --tx-in 82683878a41adcf19330e025b1fe08f342820ac487dc017b3a02c17951dcd630#0 \
  --tx-in ae1fa36e2f82b9d031169b8aaa2f139dcd62d93820122c3ff37cd0ee4386ef64#0 \
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
