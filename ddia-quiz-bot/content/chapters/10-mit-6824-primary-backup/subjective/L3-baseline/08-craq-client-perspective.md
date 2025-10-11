---
id: craq-subjective-L3-008
type: subjective
level: L3
category: baseline
topic: craq
subtopic: client-behavior
estimated_time: 5-7 minutes
---

# question_title - CRAQ from a Client's Perspective

## main_question - Core Question
"As a client integrating with a CRAQ cluster, what call sequencing guarantees do you rely on to ensure idempotent operations, and how does this compare to the client requirements described for primary-backup and log-based replication in DDIA?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Idempotent Retries**: Clients must attach unique identifiers to handle replays
- **Monotonic Reads**: Reads may target any clean replica yet maintain order
- **Write Acknowledgment**: Success indicates tail commit, aligning with RSM semantics
- **DDIA Connection**: Mirrors client design guidance for primary-backup and log shipping

### expected_keywords
- Primary keywords: idempotent, request ID, retry, monotonic, tail commit
- Technical terms: duplicate suppression, linearizable contract, session consistency

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Client Libraries**: Hide replica routing using manager metadata
- **Backoff Strategy**: Align with DDIA's discussion on handling transient failure
- **Monitoring Hooks**: Surfacing dirty-read rejections as metrics
- **Comparison to Kafka Producers**: Similar use of sequence numbers

### bonus_keywords
- Implementation: UUID, sequence token, deduplication table
- Scenarios: payment retry, inventory decrement, read-after-write verification
- Trade-offs: Metadata storage, client complexity, coupling to replication semantics

## sample_excellent - Example Excellence
"Clients treat CRAQ almost like a replicated state machine from DDIA: a successful write response means the update reached the tail, so they can assume it is durable and will be visible on subsequent reads. Because the head might replay writes after a failure, clients send idempotency keys—just as DDIA recommends for primary-backup log shipping—to prevent double application. For reads, the client library chooses any clean replica without breaking monotonic read guarantees."

## sample_acceptable - Minimum Acceptable
"Clients rely on CRAQ to only acknowledge writes after the tail commits them and should include idempotency tokens so retries are safe, similar to primary-backup rules in DDIA."

## common_mistakes - Watch Out For
- Assuming acknowledgments can be trusted before tail commit
- Forgetting to mention idempotent retries
- Thinking read-after-write needs client stickiness
- Ignoring parallels with other DDIA replication patterns

## follow_up_excellent - Depth Probe
**Question**: "How would you design a client SDK that automatically routes reads to clean replicas while enforcing idempotency?"
- **Looking for**: Metadata cache, retry policy, dedupe store, manager integration
- **Red flags**: Hard-coding replica addresses, ignoring clean flag updates

## follow_up_partial - Guided Probe  
**Question**: "If a client sees a 'dirty replica' error, what should it do?"
- **Hint embedded**: Retry against tail or wait for clean ack
- **Concept testing**: Respecting clean-only reads

## follow_up_weak - Foundation Check
**Question**: "Why do online stores give you an order confirmation number, and how is that similar to CRAQ write retries?"
- **Simplification**: Idempotent confirmation IDs
- **Building block**: Safe retry semantics

## bar_raiser_question - L3→L4 Challenge
"Compare CRAQ's idempotency requirements to the write sequencing rules for Kafka producers in DDIA's streaming chapter. What client guarantees let both systems provide exactly-once semantics?"

### bar_raiser_concepts
- Producer sequence numbers vs CRAQ request IDs
- Broker acknowledgment vs tail commit
- Transactional fencing
- Exactly-once operational patterns

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 2-3 min answer + 3-4 min discussion
- **Common next topics**: Client SDK design, telemetry, retry storms
