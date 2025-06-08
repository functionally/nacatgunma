# Formal Model of a Provenance-Aware Subjective DAG Ledger with Acceptance and Rejection


## Fundamental types


We define a **block** as a signed message that declares:

- The block’s **accepted parents** in the DAG.
- The set of **blocks rejected** by this block.
- A **CID of arbitrary body data**.
- Metadata to interpret the body (schema and media type).

Each block is composed of a payload and a signature, matching the following conceptual structure:

```go
type Header struct {
	Payload   Payload       // metadata for block body
	Issuer    string        // DID of the issuer
	Signature []byte        // signature by the issuer
}

type Payload struct {
	Version   int64
	Accept    []cid.Cid     // accepted parent blocks
	Reject    []cid.Cid     // rejected blocks
	Body      cid.Cid       // content-addressed arbitrary data
	SchemaUri string        // schema URI
	MediaType string        // MIME type
	Comment   string        // arbitrary commentary
}
```


## Blocks

Let $\mathcal{C}$ be the set of content identifiers (e.g., CIDs), and let $\mathcal{P}$ be the set of participants with associated public keys.

Each block is a tuple:

$$
B = (\text{id}, \mathsf{accept}(B), \mathsf{reject}(B), \mathsf{body}(B), \mathsf{schema}(B), \mathsf{media}(B), \mathsf{comment}(B), \mathsf{issuer}(B), \mathsf{sig}(B))
$$

where:
- $\text{id} \in \mathcal{C}$: the content identifier of the block itself.
- $\mathsf{accept}(B) \subset \mathcal{C}$: set of parent block CIDs that this block explicitly extends.
- $\mathsf{reject}(B) \subset \mathcal{C}$: set of block CIDs that this block rejects.
- $\mathsf{body}(B) \in \mathcal{C}$: CID of the block’s body (arbitrary content).
- $\mathsf{schema}(B) \in \mathbb{S}$: schema URI string.
- $\mathsf{media}(B) \in \mathbb{M}$: media type string.
- $\mathsf{comment}(B)$: comment string.
- $\mathsf{issuer}(B) \in \mathcal{P}$: DID of the block’s issuer.
- $\mathsf{sig}(B)$: digital signature over the payload fields, created using the issuer’s secret key.


## DAG Structure and validity

Let the **global DAG** be a graph $\mathcal{D}$ where nodes are block identifiers and edges are induced by $\mathsf{accept}(B)$.

A block $B$ is **valid** iff:

1. All parent CIDs in $\mathsf{accept}(B)$ exist in the DAG.
2. The signature $\mathsf{sig}(B)$ verifies under $\mathsf{issuer}(B)$ against a canonical serialization of the payload.
3. $B$ does not create a cycle when added to $\mathcal{D}$.


## Participant trust and views

Each participant $u \in \mathcal{P}$ maintains:

- A **trust function**:

$$
T\_u : \mathcal{P} \to [0,1]
$$

- A set of **trusted tips** $\mathcal{T}_u \subseteq \mathcal{C}$ used as starting points for state computation.


## Subjective traversal with rejection

For a participant $u$, the **visible subgraph** $\mathcal{D}_u \subseteq \mathcal{D}$ is constructed as follows:

1. Start from trusted tips $\mathcal{T}_u$.
2. Traverse backward along $\mathsf{accept}$ edges recursively.
3. Whenever a block $B$ is encountered:
   - If any ancestor block $A$ declares $B \in \mathsf{reject}(A)$, then:
     - Omit $B$ and all descendants reachable via $\mathsf{accept}$ edges **unless** they are also reachable via a separate path that does not pass through a rejecting ancestor.

This defines **rejection as a pruning operator** on traversal.

Formally, the **visible blocks** for participant $u$ are:

$$
\mathcal{V}_u = \{ B \in \mathcal{D} \mid B \text{ is reachable from } \mathcal{T}_u \text{ via non-rejected paths} \}
$$

A path is **non-rejected** if no intermediate block $A$ along the path has $\mathsf{reject}(A) \ni B'$ for any $B'$ also on the path.

The following non-normative pseudo-code sketches the subjective DAG traversal and rejection. See [Pruning.lean](Pruning.lean) for Lean4 code and examples of pruning rejected blocks. One can experiment with the pruning algorithm [here](https://live.lean-lang.org/#codez=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).

```lean
import Std.Data.HashMap
import Std.Data.HashSet

open Std

abbrev Cid := String

structure Block where
  cid     : Cid
  accept  : List Cid
  reject  : List Cid

abbrev BlockMap := HashMap Cid Block
abbrev RejectionContext := HashMap Cid (HashSet Cid)

partial def computeVisibleView
  (blocks : BlockMap)
  (trustedTips : List Cid)
  : HashSet Cid :=
  let rec visit
    (stack : List Cid)
    (visited : HashSet Cid)
    (visible : HashSet Cid)
    (context : RejectionContext)
    : HashSet Cid :=
      match stack with
      | [] => visible
      | cid :: rest =>
        if visited.contains cid then
          visit rest visited visible context
        else 
          match blocks.get? cid with
          | none => visit rest (visited.insert cid) visible context
          | some block =>
            let inheritedRej :=
              blocks.fold (init := ∅) fun acc _ child =>
                if child.accept.contains cid then
                  acc ∪ context.getD child.cid ∅
                else acc
            let fullRej := inheritedRej ∪ HashSet.ofList block.reject
            let context' := context.insert cid fullRej
            if fullRej.contains cid then
              visit rest (visited.insert cid) visible context'
            else
              let visited' := visited.insert cid
              let visible' := visible.insert cid
              let stack' := block.accept ++ rest
              visit stack' visited' visible' context'
  visit trustedTips ∅ ∅ (HashMap.emptyWithCapacity 10)
```


## State interpretation

Each participant $u$ uses their local visible DAG $\mathcal{V}_u$ to compute a **state view**, interpreted according to application-specific logic on:

- Block bodies $\mathsf{body}(B)$, using $\mathsf{schema}(B)$ and $\mathsf{media}(B)$
- Provenance (issuer of each block)
- Ordering (topological or heuristic over the DAG)

This may include reconstructing RDF graphs, applying CRDT logic, filtering by issuer trust, etc.


## Rejection semantics

- Rejection is **local** and **non-monotonic**: a descendant of a rejected block can later be re-accepted explicitly.
- Rejection is **non-destructive**: rejected blocks are still available in the DAG but **ignored** in a participant’s traversal.


## Summary

| Concept   | Formally                                                        |
|-----------|-----------------------------------------------------------------|
| Block     | Tuple of accepted/rejected parents, body CID, issuer, signature |
| DAG       | Induced by accepted parent links                                |
| Trust     | Local function $T\_u : \mathcal{P} \to [0,1]$                   |
| View      | Blocks visible from tips via non-rejected paths                 |
| Rejection | Explicitly prunes traversal unless overridden later             |
| State     | Application-defined interpretation of visible block bodies      |
