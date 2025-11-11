---
id: farm-subjective-L5-003
type: subjective
level: L5
category: baseline
topic: scalability
subtopic: sharding-parallelism
estimated_time: 10-12 minutes
---

# question_title - Sharding and Transaction Coordination

## main_question - Core Question
"Explain how FaRM's 90-way sharding provides parallelism while maintaining transactional consistency. Describe how a transaction coordinator handles a transaction that reads from 5 shards and writes to 3 shards. What is the coordination overhead, and how does it scale with the number of shards accessed? Compare this to a single-shard transaction."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Each shard is a primary-backup pair; object IDs encode shard/region number
- Transaction spans multiple shards requires coordinator to send messages to all involved primaries
- LOCK messages go to primaries of all shards in write-set
- VALIDATE or LOCK messages for all shards in read-set
- Coordination overhead: network round trips to all participants, increases with shard count
- Single-shard transaction: all operations local, minimal coordination

### expected_keywords
- sharding, coordinator, multi-shard transaction, network round trip, coordination overhead, parallelism

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Massive sharding (90-way) enables high parallelism for non-conflicting transactions
- Load balancing across shards spreads hot objects
- Configuration manager maintains shard-to-server mapping
- Object placement strategies affect transaction span
- Trade-off: parallelism vs coordination cost

### bonus_keywords
- load balancing, configuration manager, object placement, co-location, cross-shard coordination

## sample_excellent - Example Excellence
"FaRM shards data 90 ways; each object's region number determines its primary/backup. Multi-shard transaction: coordinator sends LOCK to 3 primaries (write-set), VALIDATE to 5 primaries (read-only objects). One network round trip for LOCK phase, one for COMMIT-PRIMARY. Overhead scales linearly with shard count: N shards = N messages per phase. Single-shard transaction: RDMA read, CAS lock, local apply—sub-microsecond. 90-way sharding parallelizes independent transactions massively but penalizes cross-shard transactions with coordination latency. Smart object placement co-locates frequently-accessed objects to minimize cross-shard transactions."

## sample_acceptable - Minimum Acceptable
"FaRM shards data across many servers. Transactions touching multiple shards need to coordinate with all involved shards, sending messages to each. More shards accessed means more coordination overhead. Single-shard transactions are faster."

## common_mistakes - Watch Out For
- Not explaining how object IDs encode shard information
- Missing the distinction between read-set and write-set in coordination
- Claiming sharding hurts performance (it enables parallelism)
- Not quantifying coordination overhead (network round trips)

## follow_up_excellent - Depth Probe
**Question**: "Design an object placement policy that minimizes cross-shard transactions for a social network workload (users, posts, comments, likes)."
- **Looking for**: Co-locate user's data on same shard, partition by user ID, replication for popular content, denormalization strategies

## follow_up_partial - Guided Probe
**Question**: "How does massive sharding (90-way) improve throughput despite coordination costs?"
- **Hint embedded**: Parallelism for non-conflicting transactions, CPU/memory distributed

## follow_up_weak - Foundation Check
**Question**: "What is sharding and why is it used in distributed databases?"
