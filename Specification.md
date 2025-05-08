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
}
```


## Blocks

Let $\mathcal{C}$ be the set of content identifiers (e.g., CIDs), and let $\mathcal{P}$ be the set of participants with associated public keys.

Each block is a tuple:

$$
B = (\text{id}, \mathsf{accept}(B), \mathsf{reject}(B), \mathsf{body}(B), \mathsf{schema}(B), \mathsf{media}(B), \mathsf{issuer}(B), \mathsf{sig}(B))
$$

where:
- $\text{id} \in \mathcal{C}$: the content identifier of the block itself.
- $\mathsf{accept}(B) \subset \mathcal{C}$: set of parent block CIDs that this block explicitly extends.
- $\mathsf{reject}(B) \subset \mathcal{C}$: set of block CIDs that this block rejects.
- $\mathsf{body}(B) \in \mathcal{C}$: CID of the block’s body (arbitrary content).
- $\mathsf{schema}(B) \in \mathbb{S}$: schema URI string.
- $\mathsf{media}(B) \in \mathbb{M}$: media type string.
- $\mathsf{issuer}(B) \in \mathcal{P}$: public key of the block’s issuer.
- $\mathsf{sig}(B)$$: digital signature over the payload fields, created using the issuer’s secret key.


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
    T_u : \mathcal{P} \to [0,1]
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

The following non-normative pseudo-code sketches the subjective DAG traversal and rejection.

```lean
structure Block where
  cid     : Cid          -- self
  accept  : List Cid     -- parents
  reject  : List Cid     -- rejected roots

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
          match blocks.find? cid with
          | none => visit rest (visited.insert cid) visible context
          | some block =>
            -- gather inherited rejections from children
            let inheritedRej :=
              blocks.fold (init := ∅) fun acc child =>
                if child.accept.contains cid then
                  acc ∪ context.findD child.cid ∅
                else acc

            let fullRej := inheritedRej ∪ block.reject.toHashSet
            let context' := context.insert cid fullRej

            if fullRej.contains cid then
              visit rest (visited.insert cid) visible context'
            else
              let visited' := visited.insert cid
              let visible' := visible.insert cid
              let stack' := block.accept ++ rest
              visit stack' visited' visible' context'

  visit trustedTips ∅ ∅ HashMap.empty
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


### Summary

| Concept   | Formal                                                          |
|-----------|-----------------------------------------------------------------|
| Block     | Tuple of accepted/rejected parents, body CID, issuer, signature |
| DAG       | Induced by accepted parent links                                |
| Trust     | Local function $T_u : \mathcal{P} \to [0,1]$                    |
| View      | Blocks visible from tips via non-rejected paths                 |
| Rejection | Explicitly prunes traversal unless overridden later             |
| State     | Application-defined interpretation of visible block bodies      |
