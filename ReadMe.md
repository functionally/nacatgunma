# Nacatgunma

Nacatgunma is an experimental blockchain that supports fluid consensus.


## CLI examples


### Build the CLI executable

```bash
go build -o nacatgunma main.go
```


### Generate a private key

```bash
nacatgunma key generate \
  --key-file private.pem
```

```console
did:key:z6Mkrpqbu1hJWGHCuTa9Z9SDaCp8VJa82dLYhhCWYiyye6Nr
```

```bash
cat private.pem
```

```console
-----BEGIN ED25519 PRIVATE KEY-----
MC4CAQAwBQYDK2VwBCIEIDwz2pHLuwjP+rIrsLEwWhxoHq5iyZvWGFy/k44sHFCR
-----END ED25519 PRIVATE KEY-----
```


### Resolve the DID for a public key

```bash
nacatgunma key resolve \
  --key-did "did:key:z6Mkrpqbu1hJWGHCuTa9Z9SDaCp8VJa82dLYhhCWYiyye6Nr" \
  --output-file resolution.json

json2yaml resolution.json
```

```console
Context:
- https://w3id.org/did-resolution/v1
DIDDocument:
  '@context':
  - https://w3id.org/did/v1
  id: did:key:z6Mkrpqbu1hJWGHCuTa9Z9SDaCp8VJa82dLYhhCWYiyye6Nr
  verificationMethod:
  - controller: did:key:z6Mkrpqbu1hJWGHCuTa9Z9SDaCp8VJa82dLYhhCWYiyye6Nr
    id: did:key:z6Mkrpqbu1hJWGHCuTa9Z9SDaCp8VJa82dLYhhCWYiyye6Nr#z6Mkrpqbu1hJWGHCuTa9Z9SDaCp8VJa82dLYhhCWYiyye6Nr
    publicKeyBase58: DNaZJmSsAinjnxjSsaUNj7G8fjJGck6C1gHaiT1xisbU
    type: Ed25519VerificationKey2018
  authentication:
  - did:key:z6Mkrpqbu1hJWGHCuTa9Z9SDaCp8VJa82dLYhhCWYiyye6Nr#z6Mkrpqbu1hJWGHCuTa9Z9SDaCp8VJa82dLYhhCWYiyye6Nr
  assertionMethod:
  - did:key:z6Mkrpqbu1hJWGHCuTa9Z9SDaCp8VJa82dLYhhCWYiyye6Nr#z6Mkrpqbu1hJWGHCuTa9Z9SDaCp8VJa82dLYhhCWYiyye6Nr
  capabilityDelegation:
  - did:key:z6Mkrpqbu1hJWGHCuTa9Z9SDaCp8VJa82dLYhhCWYiyye6Nr#z6Mkrpqbu1hJWGHCuTa9Z9SDaCp8VJa82dLYhhCWYiyye6Nr
  capabilityInvocation:
  - did:key:z6Mkrpqbu1hJWGHCuTa9Z9SDaCp8VJa82dLYhhCWYiyye6Nr#z6Mkrpqbu1hJWGHCuTa9Z9SDaCp8VJa82dLYhhCWYiyye6Nr
  keyAgreement:
  - controller: did:key:z6Mkrpqbu1hJWGHCuTa9Z9SDaCp8VJa82dLYhhCWYiyye6Nr
    id: did:key:z6Mkrpqbu1hJWGHCuTa9Z9SDaCp8VJa82dLYhhCWYiyye6Nr#z6LSpeMaM7nxQ7xYbC4wK2nNxZwqDpKQRNcPy2u3E9nvN961
    publicKeyBase58: DyBQpoz6JfEoVohAnPGRdyjMNfnHimSF64BMjh9PemKF
    type: X25519KeyAgreementKey2019
  created: '2025-04-26T08:00:57.983089098-06:00'
  updated: '2025-04-26T08:00:57.983089098-06:00'
DocumentMetadata: null
```


### Build a block header

```bash
nacatgunma header build \
  --key-file private.pem \
  --body bafyreiea2su23cm4nbfl3675m442gp5yo5qmghspjikeeeioudyls2jjtm \
  --accept bafyreid32eo34hcksuilttqhoyhssoz36r5umtk6zallvgoxqlehsqltru \
  --header-file header.cbor
```

```console
bafyreib5fuk4qex34is3pt52ij4jddlnsevkys7jwa6v2lp2qrs2eoq5he
```

```bash
cbordump header.cbor 
```

```console
{"Issuer": "did:key:z6Mkrpqbu1hJWGHCuTa9Z9SDaCp8VJa82dLYhhCWYiyye6Nr", "Payload": {"Body": 42(h'000171122080d4a9ad899c684abdfbfd6739a33fb87760c31e4f4a1442110ea0f0b969299b'), "Accept": [42(h'00017112207bd11dbe1c4a9510b9ce07760f293b3bf47b464d5ec816ba99d782c87941738d')], "Reject": [], "Schema": "DAG-CBOR", "Version": 1, "MediaType": "application/cbor"}, "Signature": h'6801ce78df1867c633fdc4b8a13fd7a5b81c93751ee9c9688e5b35a7f4e74756eefc4e6b22241f13c721fe9185b8a1fa2b41752ae5254f78bcf9025ed837990e'}
```


### Verify a block header

```bash
nacatgunma header verify \
  --header-file header.cbor
```

```console
Verified signature by did:key:z6Mkrpqbu1hJWGHCuTa9Z9SDaCp8VJa82dLYhhCWYiyye6Nr
```

```bash
echo $?
```

```console
0
```


### Export a block header as JSON

```bash
nacatgunma header export \
  --header-file header.cbor \
  --output-file header.json

json2yaml header.json
```

```console
Payload:
  Version: 1
  Accept:
  - /: QmXUuPakG4UHBs7BBQKS6o3vKeZJkRy7yTY4Pu1iGSr2PY
  Reject: null
  Body:
    /: QmXUuPakG4UHBs7BBQKS6o3vKeZJkRy7yTY4Pu1iGSr2PY
  Schema: DAG-CBOR
  MediaType: application/cbor
Issuer: did:key:z6Mkrpqbu1hJWGHCuTa9Z9SDaCp8VJa82dLYhhCWYiyye6Nr
Signature: /Up9zvN/8Nf9uR0kcC5PSgLR+c4nYx2GL4XVKYR5hh3fsGsvwECGnw0E5kee9ZDj2oZlLy826bKcPQV+C4DxBw==
```


### Create a block body from RDF N-quads

```bash
nacatgunma body rdf \
  --rdf-file body.nq \
  --base-uri http://example.org/person \
  --body-file body.cbor
```

```console
bafyreiea2su23cm4nbfl3675m442gp5yo5qmghspjikeeeioudyls2jjtm
```

```bash
cbordump body.cbor
```

```console
{"@id": "g1", "@graph": [{"@id": "#1234", "name": "Alice", "knows": {"@id": "#5678"}}, {"@id": "#5678", "name": "Bob"}], "@context": {"name": "http://schema.org/name", "knows": "http://schema.org/knows"}}
```


### Export a block body as JSON

```bash
nacatgunma body export \
  --body-file body.cbor \
  --output-file body.json

json2yaml body.json
```

```console
'@context':
  knows: http://schema.org/knows
  name: http://schema.org/name
'@graph':
- '@id': '#1234'
  knows:
    '@id': '#5678'
  name: Alice
- '@id': '#5678'
  name: Bob
'@id': g1
```


### Store a block on IPFS

```bash
nacatgunma ipfs store \
  --key-file private.pem \
  --body-file body.cbor \
  --accept bafyreid32eo34hcksuilttqhoyhssoz36r5umtk6zallvgoxqlehsqltru
```

```console
bafyreib5fuk4qex34is3pt52ij4jddlnsevkys7jwa6v2lp2qrs2eoq5he
```


### Fetch a block from IPFS

```bash
nacatgunma ipfs fetch \
  --header-cid bafyreib5fuk4qex34is3pt52ij4jddlnsevkys7jwa6v2lp2qrs2eoq5he \
  --header-file h.cbor \
  --body-file b.cbor
```
