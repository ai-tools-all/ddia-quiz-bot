---
id: occ-subjective-L4-002
type: subjective
level: L4
category: baseline
topic: occ
subtopic: conflicts-and-retries
estimated_time: 9-11 minutes
---

# question_title - Conflict Analysis and Retries

## main_question - Core Question
"Given overlapping transactions T1 and T2 with read/write sets: T1 reads A,B and writes B; T2 reads B,C and writes B. Analyze possible interleavings under OCC: which transaction aborts and why? Propose a retry/backoff strategy to reduce repeated conflicts."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Validation detects write-write and read-write conflicts at B
- First committer succeeds; the other aborts on version/lock conflict
- Backoff reduces collision probability on hot keys
- Shorter critical sections (smaller write-set) reduce conflicts

### expected_keywords
- read set, write set, validation, conflict, backoff, hot key

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Randomized exponential backoff vs fixed delay
- Priority schemes (e.g., older transaction wins)
- Application-level batching/combining updates

### bonus_keywords
- starvation, fairness, contention management

## sample_excellent - Example Excellence
"Both touch B; if T1 locks/validates first, it commits, incrementing B’s version. T2’s validation observes B’s lock or version change, so T2 aborts. To avoid livelock, use exponential backoff with jitter, optionally prioritizing long-waiting transactions. Reduce write-set scope or batch writes to lessen conflicts."

## sample_acceptable - Minimum Acceptable
"One commits first; the other aborts on validation. Use backoff to reduce repeated collisions."

## common_mistakes - Watch Out For
- Claiming both can commit
- Ignoring write-write conflicts
- Deterministic retry intervals that synchronize collisions

## follow_up_excellent - Depth Probe
**Question**: "When would you consider introducing a lightweight lock to guard the hottest key?"
- **Looking for**: Pathological hotspots where OCC aborts dominate

## follow_up_partial - Guided Probe  
**Question**: "Why is randomized backoff preferable to a fixed delay?"
- **Hint embedded**: Avoid re-collision

## follow_up_weak - Foundation Check
**Question**: "What’s a write-write conflict?"
