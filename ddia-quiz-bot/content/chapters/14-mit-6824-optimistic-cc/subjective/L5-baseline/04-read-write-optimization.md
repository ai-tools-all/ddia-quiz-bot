---
id: farm-subjective-L5-004
type: subjective
level: L5
category: baseline
topic: optimization
subtopic: read-write-mixed
estimated_time: 10-12 minutes
---

# question_title - Optimizing Mixed Read-Write Transactions

## main_question - Core Question
"Design optimization strategies for a transaction that reads 100 objects but writes only 3 of them. Compare three approaches: (1) using LOCK for all objects, (2) using VALIDATE for reads and LOCK for writes, (3) caching read values and validating at commit. For each approach, analyze latency, throughput, abort rate, and consistency guarantees."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Approach 1 (all LOCK): high overhead, log appends for all 100 objects, strong isolation
- Approach 2 (VALIDATE/LOCK): VALIDATE for 97 reads (one-sided RDMA), LOCK for 3 writes, optimal
- Approach 3 (caching): eliminate repeated network reads, but must validate all cached objects
- VALIDATE avoids log writes and primary CPU, dramatically faster for reads
- All approaches provide serializability if validation includes all read objects

### expected_keywords
- VALIDATE, LOCK, optimization, log append, one-sided RDMA, validation overhead, serializability

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Repeated reads benefit from caching but increase abort risk (longer validation window)
- Read-heavy workloads benefit most from VALIDATE optimization
- Write amplification problem with locking read-only objects
- Application-level hints about read-only vs read-write objects

### bonus_keywords
- caching, write amplification, read-heavy, validation window, application hints, bloom filters

## sample_excellent - Example Excellence
"Approach 1: LOCK all 103 objects—100 log appends, primary processing, ~100μs. Approach 2 (best): VALIDATE 97 reads via one-sided RDMA header checks (~1μs each), LOCK 3 writes (~10μs each), total ~130μs vs 1000μs. Approach 3: cache 97 values locally, LOCK 3, validate 97 headers at commit—saves repeated RDMA reads if objects accessed multiple times, but increases abort risk (longer window for conflicts). All provide serializability since validation checks all read versions. Approach 2 optimal for single-pass reads; Approach 3 better if objects accessed 3+ times."

## sample_acceptable - Minimum Acceptable
"Use VALIDATE for objects you only read and LOCK for objects you write. This is faster because VALIDATE uses one-sided RDMA and doesn't write logs. All approaches maintain consistency."

## common_mistakes - Watch Out For
- Thinking VALIDATE compromises consistency (it doesn't)
- Not quantifying the performance difference
- Missing the write amplification problem
- Forgetting that cached values still need validation

## follow_up_excellent - Depth Probe
**Question**: "How would you extend VALIDATE to support range queries efficiently while maintaining serializability?"
- **Looking for**: Index versioning, predicate locks, range version counters, gap locks

## follow_up_partial - Guided Probe
**Question**: "Why does VALIDATE avoid log appends while LOCK requires them?"
- **Hint embedded**: Writes need durability and recovery information, reads don't

## follow_up_weak - Foundation Check
**Question**: "What is write amplification and why is it a problem?"
