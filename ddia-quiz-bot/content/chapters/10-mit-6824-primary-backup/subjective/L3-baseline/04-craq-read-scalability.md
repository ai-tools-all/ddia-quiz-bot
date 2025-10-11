---
id: craq-subjective-L3-004
type: subjective
level: L3
category: baseline
topic: craq
subtopic: read-throughput
estimated_time: 5-7 minutes
---

# question_title - CRAQ Read Scaling Fundamentals

## main_question - Core Question
"In CRAQ, why can clients read from any replica without breaking linearizability? Describe how this differs from classic chain replication and relate it to the 'reading your own writes' consistency pattern from DDIA."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Clean vs Dirty State**: Only replicas marked clean may serve reads
- **Tail Confirmation**: Writes become clean after tail acknowledgment
- **Linearizability Preservation**: Read path respects write order
- **DDIA Connection**: Mirrors 'reading your own writes' by ensuring visibility sequencing

### expected_keywords
- Primary keywords: clean replica, dirty replica, tail acknowledgment, linearizable read
- Technical terms: chain replication, read-after-write, visibility, metadata flag

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Version Counters**: Metadata tying reads to acknowledged versions
- **Replica Caching**: Faster local reads once clean
- **Replication Lag Analogy**: CRAQ avoids stale reads seen in async follower setups
- **Client Stickiness**: Not required unlike leader/follower designs

### bonus_keywords
- Implementation: lease, version vector, dirty flag propagation
- Scenarios: failover read safety, cross-datacenter reads, cache invalidation
- Trade-offs: additional metadata, delayed read availability until clean

## sample_excellent - Example Excellence
"CRAQ tags each replica's object state as dirty after the head forwards a write. The object only becomes clean once the update reaches the tail and the tail sends the confirmation back through the chain. Only clean copies are allowed to answer reads, so a client that randomly picks a middle replica will still see linearizable data: if a newer write exists, that replica is still marked dirty and must defer to the tail. Classic chain replication routes every read to the tail, which is simpler but a throughput bottleneck. CRAQ's clean/dirty metadata lets us spray reads across the chain while preserving the single-system-order semantics DDIA highlights in the 'reading your own writes' consistency rule." 

## sample_acceptable - Minimum Acceptable
"CRAQ keeps replicas dirty until the tail confirms the write, and only clean replicas can answer reads. That way, even if you read from the middle of the chain, you won't see stale data. Chain replication sent every read to the tail; CRAQ spreads reads out but still respects the read-after-write consistency pattern from DDIA."

## common_mistakes - Watch Out For
- Thinking clean/dirty is optional metadata
- Assuming any dirty replica can answer with older state
- Ignoring tail acknowledgment step
- Forgetting to compare with standard chain replication

## follow_up_excellent - Depth Probe
**Question**: "If a client issues a write then immediately reads from a middle replica, outline the message flow that guarantees the read sees the new value."
- **Looking for**: Head marks dirty, tail commits, acknowledgment flows backward, clean flag flip, read served
- **Red flags**: Letting dirty replica respond before tail ack

## follow_up_partial - Guided Probe  
**Question**: "How does CRAQ avoid the stale-read problem found in asynchronous follower replication?"
- **Hint embedded**: Clean flag withheld until commit
- **Concept testing**: Relationship to replication lag

## follow_up_weak - Foundation Check
**Question**: "If a book is marked 'unchecked' on a library shelf, should a reader trust it yet? How does CRAQ use a similar idea?"
- **Simplification**: Dirty vs clean books
- **Building block**: Trust gating via metadata

## bar_raiser_question - L3→L4 Challenge
"Suppose network latency between head and tail spikes, delaying tail acknowledgments. How does that affect CRAQ read throughput compared to traditional chain replication, and what mitigation borrowed from DDIA's caching chapter could help?"

### bar_raiser_concepts
- Tail latency as throughput limiter
- Impact on clean flag propagation
- Cache freshness vs staleness trade-offs
- Hint toward invalidation plus version pinning

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 2-3 min answer + 3-4 min discussion
- **Common next topics**: Dirty/clean reconciliation, configuration manager, replication lag handling
