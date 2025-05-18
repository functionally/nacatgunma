# Using Cardand for Layer 1

Cardano can be used as a Layer 1 by Nacatgunma. Users can publicize their new tips via Cardano smart-contract transactions.


## Plutus smart contract

The Plutus smart contract for tracking the tip is written in [Pluto](https://github.com/Plutonomicon/pluto) and compiled to [Untyped Plutus Core (UPLC)](https://plutonomicon.github.io/plutonomicon/uplc).

- [Pluto source code](onchain/script-0.pluto)
- [UPLC in CBOR](onchain/script-0.cbor)
- [Plutus text envelope](onchain/script-0.plutus)

The contract has the following addresses:

- Mainnet: [addr1w8lyu0uj30gyytukg25ynfypvqlw7tt4duuu7lqd09qrnugm34xp8](https://cardanoscan.io/address/71fe4e3f928bd0422f9642a849a481603eef2d756f39cf7c0d794039f1)
- Preprod: [addr\_test1wrlyu0uj30gyytukg25ynfypvqlw7tt4duuu7lqd09qrnugqep6wz](https://preprod.cardanoscan.io/address/https://preprod.cardanoscan.io/address/70fe4e3f928bd0422f9642a849a481603eef2d756f39cf7c0d794039f1)


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

1. [c4fb42bdca24628f96d3268c71253ebb9565453bf29de9b4ace70521a9bb9959](https://preprod.cardanoscan.io/transaction/c4fb42bdca24628f96d3268c71253ebb9565453bf29de9b4ace70521a9bb9959?tab=utxo)
2. [ae1fa36e2f82b9d031169b8aaa2f139dcd62d93820122c3ff37cd0ee4386ef64](https://preprod.cardanoscan.io/transaction/ae1fa36e2f82b9d031169b8aaa2f139dcd62d93820122c3ff37cd0ee4386ef64?tab=utxo)
