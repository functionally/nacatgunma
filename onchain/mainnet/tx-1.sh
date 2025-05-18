export PATH=/extra/iohk/bin:$PATH
export CARDANO_NODE_SOCKET_PATH=/extra/iohk/networks/mainnet/node.socket
export CARDANO_NODE_NETWORK_ID=mainnet
NETWORK=mainnet

cardano-cli conway query utxo --address $(cat nacatgunma.$NETWORK.address)

cardano-cli conway query utxo --address $(cat controllers.$NETWORK.address)

cardano-cli conway query utxo --address $(cat script-0.$NETWORK.address)

cardano-cli conway transaction build \
  --tx-in-collateral 069bc8f72c074c9c4dc2b25a0fcdc0158a88260a37b63a9f16beca04ff6e227f#1 \
  --tx-in f5adb165a5e32d4cacfa7e3e87ba126832ee2df23a24d45ed7e306b003704da2#1 \
    --tx-in-script-file controllers.script \
  --tx-in f5adb165a5e32d4cacfa7e3e87ba126832ee2df23a24d45ed7e306b003704da2#0 \
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
