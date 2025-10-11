---
id: craq-subjective-L4-005
type: subjective
level: L4
category: baseline
topic: craq
subtopic: write-path
estimated_time: 6-8 minutes
---

# question_title - Comparing CRAQ Writes with Output Commit

## main_question - Core Question
"Explain how CRAQ's write acknowledgement path relates to the output commit problem covered in DDIA's VMware FT discussion. What guarantees do both systems enforce before revealing results externally, and where do they differ?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Tail Acknowledgment**: CRAQ waits for tail before confirming write
- **Output Commit Analogy**: External visibility only after backup aligned
- **Difference in Scope**: CRAQ focuses on state visibility; VMware FT deals with I/O outputs
- **Consistency Assurance**: Both prevent state rollback seen by clients

### expected_keywords
- Primary keywords: tail ack, output commit, external visibility, rollback
- Technical terms: replicated log, deterministic replay, response gating, failover

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Latency Implications**: Added round trip, potential batching
- **Failure Handling**: Replay ensures same response re-sent, similar to FT
- **Client Acknowledgment Strategy**: Exactly-once semantics via idempotent responses
- **Operational Differences**: CRAQ replicates data structures, FT replicates VM I/O stream

### bonus_keywords
- Implementation: write barrier, ack quorum, output buffering, commit index
- Scenarios: network hiccup, tail crash, client retry
- Trade-offs: throughput vs correctness, synchronous vs asynchronous commit

## sample_excellent - Example Excellence
"Tailing a write in CRAQ mirrors VMware FT's output commit solution from DDIA: neither system lets the outside world observe a change until the backup (CRAQ's tail or FT's secondary VM) has received the state transition. CRAQ sends the client's acknowledgment only after the tail applies the update, ensuring any failover will reproduce the same log entry. VMware FT buffers outbound packets until the secondary has replayed the instruction stream. The difference is granularity—CRAQ deals with key-value records while FT gating happens on arbitrary VM outputs—but both guarantee that failover won't force clients to observe a rollback." 

## sample_acceptable - Minimum Acceptable
"CRAQ delays write acknowledgments until the tail confirms the change, which is similar to VMware FT's output commit rule: don't show results externally until the backup is ready. The difference is CRAQ cares about data updates while FT blocks network outputs." 

## common_mistakes - Watch Out For
- Claiming CRAQ can acknowledge before tail commit
- Ignoring VMware FT's buffering of outputs
- Missing the shared goal of preventing rollback
- Overlooking latency trade-offs

## follow_up_excellent - Depth Probe
**Question**: "If clients need faster perceived acknowledgment, how could you layer a speculative response without violating the guarantees?"
- **Looking for**: Dual acknowledgments, promise vs commit, client-visible state gating
- **Red flags**: Allowing speculative results to leak without reconciliation

## follow_up_partial - Guided Probe  
**Question**: "What happens if the tail crashes after receiving the write but before acknowledging it?"
- **Hint embedded**: Replacement tail catch-up, need for idempotent resend
- **Concept testing**: Failure semantics

## follow_up_weak - Foundation Check
**Question**: "Why do we wait for our teammate to confirm they copied the homework before handing it in?"
- **Simplification**: Output commit analogy
- **Building block**: Ensuring backup copy exists

## bar_raiser_question - L4→L5 Challenge
"Design an optimization that overlaps CRAQ tail acknowledgments with downstream cache invalidations, similar to DDIA's discussion of write pipelines."

### bar_raiser_concepts
- Overlap of replication and external side effects
- Cache invalidation timing
- Ensuring ordering between write commit and invalidation
- Pipelining for latency reduction

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 3-4 min answer + 3-4 min discussion
- **Common next topics**: Speculative execution, write pipelines, cache coordination
