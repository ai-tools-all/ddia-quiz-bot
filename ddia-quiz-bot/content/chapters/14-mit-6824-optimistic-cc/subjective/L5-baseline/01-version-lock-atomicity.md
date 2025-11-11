---
id: farm-subjective-L5-001
type: subjective
level: L5
category: baseline
topic: occ
subtopic: atomic-operations
estimated_time: 10-12 minutes
---

# question_title - Atomic Version Check and Lock

## main_question - Core Question
"Design the object header structure and atomic operations needed for FaRM's OCC protocol. Explain why the version check and lock bit setting must be atomic, what hardware primitive achieves this, and what race conditions would occur if these operations were separate. Include how to handle lock clearing and version incrementing after commit."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Header fields: version number (monotonic counter), lock bit (typically high-order bit)
- Compare-and-swap (CAS) to atomically verify version and set lock in one operation
- Race condition: TOCTOU (time-of-check-time-of-use) if check and set are separate
- After write applied: atomic update to increment version and clear lock

### expected_keywords
- CAS, atomic, TOCTOU, race condition, compare-and-swap, header, version bump, lock clear

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Pack lock bit and version in single word (e.g., 64-bit: 1 bit lock, 63 bits version)
- Alignment for efficient NIC access
- Idempotent operations for retry handling
- Alternative: lock as separate bit with double-CAS or 128-bit CAS

### bonus_keywords
- word size, alignment, idempotent, retry, packing, bit manipulation

## sample_excellent - Example Excellence
"Pack lock bit in MSB of version word. LOCK uses CAS(expected=ver||0, new=ver||1) to atomically verify version matches and set lock; CAS failure means version changed or already locked—abort. After applying writes, store (ver+1||0) to increment version and clear lock. Atomicity is critical: if check and set were separate, T1 could check version, T2 commits between check and lock set, then T1 sets lock on wrong version—violating serializability."

## sample_acceptable - Minimum Acceptable
"Use CAS to check version and set lock atomically. If separate, another transaction could commit between the check and lock, causing incorrect commit. After commit, increment version and clear lock."

## common_mistakes - Watch Out For
- Separate non-atomic check and lock operations
- Forgetting version increment or lock clear
- Not explaining the TOCTOU race
- Missing the role of CAS hardware primitive

## follow_up_excellent - Depth Probe
**Question**: "How would you handle a crash that leaves a lock bit set but no transaction to clear it?"
- **Looking for**: Timeouts/leases, recovery scans, epoch numbers, transaction monitoring

## follow_up_partial - Guided Probe
**Question**: "Why pack lock bit and version into a single word rather than separate fields?"
- **Hint embedded**: Single CAS operation, hardware atomicity limits

## follow_up_weak - Foundation Check
**Question**: "What is a compare-and-swap (CAS) operation?"
