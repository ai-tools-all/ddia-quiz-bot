---
id: farm-subjective-L3-004
type: subjective
level: L3
category: baseline
topic: correctness
subtopic: serializability
estimated_time: 6-8 minutes
---

# question_title - Serializability Through Validation

## main_question - Core Question
"Explain how FaRM's version-based validation at commit time ensures serializability. Walk through a specific example with two transactions T1 and T2 that both read X and Y, but T1 writes X while T2 writes Y. Show how the validation prevents non-serializable interleavings."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Each transaction records versions of all objects read
- At commit, validation checks versions haven't changed and no lock bits set
- If validation succeeds for all objects, transaction commits; else aborts
- Validation ensures transaction saw consistent snapshot and no concurrent modifications
- Example: T1 and T2 both read X@v1,Y@v2; commits serialize based on lock acquisition order

### expected_keywords
- version validation, serializability, consistent snapshot, isolation, commit order, lock acquisition

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Corresponds to optimistic read phase + pessimistic commit phase
- Version validation is equivalent to read-set conflict detection
- Lock bits provide write-set conflict detection
- Combination gives full serializability (not just snapshot isolation)

### bonus_keywords
- read-set, write-set, snapshot isolation, conflict serializability, commit timestamp

## sample_excellent - Example Excellence
"T1 and T2 both read X@v1,Y@v2 via RDMA. T1 commits first: locks X, validates Y@v2 (ok), writes X@v2, clears lock. T2 then commits: locks Y, validates X—finds X@v2 (expected v1)—aborts. This enforces T1→T2 serial order. If both validated successfully, they'd be serializable since they modify disjoint objects. Version checks ensure each transaction operates on a consistent snapshot without seeing partial effects of concurrent transactions."

## sample_acceptable - Minimum Acceptable
"Transactions record versions when reading. At commit, they check versions haven't changed. If changed, abort. This ensures serializability by preventing transactions from using stale data."

## common_mistakes - Watch Out For
- Claiming OCC only provides snapshot isolation (it provides serializability)
- Not explaining how version checks detect conflicts
- Missing the connection between validation and consistent snapshots
- Confusing read-write conflicts with write-write conflicts

## follow_up_excellent - Depth Probe
**Question**: "Would FaRM's protocol still be serializable if it only checked versions for objects in the write-set, not read-set?"
- **Looking for**: No—missing read-write conflicts, would degrade to snapshot isolation

## follow_up_partial - Guided Probe
**Question**: "What would happen if versions were checked before locking instead of atomically with locking?"
- **Hint embedded**: TOCTOU race condition

## follow_up_weak - Foundation Check
**Question**: "What does serializability mean in database transactions?"
