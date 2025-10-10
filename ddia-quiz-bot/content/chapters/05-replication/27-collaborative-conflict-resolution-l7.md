---
id: ch05-collaborative-conflict-resolution-l7
day: 27
level: L7
tags: [replication, conflicts, CRDTs, collaboration, principal-engineer]
related_stories: []
---

# Conflict Resolution in Collaborative Systems

## question
You're building a collaborative document editing platform competing with Google Docs. Unlike Google's centralized approach, you want a decentralized architecture where documents can be edited offline and sync later, supporting peer-to-peer collaboration without requiring constant server connectivity. Design a conflict resolution system that provides good user experience while maintaining document integrity. Consider rich text formatting, embedded images, real-time collaboration when online, and offline editing capabilities.

## expected_concepts
- Operational transformation vs CRDTs trade-offs
- Hybrid online/offline architecture
- Intent preservation in conflict resolution
- Vector clocks and causal ordering
- Merkle trees for efficient sync
- Semantic vs syntactic conflicts
- User experience during conflicts

## answer
Core Architecture: Use a hybrid CRDT approach - the document is a tree of CRDTs where paragraphs are Last-Write-Win registers identified by unique IDs, text within paragraphs uses RGA (Replicated Growable Array) for character-level editing, and formatting uses multi-value registers with user-based resolution.

Conflict Resolution Layers: (1) Structural conflicts (paragraph reordering, deletion) - use tree CRDTs with tombstones, preserving both orders until explicitly resolved. (2) Textual conflicts - RGA ensures convergence, but highlight areas with concurrent edits for review. (3) Semantic conflicts (two users changing a date to different values) - detect using domain-aware heuristics, surface to users for resolution.

Offline/Online Modes: When online, use operational transformation for immediate convergence with server acting as arbiter. When offline, accumulate CRDT operations. On reconnection, perform three-way merge: local changes, remote changes, and last known common state. Use Merkle trees to efficiently identify divergent sections.

User Experience: Implement "conflict-free replicated UI" - show local changes immediately, indicate sync status per section, highlight unresolved conflicts in margin. Provide "conflict resolution mode" where users can see all versions and choose resolution. Maintain "blame view" showing who made what changes when.

Intent Preservation: Record high-level operations (e.g., "move paragraph" not just delete+insert). Use operational transformation to rebase operations when possible. Implement "semantic locks" - temporary exclusive access for structural changes.

Key Innovation: Introduce "confidence scoring" for automatic resolution - higher confidence for older, stable text, lower for recently edited sections. Let users set "trust levels" for collaborators, influencing automatic resolution.

## hook
How do you handle the case where users make conflicting changes to the same image embedded in the document?

## follow_up
Your platform becomes successful, and enterprise customers want to use it for legal documents where every change must be auditable, reversible, and legally admissible. However, your CRDT-based system makes it difficult to establish a canonical "version history" since different nodes might see operations in different orders. How do you adapt your architecture to meet these regulatory requirements without losing the benefits of decentralized collaboration?

## follow_up_answer
Implement a "legal ledger" layer on top of the CRDT system: (1) Authoritative Event Log: While CRDTs handle convergence, maintain a separate cryptographically-signed event log using blockchain-inspired hash chaining. Each operation includes: timestamp, author signature, hash of previous operation, and witness signatures from online peers. (2) Periodic Checkpointing: Every N operations or T time, generate a "legal checkpoint" - a deterministic serialization of document state with all collaborators' signatures. These become the official versions for legal purposes.

Dual-Mode Operation: Introduce "regulated mode" for legal documents - requires minimum quorum of witnesses for operations, enforces stricter ordering via temporary coordination server, and prohibits offline editing. Regular mode remains fully decentralized. Documents can graduate from regular to regulated mode but not vice versa.

Audit Trail Construction: Build "causal history reconstructor" that can deterministically order operations post-hoc using vector clocks and witness attestations. Generate "legal diff view" showing court-admissible change history. Implement "chain of custody" tracking - who had access when, what operations were possible.

Compliance Infrastructure: Deploy "compliance nodes" that act as trusted witnesses, maintained by legal entity. These nodes: continuously validate operation ordering, detect and flag suspicious patterns (backdated edits, signature mismatches), and generate compliance certificates for document states.

Critical Trade-off: Accept that regulated mode sacrifices some distributed properties for legal compliance. The innovation is making this modal - documents that need compliance get it, while others maintain full distributed benefits. This selective approach keeps the platform viable for both casual and regulated use cases.
