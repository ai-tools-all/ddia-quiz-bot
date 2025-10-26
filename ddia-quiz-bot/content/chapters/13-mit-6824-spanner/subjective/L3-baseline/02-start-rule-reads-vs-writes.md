---
id: spanner-subjective-L3-002
type: subjective
level: L3
category: baseline
topic: spanner
subtopic: start-rule
estimated_time: 7-9 minutes
---

# question_title - Start Rule: Read-Only vs Read-Write

## main_question - Core Question
"Explain Spanner’s start rule for read-only vs read-write transactions. Use a timeline example to show how choosing TT.latest and (for writes) commit-wait together preserve external consistency. What breaks if a read-only picks TT.earliest instead?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Read-only picks timestamp using TrueTime.latest at start; no locks
- Read-write chooses commit ts using TrueTime.latest and performs commit-wait
- External consistency aligns real-time order with serialization order
- Reason to avoid TT.earliest for stamping snapshots

### expected_keywords
- start rule, TT.latest, TT.earliest, read-only, read-write, external consistency

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Bounded uncertainty ε and its effect on latencies
- Why read-only needs a single global snapshot timestamp across shards
- Interaction with replica freshness/applied-through watermarks

### bonus_keywords
- snapshot, monotonicity, timestamp assignment, uncertainty bounds

## sample_excellent - Example Excellence
"At start, a read-only (RO) uses ts_ro = TT.latest, ensuring its snapshot time is not in the future. A read-write (RW) picks ts_rw ≈ TT.latest during commit and then waits until TT.earliest > ts_rw (commit-wait) before releasing locks, guaranteeing ts_rw is in the past. Thus if RW1 completes before RO2 starts, ts_rw < ts_ro and RO2 sees RW1. If RO used TT.earliest, its timestamp could precede a just-finished commit whose ts is between earliest and latest, causing RO to miss it and violate external consistency."

## sample_acceptable - Minimum Acceptable
"RO uses TT.latest to choose a safe snapshot; RW uses TT.latest plus commit-wait. This preserves real-time order. Using TT.earliest for RO risks reading before a recent commit’s timestamp."

## common_mistakes - Watch Out For
- Saying RO needs read locks
- Confusing commit-wait with waiting for all replicas
- Claiming TT returns a single exact time

## follow_up_excellent - Depth Probe
**Question**: "How would you prove RO’s timestamp is never in the future relative to start?"
- **Looking for**: TT.latest bound reasoning

## follow_up_partial - Guided Probe  
**Question**: "Why isn’t commit-wait required for read-only?"
- **Hint embedded**: No publish of new writes; just reading a chosen snapshot

## follow_up_weak - Foundation Check
**Question**: "What does an uncertainty bound ε mean operationally?"
