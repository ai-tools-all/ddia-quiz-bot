---
id: craq-subjective-L4-004
type: subjective
level: L4
category: baseline
topic: craq
subtopic: dirty-clean-reconciliation
estimated_time: 6-8 minutes
---

# question_title - Handling Dirty Replica Recovery

## main_question - Core Question
"Outline the steps a CRAQ replica takes to transition from dirty to clean after catching up from a snapshot. Compare this to the log compaction and state transfer mechanisms Kleppmann describes for replicated state machines." 

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Snapshot Install**: Replica applies snapshot to reach recent baseline
- **Log Replay**: Replays remaining tail-confirmed writes to achieve clean state
- **Clean Flag Switch**: Only flips after verifying complete prefix
- **DDIA Link**: Mirrors RSM catch-up (snapshot + log tail) ensuring deterministic state

### expected_keywords
- Primary keywords: snapshot, log replay, clean state, deterministic apply
- Technical terms: checkpoint, log compaction, state transfer, idempotent replay

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Checksums**: Validating snapshot integrity (ties to fail-stop conversion)
- **Version Negotiation**: Ensuring schema compatibility per DDIA encoding chapter
- **Bandwidth vs Latency Trade-off**: Snapshot reduces replay length
- **Operational Automation**: Health check gating before rejoining read pool

### bonus_keywords
- Implementation: snapshot sequence number, catch-up marker, commit index
- Scenarios: high churn leadership, cold-start replica, bulk resync after outage
- Trade-offs: snapshot size, downtime, disk IO pressure

## sample_excellent - Example Excellence
"A CRAQ replica recovering from scratch downloads the latest committed snapshot (say every 10 MB) to approximate the tail's state, then replays log entries after the snapshot point until it reaches the tail's commit index. Only once the replica has processed all tail-confirmed writes does it mark those objects clean and rejoin the read pool. This is identical to the snapshot+log catch-up sequence DDIA prescribes for replicated state machines: snapshots bound recovery time while replay ensures deterministic convergence."

## sample_acceptable - Minimum Acceptable
"Recovering CRAQ replicas install a snapshot, replay newer writes, and only become clean after all tail-confirmed entries are applied—just like the RSM recovery flow DDIA describes."

## common_mistakes - Watch Out For
- Skipping snapshot verification
- Allowing replica to serve reads before full replay
- Forgetting to mention tail confirmation vs arbitrary log end
- Ignoring schema/version alignment concerns

## follow_up_excellent - Depth Probe
**Question**: "How do you guarantee the snapshot and log tail were produced from the same epoch, and what DDIA consistency principle does this protect?"
- **Looking for**: Snapshot metadata, epoch/fencing tokens, linearizability
- **Red flags**: Mixing data from different chain configurations

## follow_up_partial - Guided Probe  
**Question**: "If the snapshot is large, what heuristics can you adopt to keep recovery time acceptable?"
- **Hint embedded**: Incremental snapshots, background transfer, compression
- **Concept testing**: Operational trade-offs

## follow_up_weak - Foundation Check
**Question**: "Why do you copy someone else's homework before filling in the latest answers yourself?"
- **Simplification**: Snapshot + log replay analogy
- **Building block**: Baseline first, then incremental updates

## bar_raiser_question - L4→L5 Challenge
"Design a schema evolution strategy for CRAQ snapshots that aligns with DDIA's encoding and evolution guidance. How do you ensure old replicas can still recover after a format change?"

### bar_raiser_concepts
- Backward/forward compatibility
- Snapshot versioning, feature flags
- Dual-write window, transformation pipelines
- Operational rollout sequencing

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 3-4 min answer + 3-4 min discussion
- **Common next topics**: Snapshot scheduling, schema management, rolling upgrades
