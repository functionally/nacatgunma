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
  --body QmXUuPakG4UHBs7BBQKS6o3vKeZJkRy7yTY4Pu1iGSr2PY \
  --accept QmXUuPakG4UHBs7BBQKS6o3vKeZJkRy7yTY4Pu1iGSr2PY \
  --accept QmXUuPakG4UHBs7BBQKS6o3vKeZJkRy7yTY4Pu1iGSr2PY \
  --header-file header.cbor
```

```console
QmTZaQgzYc8Jupb1zcmw4C66ign47n32pHyRL9hUreQvAN
```

```bash
cbordump header.cbor 
```

```console
{"Issuer": "did:key:z6Mkrpqbu1hJWGHCuTa9Z9SDaCp8VJa82dLYhhCWYiyye6Nr", "Payload": {"Body": 42(h'00122087d66da12bd9f0855fed372f7c05f2dae6fd1cec5e436dcfc26b93efc5413317'), "Accept": [42(h'00122087d66da12bd9f0855fed372f7c05f2dae6fd1cec5e436dcfc26b93efc5413317')], "Reject": [], "Schema": "DAG-CBOR", "Version": 1, "MediaType": "application/cbor"}, "Signature": h'fd4a7dcef37ff0d7fdb91d24702e4f4a02d1f9ce27631d862f85d5298479861ddfb06b2fc040869f0d04e6479ef590e3da86652f2f36e9b29c3d057e0b80f107'}
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
QmP8CHrZKFaGrG4HFBq8zzDfbpywEpU2ASPK4hSLFHtQPw
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
