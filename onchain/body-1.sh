#!/usr/bin/env bash

set -e

EXE=../nacatgunma

PRIV_KEY=key-0.pem

COMMENT="Specification block for the Nacatgunma blockchain."

BLOCK=block-1

BODY=body-1
BODY_IRI=urn:uuid:e7c8a7a8-eecb-4474-af36-a0ca474a2af5
BODY_TTL="$BODY.ttl"
BODY_NQ="$BODY.nq"
BODY_CBOR="$BODY.cbor"

HEADER=header-1
HEADER_CBOR="$HEADER.cbor"

ACCEPT_CID="$(cat header-0.cid)"

rapper -q -i turtle -o nquads "$BODY_TTL" "$BODY_IRI" > "$BODY_NQ"

BODY_CID=$("$EXE" body rdf --base-uri "$BODY_IRI" --rdf-file "$BODY_NQ" --body-file "$BODY_CBOR")

HEADER_CID=$("$EXE" header build --key-file "$PRIV_KEY" --accept "$ACCEPT_CID" --body "$BODY_CID" --comment "$COMMENT" --header-file "$HEADER_CBOR")

ipfs dag put --input-codec dag-cbor --store-codec dag-cbor --pin=false "$HEADER_CBOR" "$BODY_CBOR"

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
