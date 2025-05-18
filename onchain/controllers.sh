#!/usr/bin/env bash

cardano-cli address build --payment-script-file controllers.script --mainnet > controllers.mainnet.address

cardano-cli address build --payment-script-file controllers.script --testnet-magic 1 > controllers.testnet.address

marlowe-cli util decode-bech32 $(cat controllers.mainnet.address) 2>/dev/null | tail -c +3 > controllers.hash

