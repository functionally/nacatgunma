# Extension: Specialization to Signed RDF Graphs


## RDF graph payload

Assume that the content identifier $\mathsf{body}(B) \in \mathcal{C}$ references a concrete RDF graph: $\mathsf{G}(B) \subseteq \mathcal{T}$ where $\mathcal{T}$ is the set of RDF triples of the form $(s, p, o)$.

We interpret $\mathsf{G}(B)$ as the **payload graph** authored by $\mathsf{issuer}(B)$.


## Reification

We define a reification function $\mathsf{reify}(B)$ that maps the signed graph to a **named graph** or **quads** with provenance.

For each triple $t \in \mathsf{G}(B)$, define the reified form:

$$
\mathsf{reify}(B) = \{ (s, p, o, g\_B) \mid (s, p, o) \in \mathsf{G}(B) \}
$$

where $g\_B \in \mathcal{U}$ is a unique graph name associated with the block $B$, such as:

$$
g\_B := \texttt{did:} \mathsf{issuer}(B) \\# \mathsf{id}(B)
$$

Alternatively, if using RDF-star or PROV-O, the reification can instead be:

- A named graph or blank node denoting the assertion
- Annotated with:
  - $\texttt{prov:wasAttributedTo} \ \mathsf{issuer}(B)$
  - $\texttt{prov:wasDerivedFrom} \ \mathsf{id}(B)$

The exact form depends on the RDF dialect used, but semantically the payload is treated as **authored assertions by the signer**.


## State computation

Given a participant $u$, let $\mathcal{V}_u$ be the set of blocks **visible** under subjective traversal (see prior section with accept/reject logic).

Define the participant’s **assembled state graph** $\mathsf{S}_u$ as:

$$
\mathsf{S}_u := \cup_{B \in \mathcal{V}_u} \mathsf{reify}(B)
$$

That is, the total RDF graph is the union of all reified block payloads visible under $u$'s trust and traversal policy.

If RDF-star is used, this could equivalently be:

$$
\mathsf{S}_u := \cup_{B \in \mathcal{V}_u} \left\{ \ll t \gg \ \mathsf{prov:wasAttributedTo} \ \mathsf{issuer}(B) \right\}
\quad \text{for each } t \in \mathsf{G}(B)
$$


## Conflict handling

If blocks in $\mathcal{V}_u$ make **conflicting assertions** (e.g., two blocks claim different values for the same triple), this is **not a protocol error** — the provenance-aware model allows:

- Leaving conflicting assertions coexisting with different graph names
- Using downstream logic to resolve (e.g., trust-weighted resolution, query filtering)

Thus, **state** is:
- **Provenance-rich**: all facts have authorship
- **Non-authoritative**: clients determine what to trust or prioritize


## Summary

| Concept             | Description                                                                |
|---------------------|----------------------------------------------------------------------------|
| $\mathsf{G}(B)$     | RDF graph payload of block $B$                                             |
| $\mathsf{reify}(B)$ | Provenance-aware named graph or RDF-star reification                       |
| $\mathsf{S}_u$      | State for participant $u$: union of all reified payloads of visible blocks |
| Conflicts           | Allowed; resolution is subjective and external                             |
