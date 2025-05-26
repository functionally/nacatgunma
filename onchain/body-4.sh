#!/usr/bin/env bash

set -e

EXE=../nacatgunma

PRIV_KEY=key-0.pem

COMMENT='Merge forks'

BLOCK=block-4

HEADER=header-4
HEADER_CBOR="$HEADER.cbor"

ACCEPT_CID="$(cat header-3.cid)"

BODY_CID=bafybeif7ztnhq65lumvvtr4ekcwd2ifwgm3awq4zfr3srh462rwyinlb4y

HEADER_CID=$(
"$EXE" header build \
  --key-file "$PRIV_KEY" \
  --accept "$ACCEPT_CID" \
  --accept bafyreihqvvqxrpb4htkuwchidf3qepvn7csyl7u7hwhcofz7amvniyeliy \
  --body  "$BODY_CID" \
  --comment "$COMMENT" \
  --schema "" \
  --media-type "application/x.empty" \
  --header-file "$HEADER_CBOR"
)

ipfs dag put --input-codec dag-cbor --store-codec dag-cbor --pin=false "$HEADER_CBOR"

echo "$HEADER_CID" > "$HEADER.cid"
echo "$BODY_CID" > "$BODY.cid"

ipfs cid format -b base16 -f %s "$HEADER_CID" | tail -c +2 > "$HEADER.cid16"

$EXE cardano inputs \
  --script \
  --credential-hash "$(cat controllers.hash)" \
  --header-cid "$HEADER_CID" \
  --datum-file "$BLOCK.datum" \
  --redeemer-file /dev/null \
  --metadata-file "$BLOCK.json"
