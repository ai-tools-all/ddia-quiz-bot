---
id: occ-subjective-L6-001
type: subjective
level: L6
category: baseline
topic: occ
subtopic: rdma-pipeline
estimated_time: 12-15 minutes
---

# question_title - RDMA Pipeline and Latency Budget

## main_question - Core Question
"Outline an end-to-end latency budget for a short read-write transaction (two reads, one write) targeting ~tens of microseconds. Identify where one-sided RDMA, validation, lock CAS, and commit messages fit. Propose two optimizations to improve P99 without hurting correctness."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- One-sided RDMA reads in execution path
- Lock CAS and validation at commit
- Primary commit then backup propagation
- Tail latency dominated by network/memory queues and retries

### expected_keywords
- one-sided RDMA, CAS, commit, backup, P50 vs P99

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Coalescing validations, parallel lock acquisition, batching
- NUMA/affinity and queue pair placement
- Backoff tuning for reduced re-collisions

### bonus_keywords
- qp, NIC queues, coalescing, jitter

## sample_excellent - Example Excellence
"Execution: 2 RDMA reads (~µs). Commit: parallel LOCK CAS to primaries (validate+lock), on success write+version bump, then send backup updates. P99 risk from queueing and retries; optimize with parallel lock acquisition and header-batch validation, and pin coordinator threads to reduce cache/QP churn."

## sample_acceptable - Minimum Acceptable
"Reads via RDMA, then lock/validate, commit primaries, send to backups. Improve P99 with batching/parallelization."

## common_mistakes - Watch Out For
- Serializing lock acquisitions unnecessarily
- Waiting for backups before releasing primaries’ locks without reason
- Ignoring retry impact on tails

## follow_up_excellent - Depth Probe
**Question**: "How would you detect and mitigate QP saturation?"
- **Looking for**: metrics, backpressure, queue scaling

## follow_up_partial - Guided Probe  
**Question**: "What happens to P99 if you serialize lock CAS across shards?"
- **Hint embedded**: Head-of-line blocking

## follow_up_weak - Foundation Check
**Question**: "What is one-sided RDMA?"
