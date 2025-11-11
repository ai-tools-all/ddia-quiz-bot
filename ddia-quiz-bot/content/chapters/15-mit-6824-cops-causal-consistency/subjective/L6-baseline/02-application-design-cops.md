---
id: cops-subjective-L6-002
type: subjective
level: L6
category: baseline
topic: cops
subtopic: application-design
estimated_time: 12-15 minutes
---

# question_title - Application Design with COPS Semantics

## main_question - Core Question
"You're building a collaborative document editing application (like Google Docs) on top of COPS. Analyze the consistency challenges for operations like: (1) creating a document, (2) adding it to a folder, (3) granting user permissions, and (4) users editing the document. What works well with COPS, what doesn't, and how would you architect the application layer to handle gaps?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Causal operations (create doc → add to folder → grant permission) work well with COPS
- COPS ensures users see doc before seeing it in folder, see permissions before accessing
- Concurrent edits to same document problematic: LWW loses edits
- Need operational transformation or CRDT at application layer for collaborative editing
- Cross-object dependencies (folder structure + permissions + content) tracked via context

### expected_keywords
- causal chain, cross-key dependencies, concurrent edits, application-layer merge, operational transformation

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Separating immutable data (doc creation) from mutable data (edits)
- Using COPS-GT (get transactions) for multi-key reads
- Designing conflict-free operations (append-only logs)
- Access control challenges: permission revocation propagation
- Optimistic UI with eventual convergence

### bonus_keywords
- immutability, COPS-GT, CRDT for text, append-only, access control propagation, optimistic rendering

## sample_excellent - Example Excellence
"COPS handles the causal chain well: create-doc → add-to-folder → grant-permission ensures users never see folder entry before doc exists, or access doc before permissions arrive. This prevents critical authorization failures. However, concurrent document edits fail with COPS's LWW—if two users edit simultaneously, one edit is lost.

Application layer architecture:
1. **Immutable document creation**: Store initial doc as immutable blob; COPS causal tracking prevents dangling folder references.
2. **Edit operations as append-only log**: Each edit is a new versioned entry (op1, op2, ...) rather than overwriting content. Use operational transformation (OT) or CRDT (Logoot, WOOT) to merge concurrent edits deterministically.
3. **Permission model**: Store permissions separately with version dependencies. Before rendering doc, check permissions are satisfied. Use COPS-GT for atomic multi-key reads (doc + permissions).
4. **Optimistic UI**: Show local edits immediately (COPS local write), sync causally across datacenters.
5. **Access revocation**: Add reverse dependencies—permission revokes must be visible before showing doc updates.

Gaps: COPS lacks transactions across multiple keys (use COPS-GT extension), and no built-in merge for conflicting edits (need app logic). Monitoring dependency wait times critical for user experience (long waits → staleness indicators)."

## sample_acceptable - Minimum Acceptable
"COPS ensures causal operations like create→add to folder→grant permission happen in order. Concurrent edits are a problem due to LWW. Application should use CRDTs or OT for merging edits, and possibly COPS-GT for multi-key consistency. Separate metadata (permissions) from content."

## common_mistakes - Watch Out For
- Not identifying which operations benefit from causal consistency
- Missing the collaborative editing challenge (concurrent edits)
- Not proposing concrete application-layer solutions (CRDTs/OT)
- Ignoring permission/access control propagation issues

## follow_up_excellent - Depth Probe
**Question**: "How would you handle a scenario where a user's permission is revoked, but they have unsync'd local edits in flight?"
- **Looking for**: Reverse causality, conflict resolution, authorization checks at commit time

## follow_up_partial - Guided Probe
**Question**: "Why is LWW insufficient for collaborative text editing?"
- **Hint embedded**: Multiple users typing simultaneously

## follow_up_weak - Foundation Check
**Question**: "What benefit does COPS provide for the create-doc → add-to-folder operation?"
