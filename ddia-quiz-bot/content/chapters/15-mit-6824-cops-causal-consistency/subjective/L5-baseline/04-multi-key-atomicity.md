---
id: cops-subjective-L5-004
type: subjective
level: L5
category: advanced
topic: cops
subtopic: multi-key-atomicity
estimated_time: 10-12 minutes
---

# question_title - Multi-Key Operations and Atomicity

## main_question - Core Question
"A banking application needs to transfer money: decrement account A, increment account B. Base COPS doesn't provide atomic multi-key updates. Describe the problems that arise and propose two different approaches to handle this on top of COPS. Compare their trade-offs."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Base COPS: each put is independent, no cross-key atomicity
- Problem: readers might see debit without credit (inconsistent state)
- Approach 1: Application-level two-phase commit (2PC) using COPS
- Approach 2: Use COPS-GT extension for multi-key transactions
- Trade-offs: coordination cost, latency, consistency guarantees, complexity

### expected_keywords
- atomicity, multi-key transaction, two-phase commit, COPS-GT, coordination, inconsistent state

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Alternative: compensation/reversal transactions
- Alternative: escrow/transfer record as third key (event sourcing)
- Read-after-write consistency for transaction originator
- Cross-shard coordination challenges
- Comparison with Spanner's serializable transactions

### bonus_keywords
- compensation, event sourcing, escrow, cross-shard, serializable isolation

## sample_excellent - Example Excellence
"**Problem:** With base COPS, put(A, balance-100) and put(B, balance+100) are separate operations. A reader might see A's debit before B's credit (via replication lag or causal ordering), violating invariant 'total money constant.' This is worse than lost updates—it's temporarily inconsistent state visible to queries.

**Approach 1 - Application 2PC:**
- Client writes transaction intent to TXN_LOG key
- Performs put(A) and put(B) with dependency on TXN_LOG
- Readers check TXN_LOG before trusting balances
- Trade-offs: No COPS changes needed, but application complexity high, reads must check TXN_LOG, and no isolation (concurrent txns can interfere).

**Approach 2 - COPS-GT:**
- Use COPS-GT extension for atomic multi-key writes (not just reads)
- Requires coordinator to ensure both puts visible atomically at all replicas
- Trade-offs: Better semantics (true atomicity), but requires extending COPS protocol, adds latency (multi-round), and reduces availability (coordinator becomes bottleneck).

**Better Alternative:** Event sourcing—create transfer record T:(A→B, $100) as single COPS put. A and B accounts are derived state (sum of transfers). Reads must compute balance by scanning transfers, but single-key atomicity is preserved."

## sample_acceptable - Minimum Acceptable
"Base COPS can't atomically update A and B. Readers see inconsistent state. Could use application-level 2PC with transaction log, or COPS-GT for multi-key transactions. 2PC is complex for app, COPS-GT adds latency. Could also use event sourcing with transfer records."

## common_mistakes - Watch Out For
- Not explaining the specific inconsistency problem (partial visibility)
- Proposing solutions without discussing trade-offs
- Missing that COPS-GT exists as an extension

## follow_up_excellent - Depth Probe
**Question**: "How would you handle cascading transfers (A→B→C) where each transfer depends on the previous?"
- **Looking for**: Dependency chains, transaction ordering, deadlock potential

## follow_up_partial - Guided Probe
**Question**: "Why doesn't COPS's causal consistency alone solve the atomicity problem?"
- **Hint embedded**: Causality vs atomicity

## follow_up_weak - Foundation Check
**Question**: "What does 'atomic multi-key update' mean?"
