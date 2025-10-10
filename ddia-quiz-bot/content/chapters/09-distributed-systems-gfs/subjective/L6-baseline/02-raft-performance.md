---
id: raft-subjective-L6-001
type: subjective
level: L6
category: baseline
topic: raft-consensus
subtopic: performance-optimization
estimated_time: 12-15 minutes
---

# question_title - Raft Performance Optimization and Trade-offs

## main_question - Core Question
"Raft prioritizes understandability over optimal performance. Discuss the key performance bottlenecks in vanilla Raft and explain optimization techniques used in production systems. How do these optimizations affect Raft's safety and liveness properties?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Leader Bottleneck**: All writes go through single leader
- **Serial Log Replication**: Followers process entries in order
- **Network Round Trips**: Minimum 1 RTT for commits
- **Disk Writes**: Persistent state on every update
- **Election Overhead**: Service disruption during elections

### expected_keywords
- Primary keywords: throughput, latency, bottleneck, optimization, batching
- Technical terms: pipeline, parallel apply, write coalescing, leader lease

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Pipelining**: Send multiple AppendEntries without waiting
- **Batching**: Group multiple client requests
- **Parallel Apply**: Apply to state machine while replicating
- **Read Optimization**: Lease reads, follower reads
- **Membership Changes**: Joint consensus overhead
- **Witness/Learner Nodes**: Non-voting replicas
- **Multi-Raft**: Sharding the consensus group

### bonus_keywords
- Production systems: etcd, TiKV, CockroachDB, Consul
- Techniques: Write coalescing, async apply, read index, lease read
- Metrics: Operations/sec, commit latency, election time

## sample_excellent - Example Excellence
"Raft's vanilla implementation has several performance limitations stemming from its strong consistency model and design simplicity:

**Key Bottlenecks**:
1. **Leader serialization**: Single leader processes all writes - becomes CPU/network bottleneck at scale
2. **Synchronous replication**: Wait for majority before acknowledging - adds network latency
3. **Disk persistence**: Must sync log to disk before voting/acknowledging - adds disk latency
4. **Ordered processing**: Followers apply entries sequentially - limits throughput

**Production Optimizations**:
1. **Batching and Pipelining**: Group multiple client requests into single AppendEntries. Send multiple uncommitted entries without waiting for previous acknowledgments. etcd achieves 10x throughput with batching.

2. **Parallel Apply**: Decouple log replication from state machine application. Can acknowledge client after commit but before applying, as long as reads go through log.

3. **Write Coalescing**: Combine multiple concurrent disk writes into single fsync. Amortizes disk latency across requests.

4. **Read Optimizations**: 
   - ReadIndex: Leader tracks commit index without log entry
   - Lease reads: Leader serves reads directly if within lease period
   - Follower reads: With timestamp bounds for staleness

5. **Multi-Raft/Sharding**: Partition keyspace across multiple Raft groups. TiKV uses this for horizontal scaling.

**Safety Preservation**: These optimizations maintain safety - never violate linearizability or lose committed data. Batching/pipelining only affects performance. Read optimizations use careful lease management.

**Liveness Trade-offs**: Some optimizations can affect liveness. Larger batches increase latency variance. Pipelining can cause more rollback on leader change. Generally acceptable trade-off for throughput."

## sample_acceptable - Minimum Acceptable
"Main bottlenecks in Raft include the single leader handling all writes, waiting for majority replication before committing, and required disk writes for persistence. Common optimizations include batching multiple operations together, pipelining AppendEntries RPCs without waiting for responses, and allowing reads from followers with bounded staleness. These optimizations don't violate safety but may affect timing and liveness properties."

## common_mistakes - Watch Out For
- Suggesting optimizations that violate safety
- Not recognizing leader as fundamental bottleneck
- Ignoring disk I/O impact
- Confusing Raft with eventually consistent systems

## follow_up_excellent - Depth Probe
**Question**: "How would you implement linearizable reads from followers without going through the leader? What clock assumptions would you need?"
- **Looking for**: ReadIndex propagation, hybrid logical clocks, lease mechanisms
- **Red flags**: Allowing stale reads without bounds

## follow_up_partial - Guided Probe
**Question**: "You mentioned batching. How do you decide the optimal batch size? What are the trade-offs?"
- **Hint embedded**: Latency vs throughput, fairness considerations
- **Concept testing**: Understanding practical trade-offs

## follow_up_weak - Foundation Check
**Question**: "Why does Raft require disk writes before sending vote responses?"
- **Simplification**: Persistence and safety connection
- **Building block**: Understanding durability requirements

## bar_raiser_question - L6→L7 Challenge
"Design a Raft variant optimized for wide-area networks where followers are in different geographical regions. How would you minimize cross-region traffic while maintaining linearizability?"

### bar_raiser_concepts
- Hierarchical consensus
- Regional leaders
- Quorum flexibility
- WAN-aware placement
- Adaptive protocols

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 5-6 min answer + 7-9 min discussion
- **Common next topics**: EPaxos, Flexible Paxos, Calvin, Spanner
