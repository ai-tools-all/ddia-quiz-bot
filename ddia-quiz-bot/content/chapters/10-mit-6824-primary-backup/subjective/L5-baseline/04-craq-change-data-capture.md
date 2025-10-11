---
id: craq-subjective-L5-004
type: subjective
level: L5
category: baseline
topic: craq
subtopic: change-data-capture
estimated_time: 8-10 minutes
---

# question_title - CRAQ Change Data Capture Strategy

## main_question - Core Question
"Design a change data capture (CDC) pipeline that streams CRAQ writes into Apache Kafka while maintaining exactly-once semantics. Reference DDIA's streaming chapter and explain how CRAQ's tail acknowledgments factor into the design." 

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Tail-Triggered Capture**: Emit CDC events only after tail commit to avoid duplicates
- **Idempotent Kafka Producers**: Use sequence numbers to match CRAQ request IDs
- **Transactional Boundaries**: Align CRAQ write batches with Kafka transactions
- **DDIA Link**: Apply exactly-once and transactional outbox patterns from streaming chapter

### expected_keywords
- Primary keywords: CDC, exactly-once, tail commit, transactional outbox
- Technical terms: Kafka idempotent producer, sequence number, commit log, watermark

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Outbox Table**: Store pending events alongside CRAQ state
- **Retry Handling**: Ensure replays don't duplicate downstream events
- **Ordering Guarantees**: Partition by key to preserve linearization
- **Backpressure Considerations**: Delay ack if Kafka unavailable, degrade gracefully

### bonus_keywords
- Implementation: Kafka transactions, LSU, offset commit, outbox dedupe
- Scenarios: downstream consumer crash, network partition between CRAQ and Kafka
- Trade-offs: latency vs reliability, coupling of systems

## sample_excellent - Example Excellence
"We capture CRAQ writes when the tail commits them; that commit signal triggers insertion into an outbox table with the same idempotency key used for CRAQ request dedupe. A CDC worker reads the outbox, publishes to Kafka using idempotent transactions, and only deletes entries after Kafka acknowledges the transaction, mirroring the transactional outbox pattern from DDIA. If a write is replayed due to failure, the idempotency key prevents double emission. Partitioning by primary key maintains ordering consistent with CRAQ's linearizable sequence." 

## sample_acceptable - Minimum Acceptable
"Use CRAQ's tail acknowledgment as the trigger to write to an outbox table, then stream that through Kafka with idempotent producers so each change is emitted exactly once, following DDIA's transactional outbox pattern." 

## common_mistakes - Watch Out For
- Capturing events before tail commit
- Ignoring idempotency keys when publishing to Kafka
- Not mentioning transactional or outbox pattern from DDIA
- Overlooking ordering guarantees

## follow_up_excellent - Depth Probe
**Question**: "How do you handle CDC backlog if Kafka is unavailable without blocking CRAQ writes indefinitely?"
- **Looking for**: Backpressure thresholds, degradation to at-least-once with compensating logic, alerting
- **Red flags**: Unbounded blocking, silent data loss

## follow_up_partial - Guided Probe  
**Question**: "What metadata do you store alongside the outbox entry to correlate with CRAQ requests?"
- **Hint embedded**: Request ID, version number, chain epoch
- **Concept testing**: Traceability

## follow_up_weak - Foundation Check
**Question**: "Why do you wait for a receipt before shipping an order confirmation email?"
- **Simplification**: Tail commit analog
- **Building block**: Avoid confusing duplicate notifications

## bar_raiser_question - L5→L6 Challenge
"Propose a dual-write mitigation plan if Kafka downtime exceeds the outbox retention window, incorporating saga-like compensations."

### bar_raiser_concepts
- Sagas, compensating transactions
- Retention, backpressure handling
- Business-level reconciliation
- Streaming reliability from DDIA

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 4-5 min answer + 4-5 min discussion
- **Common next topics**: Sagas, dual writes, observability of CDC
