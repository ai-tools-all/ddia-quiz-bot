---
id: occ-subjective-L5-001
type: subjective
level: L5
category: baseline
topic: occ
subtopic: header-cas
estimated_time: 10-12 minutes
---

# question_title - Object Header and CAS Design

## main_question - Core Question
"Design the object header for OCC with RDMA: include version number and a lock bit. Describe the atomic operation(s) needed to set the lock while verifying the expected version, and how to avoid races during commit."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Header fields: version (monotonic), lock bit (in high bit)
- Atomic compare-and-swap (CAS) to set lock iff version matches expected
- After applying write: increment version, clear lock atomically
- Idempotent operations to handle retries

### expected_keywords
- CAS, atomicity, version bump, lock clear, concurrency

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Embed checksum or short hash for integrity (optional)
- Align header for NIC-friendly access
- Combine lock+version in a single word for CAS

### bonus_keywords
- false sharing, cache line, endianness, word size

## sample_excellent - Example Excellence
"Pack lock in MSB and version in remaining bits. Commit uses CAS(old=ver||0, new=ver||1) on the header at the primary; on success, apply data write, then store header to (ver+1||0). If CAS fails, someone else changed it—abort. This ensures atomic lock/validate, preventing TOCTOU races."

## sample_acceptable - Minimum Acceptable
"Use a CAS to set a lock only if the version matches, then write data, then bump the version and clear the lock."

## common_mistakes - Watch Out For
- Separate non-atomic lock and version checks
- Forgetting to clear lock or bump version
- Non-idempotent updates

## follow_up_excellent - Depth Probe
**Question**: "How do you ensure a crash doesn’t leave a stuck lock?"
- **Looking for**: timeouts/leases or recovery routines

## follow_up_partial - Guided Probe  
**Question**: "Why pack lock+version into one machine word?"
- **Hint embedded**: Single CAS

## follow_up_weak - Foundation Check
**Question**: "What does CAS do?"
