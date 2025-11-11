---
id: farm-subjective-L4-001
type: subjective
level: L4
category: baseline
topic: occ
subtopic: validate-optimization
estimated_time: 8-10 minutes
---

# question_title - VALIDATE vs LOCK Messages

## main_question - Core Question
"Compare the VALIDATE operation with the LOCK operation in FaRM's commit protocol. When can VALIDATE be used instead of LOCK, what are the performance benefits, and what correctness guarantees must still be maintained?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- LOCK: used for objects in write-set; involves primary CPU, log append, version check + atomic lock set
- VALIDATE: used for read-only objects; one-sided RDMA read of header, checks version and lock bit
- Performance: VALIDATE avoids log appends and primary processing, dramatically faster for read-heavy transactions
- Correctness: VALIDATE must still abort if version changed or lock bit set, ensuring serializability

### expected_keywords
- write-set, read-set, log append, one-sided, header, serializability, read-only transaction

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Read-only transactions can use VALIDATE for all objects
- VALIDATE enables order-of-magnitude throughput improvement for read workloads
- Both operations detect concurrent conflicts, just different mechanisms
- VALIDATE still requires network round trip but avoids server work

### bonus_keywords
- throughput, read-heavy, network round trip, conflict detection, primary processing

## sample_excellent - Example Excellence
"LOCK sends messages to primaries to check versions and atomically set lock bits, appending to WAL, used for objects being written. VALIDATE performs one-sided RDMA reads of object headers to check version and lock bit, used for read-only objects. Read-only transactions use VALIDATE exclusively, avoiding expensive log writes and primary CPU usage while maintaining serializability—transactions still abort if headers show version changes or locked objects."

## sample_acceptable - Minimum Acceptable
"LOCK is for writes and involves the server CPU. VALIDATE is for reads and uses one-sided RDMA. Both check versions and lock bits. VALIDATE is faster."

## common_mistakes - Watch Out For
- Claiming VALIDATE doesn't ensure serializability
- Not recognizing that VALIDATE can also cause aborts
- Missing the log append cost savings
- Thinking VALIDATE eliminates network communication entirely

## follow_up_excellent - Depth Probe
**Question**: "Could a transaction use VALIDATE for objects it plans to write? Why or why not?"
- **Looking for**: No, writes require atomic lock + WAL entry for crash recovery and replication

## follow_up_partial - Guided Probe
**Question**: "What specific fields in the object header does VALIDATE check?"
- **Hint embedded**: Version number and lock bit

## follow_up_weak - Foundation Check
**Question**: "Why is avoiding log appends important for performance?"
