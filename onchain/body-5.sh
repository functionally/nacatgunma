#!/usr/bin/env bash

set -e

EXE=../nacatgunma

PRIV_KEY=key-0.pem

COMMENT='Update to Nacatgunma explorer web application.'

BLOCK=block-5

HEADER=header-5
HEADER_CBOR="$HEADER.cbor"

ACCEPT_CID="$(cat header-4.cid)"

BODY_CID=bafybeieuxid7jfvcqs2kv74t5fcwgdw26ja4cexe6hwpkuhbrcudaj5nfm

HEADER_CID=$(
"$EXE" header build \
  --key-file "$PRIV_KEY" \
  --accept "$ACCEPT_CID" \
  --reject bafyreiampkqrvb7xkgfnvs5ca6y4oashcyevrezl4xnboige7wymflpcc4 \
  --body  "$BODY_CID" \
  --comment "$COMMENT" \
  --schema "https://schema.org/WebApplication" \
  --media-type "text/directory" \
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
