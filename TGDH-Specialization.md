# Extension: Tree-based group Diffie-Hellman (TGDH) encryption

## TGDH encrypted payload

The TGDH scheme described here allows the decentralized addition and removal of keys for a shared secret that is recomputed upon each addition or removal.

Each node in the a TGDH key contains the following:

- BLS12-381 G1 public key
- An optional BLS12-381 Fr private key.
- Optional left and right child nodes.

The public key is simply the exponentiation of the private key, with G1 as the base.

The private key of a leaf is simply a randomly chosen member of the Fr group. The private key of a node is computed as the SHA256-based HKDF Fr hash of the product of the private key of one child and the public key of the other child. The string `nacatgunma-tgdh-bls12381g1` is used as the information field for HKDF.

These keys can be used to derive a symmetric key for AES256-GCM encryption, but a unique salt must be used each time in the symmetric-key derivation.

The reference implementation resides in [tgdh/](./tgdh/).

## Command-line tool

The `nacatgunma body tgdh` subcommands implement the TGDH specification described above.

```console
$ nacatgunma body tgdh --help

NAME:
   nacatgunma body tgdh - Tree-based group DH (BLS12-381) management subcommands

USAGE:
   nacatgunma body tgdh [command options]

COMMANDS:
   decrypt   Decrypt a file using a TGDH private key.
   encrypt   Encrypt a file using a TGDH private key.
   generate  Generate a TGDH private key.
   join      Join two TGDH keys into an aggregate TGDH key, where at least one of the keys is private.
   public    Strip private key information from a TGDH key.
   private   Apply a private TGHD key to a public one, deriving the private root.
   remove    Remove a TGDH keys from an aggregate TGDH key, where at least one of the root keys is private.
   help, h   Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

### Generate a key

```console
$ nacatgunma body tgdh generate --private-file A.pri

did:key:z3tEEhukKAKcFEGQn3fCFngDPForeCYJA2kSnhyvJmzDg4VqgxyKBQitb3LKjxLAzmdH3w

$ json2yaml A.pri

private: 126822990fac11644ded34c7e887a9c239b627f34e082c68d3e03e25b59f3b16
public: 82ecf95216b95e634bfda85171d5a10abadfc9cfb5226133c3dbe6496f11d1aebbb586ead61d382de2be22eebb646aaa
```

### Extract the public key

The public key can safely be posted on the blockchain so that other parties can include it in group encryption.

```console
$ nacatgunma body tgdh public --private-file A.pri --public-file A.pub

did:key:z3tEEhukKAKcFEGQn3fCFngDPForeCYJA2kSnhyvJmzDg4VqgxyKBQitb3LKjxLAzmdH3w

$ json2yaml A.pub

public: 82ecf95216b95e634bfda85171d5a10abadfc9cfb5226133c3dbe6496f11d1aebbb586ead61d382de2be22eebb646aaa
```

### Aggregate two keys

One of the keys joined into an aggregate must be private, so that the public key of the aggregate can be derived.

```console
$ nacatgunma body tgdh join --left-file A.pri --right-file B.pub --private-file AB.pri

did:key:z3tEFnWHMrZzbc4MqYGNF2hrH2LeaEN9FgxVsizEao8QgMCJcDQJCgYMEGLyfMBtYYFNGR

$ json2yaml AB.pri

private: 46560d7d32f9149b214d5764e1679a58e4750605780ad4cf6c8a3b7c7811a1c9
public: a058b91fe134956663aa5f5775b486622caa2f8fc00f3b71c2c0c378e6677808e307470e74bfad8306ae84d83c5098b2
left:
  private: 126822990fac11644ded34c7e887a9c239b627f34e082c68d3e03e25b59f3b16
  public: 82ecf95216b95e634bfda85171d5a10abadfc9cfb5226133c3dbe6496f11d1aebbb586ead61d382de2be22eebb646aaa
right:
  public: 930f91677ed24327207f2bf154dac458208ab664f80549305663e1c412fb2d692c45991aae4446a099fdf586daf39322
```

The public aggregate key can be shared on the blockchain so that other parties can include it in their key derivations.

```console
$ nacatgunma body tgdh public --private-file AB.pri --public-file AB.pub
did:key:z3tEFnWHMrZzbc4MqYGNF2hrH2LeaEN9FgxVsizEao8QgMCJcDQJCgYMEGLyfMBtYYFNGR


$ json2yaml AB.pub

public: a058b91fe134956663aa5f5775b486622caa2f8fc00f3b71c2c0c378e6677808e307470e74bfad8306ae84d83c5098b2
left:
  public: 82ecf95216b95e634bfda85171d5a10abadfc9cfb5226133c3dbe6496f11d1aebbb586ead61d382de2be22eebb646aaa
right:
  public: 930f91677ed24327207f2bf154dac458208ab664f80549305663e1c412fb2d692c45991aae4446a099fdf586daf39322
```

### Recover the private key of an aggregate

If one of the leaf private keys of the aggregate is known, the root private key can be derived.

```console
$ nacatgunma body tgdh private --private-file A.pri --public-file AB.pub --root-file AB.pri

did:key:z3tEFnWHMrZzbc4MqYGNF2hrH2LeaEN9FgxVsizEao8QgMCJcDQJCgYMEGLyfMBtYYFNGR

$ json2yaml AB.pri

private: 46560d7d32f9149b214d5764e1679a58e4750605780ad4cf6c8a3b7c7811a1c9
public: a058b91fe134956663aa5f5775b486622caa2f8fc00f3b71c2c0c378e6677808e307470e74bfad8306ae84d83c5098b2
left:
  private: 126822990fac11644ded34c7e887a9c239b627f34e082c68d3e03e25b59f3b16
  public: 82ecf95216b95e634bfda85171d5a10abadfc9cfb5226133c3dbe6496f11d1aebbb586ead61d382de2be22eebb646aaa
right:
  public: 930f91677ed24327207f2bf154dac458208ab664f80549305663e1c412fb2d692c45991aae4446a099fdf586daf39322
```

### Remove a party from an aggregate

At least one private key must be known to remove a party from an aggregate key.

```console
$ nacatgunma body tgdh remove --leaf-file B.pub --root-file ABCD.pri --private-file ACD.pri

did:key:z3tEGJJpro3jZseNAFjT8dH8UW3KXY2bLwmayjhUNau3EehzdWhtzMhohgSZvuwdWB4qxi

$ json2yaml ACD.pri

private: 54214cb4eee507a2ac6862d54eade6c7b79efd768af2d139129d549c28ee8ed4
public: ae5aaecbf1315d09b24801255190709a53bf170e10518dadf7bb2109a94092ecc8d5d7258c31aa3c32546c8e6303ed97
left:
  private: 4fc6442ead18c3122bc3351bedf69b115c1ec5b07e2ff7ce459d23adb82b5687
  public: ac879c4bd3f97720e0a7fd84980b3aa0f241b27bcab7c1f83399adc41704d7019ef5e90339c7f65ffa358e19f35cc986
  left:
    public: 86b53d5f07478c943ed255b1c0808d81a41ade43e88caf5e9ea8a423e683a8ed7e2d77af7c6a48ffba45ec6b18bda8bd
  right:
    private: 657090a5812aec229e621fd0bfdc01313ab40c660638616d62de356cfe9f073f
    public: b27be6d6d50c6c43758d7f4a5b413311a7e9dafaf57c562b6937e4b9f650f0c0ed008dcc76021d30751b4a8436bbf0e7
right:
  public: 82ecf95216b95e634bfda85171d5a10abadfc9cfb5226133c3dbe6496f11d1aebbb586ead61d382de2be22eebb646aaa
```

### Encrypt to an aggregate key

Each encryption uses the aggregate private key and randomly chosen salt to derive the AES256 symmetric key.

```console
$ echo "Hello, TGDH!" > hello.txt

$ nacatgunma body tgdh encrypt --private-file ACD.pri --plaintext-file hello.txt --content-type "text/plain" --jwe-file hello.jwe --jwk-file hello.jwk

$ json2yaml hello.jwe

ciphertext: MnU5R1izHU_KlaKZ3w
header:
  alg: dir
iv: pmUodiNTGFUhj8NV
protected: eyJhbGciOiJkaXIiLCJjdHkiOiJ0ZXh0L3BsYWluIiwiZW5jIjoiQTI1NkdDTSIsImtpZCI6ImRpZDprZXk6ejN0RUdKSnBybzNqWnNlTkFGalQ4ZEg4VVczS1hZMmJMd21heWpoVU5hdTNFZWh6ZFdodHpNaG9oZ1NadnV3ZFdCNHF4aSIsInAycyI6Ijc2ajlfeUxraUtYQTdBak5wVnM1Wkk3N29VUFdZQUN6UERwdHFhUWc0c0UiLCJ0eXAiOiJuYWNhdGd1bm1hLXRnZGgrc2FsdCJ9
tag: iLBi6HxnT8E7C-RiqB9UuA

$ jq -r .protected hello.jwe | basenc --decode --base64 | json2yaml

alg: dir
cty: text/plain
enc: A256GCM
kid: did:key:z3tEGJJpro3jZseNAFjT8dH8UW3KXY2bLwmayjhUNau3EehzdWhtzMhohgSZvuwdWB4qxi
p2s: 76j9_yLkiKXA7AjNpVs5ZI77oUPWYACzPDptqaQg4sE
typ: nacatgunma-tgdh+salt

$ json2yaml hello.jwk
alg: dir
k: fCFeQob5Vh0Zl9JvoLZ4KOSKnhPfi7i9QLEGuEqTi-4
key_ops:
- encrypt
- decrypt
kid: did:key:z3tEGJJpro3jZseNAFjT8dH8UW3KXY2bLwmayjhUNau3EehzdWhtzMhohgSZvuwdWB4qxi
kty: oct
use: enc
```
Supplying the content type or the JWK filename are optional. The JWK file is useful because it can be used with third-party tools to decrypt the JWE file.

### Decrypt a message

```console
$ nacatgunma body tgdh decrypt --private-file ACD.pri --jwe-file hello.jwe --plaintext-file hello.txt --headers-file headers.json

$ cat hello.txt

Hello, TGDH!

$ json2yaml headers.json

alg: dir
cty: text/plain
enc: A256GCM
kid: did:key:z3tEGJJpro3jZseNAFjT8dH8UW3KXY2bLwmayjhUNau3EehzdWhtzMhohgSZvuwdWB4qxi
p2s: 76j9_yLkiKXA7AjNpVs5ZI77oUPWYACzPDptqaQg4sE
typ: nacatgunma-tgdh+salt
```

Creating the headers file is optional.
