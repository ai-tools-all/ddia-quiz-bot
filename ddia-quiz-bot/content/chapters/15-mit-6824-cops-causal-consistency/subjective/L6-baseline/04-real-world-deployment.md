---
id: cops-subjective-L6-004
type: subjective
level: L6
category: advanced
topic: cops
subtopic: real-world-deployment
estimated_time: 12-15 minutes
---

# question_title - Real-World Deployment Decision

## main_question - Core Question
"Your team is building a real-time collaborative platform (think Figma, Notion, or Google Docs) with millions of users globally. You're evaluating COPS as the consistency model. Conduct a thorough analysis: (1) What workload characteristics make COPS a good or bad fit? (2) What specific features of the application would work well with COPS, and what wouldn't? (3) What alternative consistency models would you compare against, and under what conditions would you choose each? (4) If you choose COPS, what extensions or modifications would you need?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Workload analysis: read/write ratio, access patterns, causal dependencies in app
- COPS strengths: local writes, cross-object causality (doc → folder → permissions)
- COPS weaknesses: no atomic multi-key updates, LWW loses concurrent edits, cascading delays
- Need for CRDTs or OT on top of COPS for collaborative editing
- Alternatives: Spanner (strong consistency, higher latency), CRDTs alone (AP, weaker ordering)
- Extensions needed: COPS-GT for consistent snapshots, multi-key atomicity

### expected_keywords
- workload characteristics, causal dependencies, CRDT, operational transformation, consistency model comparison, extensions

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Hybrid approach: COPS for metadata, CRDTs for document content
- Presence and cursors (real-time updates) don't need COPS
- Access control propagation critical (permission revoke must be immediate)
- Offline editing and sync requirements
- Cost analysis: infrastructure, development, operational
- Migration path from existing system

### bonus_keywords
- hybrid consistency, presence protocol, access control, offline-first, total cost of ownership, migration strategy

## sample_excellent - Example Excellence
"**Workload Analysis:**
- **Reads:** 90% of operations (viewing docs, presence updates)—COPS local reads benefit here
- **Writes:** 10% but critical (edits, permissions, sharing)—mix of causal and concurrent
- **Access patterns:** Users work on few docs but read many (catalogs, templates)—partial replication opportunity

**COPS Fit Assessment:**

**Good Fit:**
- Document metadata operations: create doc → add to folder → set permissions (perfect causal chain)
- Comment threads: reply → notification (causal dependency)
- Version history: immutable operations (append-only, no conflicts)

**Poor Fit:**
- Real-time collaborative editing: concurrent character insertions need CRDTs/OT, not LWW
- Atomic operations: 'move doc from Folder A to Folder B' (needs multi-key atomicity)
- Presence/cursors: need immediate propagation, COPS dependency chains too slow

**Alternative Comparison:**

1. **Spanner:** Strong consistency, but 50-100ms cross-DC write latency unacceptable for typing experience. Use for: billing, subscription state.

2. **CRDTs (AP):** No coordination, great for text editing, but weaker ordering (can't guarantee 'see permission before doc'). Use for: document content only.

3. **Hybrid Architecture (Chosen):**
   - **COPS for metadata layer:** Folders, permissions, sharing, version pointers
   - **CRDTs for content layer:** Text editing (Yjs/Automerge), presence (Conflict-free Replicated Data Types)
   - **Spanner for critical state:** Subscriptions, billing, access logs

**COPS Extensions Needed:**
- **COPS-GT:** Consistent multi-key reads (get doc + permissions atomically)
- **Modified LWW:** For metadata, use app-specific tie-breaker (user ID) instead of pure timestamp
- **Dependency compression:** Batch dependencies for large documents (1000s of edits)
- **Timeout degradation:** After 5s dependency wait, show stale data with warning banner

**Operational Requirements:**
- Monitoring: Track P99 dependency wait time (SLA: <500ms)
- Causality debugger: Visualize dependency chains for troubleshooting
- Rollback capability: If COPS causes production issues, fallback to eventual consistency

**Cost Analysis:**
- Development: 6 months custom COPS implementation + tooling
- Infrastructure: 20% higher storage (version metadata), 30% higher network (dependency propagation)
- Benefit: Better UX (fewer anomalies), 50% reduction in conflict resolution support tickets"

## sample_acceptable - Minimum Acceptable
"Collaborative platform has causal dependencies (doc creation → folder → permissions) which COPS handles well. But concurrent editing needs CRDTs on top since COPS's LWW loses edits. Compare with Spanner (too slow for typing) and pure CRDTs (no causal guarantees for metadata). Hybrid approach: COPS for metadata, CRDTs for content. Need COPS-GT for multi-key reads and dependency compression for performance."

## common_mistakes - Watch Out For
- Not analyzing specific workload characteristics
- Missing that COPS alone is insufficient (needs CRDTs for editing)
- Not comparing multiple alternatives
- Ignoring cost/complexity of implementation

## follow_up_excellent - Depth Probe
**Question**: "How would you design the migration path from an existing eventual consistency system to COPS without downtime?"
- **Looking for**: Dual-write strategy, shadow mode testing, gradual rollout, rollback plan

## follow_up_partial - Guided Probe
**Question**: "Why is COPS not sufficient for collaborative text editing?"
- **Hint embedded**: Concurrent character insertions

## follow_up_weak - Foundation Check
**Question**: "What is the main benefit of COPS for the 'create doc → add to folder' workflow?"
