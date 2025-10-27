---
id: occ-subjective-L5-002
type: subjective
level: L5
category: baseline
topic: occ
subtopic: range-scans-phantoms
estimated_time: 10-12 minutes
---

# question_title - Range Scans and Phantoms in OCC

## main_question - Core Question
"You need serializable range scans under OCC. Explain the phantom problem and propose an approach to detect/prevent it in this RDMA-oriented design. Discuss trade-offs."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Phantom: new/deleted rows in range after initial read
- Per-range/index-node versioning or predicate locks
- Validate index nodes (and potentially fence keys) at commit
- Abort on node version change indicating structural modification

### expected_keywords
- index node version, predicate lock, fence key, structural changes

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Cost of validating many nodes vs concurrency
- Backfill/index maintenance considerations
- Combining coarse-grained predicate locks for hot ranges

### bonus_keywords
- B-Tree, MVCC vs OCC, granularity trade-offs

## sample_excellent - Example Excellence
"Scan records by traversing index nodes and record node versions/fence keys. At commit, validate that all touched nodes’ versions are unchanged. If any node changed, abort—this detects inserted/deleted keys (phantoms). For hot ranges, consider short-lived predicate locks to reduce repeated aborts; the trade-off is lower concurrency and extra coordination."

## sample_acceptable - Minimum Acceptable
"Use versioned index nodes and re-validate them at commit to catch range changes; abort if changed."

## common_mistakes - Watch Out For
- Ignoring range-level changes and validating only leaf records
- Holding long-lived locks defeating OCC benefits
- Missing structural changes that move keys across nodes

## follow_up_excellent - Depth Probe
**Question**: "When would you add a temporary predicate lock instead of pure validation?"
- **Looking for**: Pathological contention

## follow_up_partial - Guided Probe  
**Question**: "What’s a phantom vs a write-write conflict?"
- **Hint embedded**: Range-level anomaly

## follow_up_weak - Foundation Check
**Question**: "Why do range scans need extra care under OCC?"
