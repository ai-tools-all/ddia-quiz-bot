---
id: fault-tolerance-subjective-L5-003
type: subjective
level: L5
category: bar-raiser
topic: fault-tolerance
subtopic: performance-optimization
estimated_time: 9-12 minutes
---

# question_title - Optimizing Raft Performance

## main_question - Core Question
"Raft prioritizes understandability over performance. Describe three significant performance optimizations for Raft and their trade-offs. How would you modify Raft to handle 10,000 writes/second while maintaining correctness?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Batching**: Amortize RPC costs across multiple entries
- **Pipelining**: Send next AppendEntries before previous acknowledged
- **Parallel Dispatch**: Leader sends to all followers simultaneously
- **Bottleneck Analysis**: Leader CPU, network bandwidth, disk I/O

### expected_keywords
- Primary keywords: throughput, latency, batching, pipelining
- Technical terms: write amplification, tail latency, p99

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Leader Lease**: Bypass consensus for reads
- **Follower Reads**: Serve stale reads with bounded staleness
- **Log Compaction**: Snapshots to prevent unbounded growth
- **Hardware Optimizations**: NVMe, RDMA, kernel bypass
- **Multi-Raft**: Shard keyspace across multiple Raft groups
- **Witness Replicas**: Voting-only servers without full log

### bonus_keywords
- Implementation: io_uring, vectored I/O, zero-copy
- Related: Multi-Paxos optimizations, EPaxos
- Metrics: operations/sec, commit latency distribution

## sample_excellent - Example Excellence
"To handle 10,000 writes/second, I'd implement several optimizations: 

1. **Batching and Pipelining**: Instead of one entry per AppendEntries, batch multiple client requests. With 1ms batching window, 10 requests become 1 RPC. Pipeline by sending batch N+1 while waiting for batch N acknowledgments. Trade-off: increases mean latency but dramatically improves throughput.

2. **Parallel Commit Path**: Use separate threads for client handling, log writing, and follower replication. Leader writes to local log while simultaneously sending to followers. Trade-off: complexity in handling failures and maintaining ordering.

3. **Efficient Storage**: Use write-ahead log on NVMe with O_DIRECT to bypass kernel caches. Implement group commit where multiple entries fsync together. Trade-off: requires careful crash recovery handling.

4. **Read Optimizations**: Implement leader leases (leader assumes it's still leader for bounded time after successful heartbeat). Reads bypass consensus during lease. Trade-off: requires loosely synchronized clocks, availability impact if clock skew.

5. **Sharding**: Partition keyspace across multiple Raft groups (Multi-Raft). Each handles subset of keys. Trade-off: loses global ordering, complex for multi-key operations.

At 10k ops/sec with 5 servers, each follower processes 20k messages/sec (AppendEntries + responses). With batching of 10, reduces to 2k messages. Modern hardware can handle this with proper implementation. Key bottleneck becomes leader's CPU for message serialization and network bandwidth for replication."

## sample_acceptable - Minimum Acceptable
"Three key optimizations: 1) Batching multiple operations into single AppendEntries to reduce RPCs, trading latency for throughput. 2) Pipelining where leader sends next batch before previous fully acknowledged, improving throughput but complicating failure handling. 3) Leader leases for read optimization, where leader serves reads without consensus for a time period after successful heartbeat, requiring clock synchronization. For 10k writes/second, combine batching to reduce network overhead and parallel processing to utilize multiple CPU cores."

## common_mistakes - Watch Out For
- Optimizations that break safety properties
- Not considering end-to-end latency impact
- Ignoring hardware limitations
- Forgetting about failure handling complexity

## follow_up_excellent - Depth Probe
**Question**: "Your batching optimization delays operations by 1ms. How do you handle a client that needs strict low-latency for specific critical operations?"
- **Looking for**: Dual paths, priority queues, adaptive batching
- **Red flags**: One-size-fits-all approach

## follow_up_partial - Guided Probe  
**Question**: "You mentioned sharding across Raft groups. How do you handle operations that span multiple shards?"
- **Hint embedded**: Distributed transactions needed
- **Concept testing**: Understanding coordination complexity

## follow_up_weak - Foundation Check
**Question**: "If sending messages is slow, would sending more messages in parallel make things faster or slower?"
- **Simplification**: Bandwidth vs message overhead
- **Building block**: Batching benefits

## bar_raiser_question - L5→L6 Challenge
"Design a hybrid system that provides both linearizable operations for critical paths and eventual consistency for analytics queries. How do you maintain correctness while serving different consistency levels from the same Raft cluster?"

### bar_raiser_concepts
- Dual read paths: consensus reads vs follower reads
- Snapshot isolation for analytics using MVCC
- Read concern levels like MongoDB
- Causal consistency tokens for session guarantees
- Trade-offs between consistency, performance, and complexity

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 5-6 min answer + 4-6 min discussion
- **Common next topics**: CockroachDB architecture, etcd optimizations, TiKV design

## assistant_answer
Use batching (group commit) and pipelining to amortize RPC/fsync costs, parallelize leader I/O and replication, and employ efficient storage (preallocated WAL, NVMe, grouped fsync). Scale via multi-Raft sharding and backpressure to hit 10k writes/s while preserving ordering and durability; optimize reads with leases when safe.

## improvement_suggestions
- Specify a latency/throughput budget (e.g., batch window, fsync cadence, network RTTs) and bottleneck analysis.
- Require a rollback/compensation plan for partial batch failures and leader changes mid-pipeline.

## improvement_exercises
### exercise_1 - Performance Budget
**Question**: "Propose a target budget to reach 10k writes/s with 5 nodes (batch size, batching window, expected RTT, fsync grouping). Identify bottlenecks and metrics to monitor."

**Sample answer**: "Batch size 16–32, 0.5–1ms window; RTT ~0.5–1ms (DC-local), group fsync every 1–2ms. Bottlenecks: leader CPU serialization, NIC bandwidth, WAL fsync. Monitor: p50/p99 commit latency, in-flight bytes per follower, fsync time, dropped acks, CPU%."

### exercise_2 - Rollback Handling
**Question**: "A leader crashes mid-pipeline with partially replicated batches. Describe how a new leader safely resumes without double-applying or violating order."

**Sample answer**: "Use matchIndex/nextIndex to resume replication at last confirmed index; new leader only commits entries with majority acks, re-sends uncommitted suffix, and state machine dedup (clientID,seq) suppresses double execution. Maintain log order and re-apply idempotently."
