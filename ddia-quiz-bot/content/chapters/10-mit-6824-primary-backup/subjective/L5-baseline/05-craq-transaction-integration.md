---
id: craq-subjective-L5-005
type: subjective
level: L5
category: baseline
topic: craq
subtopic: transaction-integration
estimated_time: 8-10 minutes
---

# question_title - Integrating CRAQ with Distributed Transactions

## main_question - Core Question
"You need to coordinate a CRAQ write with a payment service running two-phase commit (2PC). Design the interaction so that payments and CRAQ state stay consistent, referencing DDIA's transactions chapter. Discuss where to place CRAQ tail commit within the 2PC phases." 

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Prepare Phase Alignment**: CRAQ write must reach dirty state before prepare completes
- **Commit Phase Trigger**: Tail acknowledgment signals safe commit
- **Failure Handling**: If CRAQ can't reach tail, 2PC must abort to keep systems aligned
- **DDIA Link**: Applying 2PC with replicated logs/recoverable participants

### expected_keywords
- Primary keywords: two-phase commit, prepare, commit, tail acknowledgment
- Technical terms: XA transaction, coordinator, participant, timeout

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Idempotency and Retries**: CRAQ dedupe to survive coordinator retries
- **Timeout Strategy**: Ensure participants don't block indefinitely
- **Compensation**: Saga fallback if 2PC unavailable
- **Operational Impact**: Rolling upgrades, coordinator failover

### bonus_keywords
- Implementation: redo log, commit index, transaction coordinator, canonical ordering
- Scenarios: payment success but CRAQ fail, partial commit, network partition
- Trade-offs: latency vs atomicity, distributed lock holding

## sample_excellent - Example Excellence
"During prepare, the coordinator asks CRAQ to stage the update. The CRAQ head logs the write, propagates it to the tail, but does not flip the clean flag or acknowledge externally. Instead, it returns a 'prepared' token once the update is durable and dirty across the chain. When the coordinator issues commit, CRAQ waits for the tail acknowledgment, marks the entry clean, and sends final success. If the tail can't be reached, CRAQ returns an abort, so the coordinator instructs the payment service to roll back, following DDIA's 2PC guidance on atomic commit with replicated participants." 

## sample_acceptable - Minimum Acceptable
"Use 2PC: CRAQ participates by logging the write during prepare and only making it clean (tail ack + client visible) during commit. If the tail can't confirm, CRAQ forces the transaction to abort so payments and state stay consistent." 

## common_mistakes - Watch Out For
- Letting CRAQ return success during prepare
- Forgetting to mention tail acknowledgment gating commit
- Ignoring coordinator-driven retries
- Not referencing DDIA's 2PC patterns

## follow_up_excellent - Depth Probe
**Question**: "How would you recover if the coordinator crashes after CRAQ prepared but before commit?"
- **Looking for**: In-doubt transactions, log inspection, waiting for coordinator, timeouts to roll back
- **Red flags**: Automatically committing without coordinator direction

## follow_up_partial - Guided Probe  
**Question**: "What metadata must CRAQ persist so it can complete the transaction after restart?"
- **Hint embedded**: Prepared entries with coordinator ID, request ID
- **Concept testing**: Durability requirements

## follow_up_weak - Foundation Check
**Question**: "Why do both parties sign a contract before money changes hands?"
- **Simplification**: Prepare vs commit analogy
- **Building block**: Agreement before final action

## bar_raiser_question - L5→L6 Challenge
"If you switch from 2PC to a saga, outline compensating actions for failed CRAQ commits while guaranteeing payment reversals."

### bar_raiser_concepts
- Sagas vs 2PC trade-offs
- Compensation sequencing
- Distributed transaction resilience
- Business-level consistency

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 4-5 min answer + 4-5 min discussion
- **Common next topics**: Saga design, coordinator durability, failure testing
