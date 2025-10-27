---
id: occ-subjective-L3-002
type: subjective
level: L3
category: baseline
topic: occ
subtopic: read-only-validation
estimated_time: 7-9 minutes
---

# question_title - Read-Only Validation Path

## main_question - Core Question
"Design the read-only transaction path using RDMA one-sided operations. What must be checked during validation to ensure serializability, and when should the client abort?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- One-sided reads of data; record version numbers
- Validation re-reads headers (version, lock bit)
- Abort if any version changed or lock bit observed
- No locks needed for read-only

### expected_keywords
- RDMA, header, version, lock bit, validation, abort

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Phantom challenges for range scans and index-level validation
- Handling cache/coherency delays; retries with backoff
- Batch header fetch to reduce round-trips

### bonus_keywords
- range scans, index node version, coalescing

## sample_excellent - Example Excellence
"RO reads via RDMA and records versions. Before responding, it RDMA-reads headers of all read objects; if any lock bit is set or version differs, it aborts and retries. This prevents returning a snapshot influenced by concurrent writes, preserving serializability without holding locks."

## sample_acceptable - Minimum Acceptable
"Read data, then re-check versions/lock bits. Abort on change; otherwise return."

## common_mistakes - Watch Out For
- Taking read locks unnecessarily
- Skipping validation because ‘reads don’t conflict’
- Ignoring lock bit during validation

## follow_up_excellent - Depth Probe
**Question**: "How would you handle range scans to avoid phantoms?"
- **Looking for**: index/node versioning, fence keys

## follow_up_partial - Guided Probe  
**Question**: "Why check the lock bit if you’re not modifying objects?"
- **Hint embedded**: Writer may be committing concurrently

## follow_up_weak - Foundation Check
**Question**: "What’s a version number used for?"
