#!/usr/bin/env nix-shell
#!nix-shell -I nixpkgs=https://github.com/NixOS/nixpkgs/archive/nixos-25.05.tar.gz
#!nix-shell -i bash -p util-linux kubo jq librdf_raptor2

set -e

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

export PATH="/extra/iohk/bin:$DIR/..:$PATH"

export CARDANO_NODE_SOCKET_PATH=/extra/iohk/networks/mainnet/node.socket
export CARDANO_NODE_NETWORK_ID=mainnet
NETWORK=mainnet

PLUTUS_FILE="$DIR/script-0.plutus"
PLUTUS_ADDR=$(cat "${PLUTUS_FILE%%.plutus}.$NETWORK.address")
TIP=$(nacatgunma cardano tips --node-socket "$CARDANO_NODE_SOCKET_PATH" --script-address $PLUTUS_ADDR | jq '.[0]')

### BEGIN CUSTOMIZATION ###
SUBMIT=0
PIN=0
SUFFIX=13
TYPE=EMPTY
### END CUSTOMIZATION ###

case "$TYPE" in

  TURTLE)
    ### BEGIN CUSTOMIZATION ###
    COMMENT=
    BODY_IRI=urn:uuid:$(uuidgen)
    BODY=
    ### END CUSTOMIZATION ###
    BODY_TTL="$BODY.ttl"
    BODY_NQ="$BODY.nq"
    BODY_CBOR="$BODY.cbor"
    rapper -q -i turtle -o nquads "$BODY_TTL" "$BODY_IRI" > "$BODY_NQ"
    BODY_CID=$(nacatgunma body rdf --base-uri "$BODY_IRI" --rdf-file "$BODY_NQ" --body-file "$BODY_CBOR")
    ipfs dag put --input-codec dag-cbor --store-codec dag-cbor --pin=false "$BODY_CBOR"
    MEDIA_TYPE=application/vnd.ipld.dag-cbor
    SCHEMA_URI=https://w3c.github.io/json-ld-cbor/
    ;;
    
  JOSE)
    ### BEGIN CUSTOMIZATION ###
    COMMENT=
    BODY_PLAINTEXT_FILE=
    BODY_CIPHERTEXT_FILE="$BODY_PLAINTEXT_FILE.jwe"
    BODY_TGDH_KEY="$DIR/brio-20250704a.tgdh.pri"
    ### END CUSTOMIZATION ###
    nacatgunma body tgdh encrypt \
      --private-file "$BODY_TGDH_KEY" \
      --plaintext-file "$BODY_PLAINTEXT_FILE" \
      --content-type "$CONTENT_TYPE" \
      --jwe-file "$BODY_CIPHERTEXT_FILE"
    BODY_CID=$(ipfs add ---pin=false --cid-version 1 "$BODY_CIPHERTEXT_FILE")
    MEDIA_TYPE=application/jose+json
    SCHEMA_URI=https://github.com/functionally/nacatgunma/blob/37b3f8da9e81d3c09ee9115023510ad0d632a4a9/TGDH-Specialization.md
    ;;
    
  FILE)
    ### BEGIN CUSTOMIZATION ###
    COMMENT=
    MEDIA_TYPE=
    SCHEMA_URI=
    BODY_FILE=
    ### END CUSTOMIZATION ###
    BODY_CID=$(ipfs add ---pin=false --cid-version 1 "$BODY_FILE")
    ;;
    
  FOLDER)
    ### BEGIN CUSTOMIZATION ###
    COMMENT=
    SCHEMA_URI=
    BODY_FOLDER=
    ### END CUSTOMIZATION ###
    BODY_CID=$(ipfs add --recursive=true --pin=false --cid-version 1 "$BODY_FOLDER")
    MEDIA_TYPE=text/directory
    ;;
    
  WEB)
    ### BEGIN CUSTOMIZATION ###
    COMMENT="Update to Nacatgunma explorer web application"
    BODY_FOLDER="$DIR/../web"
    ### END CUSTOMIZATION ###
    BODY_CID=$(ipfs add --recursive=true --pin=false --cid-version 1 "$BODY_FOLDER")
    MEDIA_TYPE=text/directory
    SCHEMA_URI=https://schema.org/WebApplication
    ;;
    
  EMPTY)
    ### BEGIN CUSTOMIZATION ###
    COMMENT="Merge forks"
    ### END CUSTOMIZATION ###
    BODY_CID=bafybeif7ztnhq65lumvvtr4ekcwd2ifwgm3awq4zfr3srh462rwyinlb4y
    MEDIA_TYPE=application/x.empty
    SCHEMA_URI=
    ;;

  *)
    false
    ;;
esac

BLOCK=block-$SUFFIX
HEADER=header-$SUFFIX

PRIV_KEY="$DIR/key-0.pem"

HEADER_CBOR="$HEADER.cbor"

ACCEPT_CID=$(echo $TIP | jq -r .HeaderCid)

HEADER_CID=$(
  nacatgunma header build \
    --key-file "$PRIV_KEY" \
    --accept "$ACCEPT_CID" \
    --body "$BODY_CID" \
    --comment "$COMMENT" \
    --media-type "$MEDIA_TYPE" \
    --schema "$SCHEMA_URI" \
    --header-file "$HEADER_CBOR"
)

ipfs dag put --input-codec dag-cbor --store-codec dag-cbor --pin=false "$HEADER_CBOR"

if [ "$PIN" == "1" ]
then
  ipfs pin add --recursive --name=nacatgunma-$SUFFIX.header $HEADER_CID
else
  echo ipfs pin add --recursive --name=nacatgunma-$SUFFIX.header $HEADER_CID
fi

ipfs cid format -b base16 -f %s "$HEADER_CID" | tail -c +2 > "$HEADER.cid16"
echo "$HEADER_CID" > "$HEADER.cid"
echo "$BODY_CID" > "$BODY.cid"

nacatgunma cardano inputs \
  --script \
  --credential-hash "$(cat controllers.hash)" \
  --header-cid "$HEADER_CID" \
  --datum-file "$BLOCK.datum" \
  --redeemer-file /dev/null \
  --metadata-file "$BLOCK.json"

COLLATERAL_SKEY=nacatgunma.skey
COLLATERAL_ADDR=$(cat "${COLLATERAL_SKEY%%.skey}.$NETWORK.address")
COLLATERAL_TXIN=$(
cardano-cli conway query utxo --address $COLLATERAL_ADDR --output-json \
| jq -r '. | to_entries | map(select(.value.datum == null and (.value.value | keys) == ["lovelace"])) | .[0] | .key'
)

CONTROLER_SCRIPT=controllers.script
CONTROLER_ADDR=$(cat "${CONTROLER_SCRIPT%%.script}.$NETWORK.address")
CONTROLER_HASH=$(cat "${CONTROLER_SCRIPT%%.script}.hash")
TOKEN="1 $CONTROLER_HASH.4e6163617467756e6d61"
CONTROLER_TXIN=$(
cardano-cli conway query utxo --address $CONTROLER_ADDR --output-json \
| jq -r '. | to_entries | map(select(.value.datum == null and (.value.value | keys) == ["lovelace"])) | .[0] | .key'
)

PLUTUS_REDEEMER="${PLUTUS_FILE%%.plutus}.redeemer"
PLUTUS_TXIN=$(echo $TIP | jq -r '.TxId')

cardano-cli conway transaction build \
  --tx-in-collateral "$COLLATERAL_TXIN"\
  --tx-in "$CONTROLER_TXIN" \
    --tx-in-script-file "$CONTROLER_SCRIPT" \
  --tx-in "$PLUTUS_TXIN" \
    --tx-in-script-file "$PLUTUS_FILE" \
    --tx-in-inline-datum-present \
    --tx-in-redeemer-file "$PLUTUS_REDEEMER" \
  --tx-out "$PLUTUS_ADDR+1500000+$TOKEN" \
    --tx-out-inline-datum-file "$BLOCK.datum" \
  --change-address "$CONTROLER_ADDR" \
  --invalid-before 58312 \
  --json-metadata-no-schema \
  --metadata-json-file "$BLOCK.json" \
  --out-file tx-$SUFFIX.unsigned

cardano-cli conway transaction sign \
  --tx-body-file tx-$SUFFIX.unsigned \
  --out-file tx-$SUFFIX.signed \
  --signing-key-file "$COLLATERAL_SKEY"

TXID="$(cardano-cli conway transaction txid --tx-file tx-$SUFFIX.signed)#0"

if [ "$SUBMIT" == "1" ]
then
  cardano-cli conway transaction submit --tx-file tx-$SUFFIX.signed
  while true
  do
    echo 'Waiting for '"$TXID"' . . .'
    sleep 20s
    if [ "$(cardano-cli conway query utxo --address $PLUTUS_ADDR --output-json | jq 'select(."'"$TXID"'")' | wc -l)" != "0" ]
    then
      echo ' . . . confirmed.'
      break
    fi
  done
else
  echo 'Transaction '"$TXID"' not yet submitted.'
fi
