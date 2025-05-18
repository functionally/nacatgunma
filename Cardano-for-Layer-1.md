# Using Cardand for Layer 1

Cardano can be used as a Layer 1 by Nacatgunma. Users can publicize their new tips via Cardano smart-contract transactions.


## Plutus smart contract

The Plutus smart contract for tracking the tip is written in [Pluto](https://github.com/Plutonomicon/pluto) and compiled to [Untyped Plutus Core (UPLC)](https://plutonomicon.github.io/plutonomicon/uplc).

- [Pluto source code](onchain/script-0.pluto)
- [UPLC in CBOR](onchain/script-0.cbor)
- [Plutus text envelope](onchain/script-0.plutus)

The contract has the following addresses:

- Mainnet: [addr1w97pkwtxuqh3x8edfjhsndd5f8ymyu2gnnwq4wx7jytt7lcnu3lv5](https://cardanoscan.io/address/717c1b3966e02f131f2d4caf09b5b449c9b271489cdc0ab8de9116bf7f)
- Preprod: [addr\_test1wp7pkwtxuqh3x8edfjhsndd5f8ymyu2gnnwq4wx7jytt7lcg59rr3](https://preprod.cardanoscan.io/address/707c1b3966e02f131f2d4caf09b5b449c9b271489cdc0ab8de9116bf7f)


## Datum

The datum specifies the hash of the credential that is authorized to spend the UTxO at the contract address. This credential may be a public key, an simple script, or a Plutus script. The datum also specifies the IPFS CID of the block at the tip.

```yaml
list:
- constructor: 1                                                                   # 0 = public key; 1 = script
  fields:                                                                          #
  - bytes: 30135f08305143796de4276083cc54e47fbcafb176df6b58ab309446                # hash of the credential (public key or script)
- bytes: 01711220cd65fa79bf59e0f3985678e872fc5bfac509004be9c8ada824a927e20548f968  # IPFS CID for the block header of the tip
```


## Redeemer

The redeemer is arbitrary. For concreteness, we use the integer of the Nacatgunma metadata key as the redeemer.

```yaml
int: 58312
```


## Metadata

For convenience only, the transaction metadata may provide a reference to this repository and repeat the IPFS CID for the block header at the tip in human-readable format.

```yaml
'58312':
  blockchain: https://github.com/functionally/nacatgunma
  header:
    ipfs: bafyreignmx5htp2z4dzzqvty5bzpyw72yueqas7jzcw2qjfje7rakshzna
```


## Transaction

For easiest on-chain discoverability of the CID of the tip, include the datum as *inline datum* in the transaction. This makes it possible to retrieve that CID and the datum directly from a UTxO query on the node; otherwise, one must use a chain follower or chain database to discover the details of the tip or find the latest tips.

See [onchain/mainnet/](onchain/mainnet/) and [onchain/testnet/](onchain/testnet/) for model transactions.

1. [d4892fc7f2f50a493f03d43dea6158e7a3d3aafdf6c036704043284a1a54bd68](https://cardanoscan.io/transaction/d4892fc7f2f50a493f03d43dea6158e7a3d3aafdf6c036704043284a1a54bd68?tab=utxo)
2. [1112d0521791e6e1439a66c7c7055b1f5fa56247e6b7e1113094e767dd413d8f](https://cardanoscan.io/transaction/1112d0521791e6e1439a66c7c7055b1f5fa56247e6b7e1113094e767dd413d8f?tab=utxo)
