---
id: farm-subjective-L5-002
type: subjective
level: L5
category: baseline
topic: occ
subtopic: phantoms
estimated_time: 10-12 minutes
---

# question_title - Phantom Reads in OCC

## main_question - Core Question
"Explain the phantom read problem in the context of FaRM's OCC protocol. Consider a transaction that reads a range of objects (e.g., all accounts in a region) and another transaction that inserts a new object in that range. How does FaRM's version-based validation handle or not handle this scenario? What modifications to the protocol would prevent phantom reads?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Phantom: new objects appearing in a range query that wasn't there during initial read
- FaRM's per-object versioning only tracks existing objects, not insertions/deletions
- Standard OCC validation would miss phantoms because new object wasn't in read-set
- Solutions: predicate locks, range versioning, or index versioning to track structural changes

### expected_keywords
- phantom read, range query, insertion, read-set, predicate lock, serializability, isolation

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Difference between snapshot isolation and full serializability
- Index or region versioning as practical solution
- Trade-off: performance vs strict serializability for range operations
- Application-level strategies to avoid phantoms

### bonus_keywords
- snapshot isolation, index, partition version, serializable snapshot isolation, precision locking

## sample_excellent - Example Excellence
"T1 scans region R finding objects {A, B}, recording versions. T2 inserts C into R and commits. T1's validation checks A and B versions (unchanged), so commits—but now sees different results if re-scanned (phantom C). Per-object versions can't detect this. Solutions: (1) predicate locks on range predicates, (2) version counter for region/index structure itself, (3) serializable snapshot isolation with write conflict detection on ranges."

## sample_acceptable - Minimum Acceptable
"Phantoms happen when new objects are inserted into a range a transaction already read. Per-object versioning doesn't catch this because the new object wasn't in the read-set. Need range or index versioning."

## common_mistakes - Watch Out For
- Confusing phantom reads with regular write-write conflicts
- Claiming OCC automatically prevents phantoms (it doesn't without extra mechanisms)
- Not explaining why per-object versions are insufficient
- Missing the connection to serializability vs snapshot isolation

## follow_up_excellent - Depth Probe
**Question**: "Design a concrete region versioning scheme: where is the version stored, when is it incremented, and how do transactions validate it?"
- **Looking for**: Version in partition metadata, increment on insert/delete, transactions read region version and validate like object versions

## follow_up_partial - Guided Probe
**Question**: "Why doesn't validating individual object versions catch phantom insertions?"
- **Hint embedded**: The new object wasn't part of the original read-set

## follow_up_weak - Foundation Check
**Question**: "What is a phantom read and why is it a problem?"
