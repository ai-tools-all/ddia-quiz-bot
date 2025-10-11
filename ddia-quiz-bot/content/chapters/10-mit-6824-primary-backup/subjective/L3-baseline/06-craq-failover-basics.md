---
id: craq-subjective-L3-006
type: subjective
level: L3
category: baseline
topic: craq
subtopic: failover-basics
estimated_time: 5-7 minutes
---

# question_title - CRAQ Failover Fundamentals

## main_question - Core Question
"A middle replica in CRAQ fails while holding dirty state. Describe how the chain continues serving reads and writes without violating the linearizable guarantees discussed in DDIA's consensus chapter."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Dirty State Handling**: Successor must fetch missing writes before becoming clean
- **Configuration Update**: Manager redefines chain without failed node
- **Linearizable Order Preservation**: Writes replayed in-order through remaining nodes
- **Consensus Analogy**: Similar to Raft log catch-up ensuring committed prefix consistency

### expected_keywords
- Primary keywords: dirty replica, replay, reconfiguration, linearizability
- Technical terms: log backfill, catch-up, committed prefix, fencing token

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Stall Reads**: Reads may temporarily route to tail until replica cleans up
- **Telemetry**: Detect dirty backlog using metrics akin to replication lag dashboards
- **Client Retries**: Idempotent writes / duplicate suppression at head
- **Comparison to Primary-Backup**: Replay step similar to state transfer

### bonus_keywords
- Implementation: backlog queue, snapshot install, handoff token
- Scenarios: fail-stop assumption, network partition, slow replica promotion
- Trade-offs: recovery latency vs availability

## sample_excellent - Example Excellence
"When a middle CRAQ replica fails, the configuration manager removes it and reconnects the surviving nodes. Any dirty writes that hadn't reached the tail are still in-flight; the upstream node keeps them in its outbound queue and forwards them to the next downstream node, ensuring the committed prefix remains identical—just like the log catch-up guarantee Raft gives in DDIA. Until the replacement finishes replaying and marks the entries clean, reads fall back to other clean replicas (often the tail). This preserves linearizability at the cost of a short performance dip." 

## sample_acceptable - Minimum Acceptable
"If a CRAQ replica fails while dirty, the chain reconfigures and the upstream node resends those writes to the next replica so the order stays the same. Reads may have to go to the tail until the replacement is clean."

## common_mistakes - Watch Out For
- Assuming dirty writes are lost
- Forgetting configuration manager involvement
- Letting replacement serve reads before catching up
- Ignoring similarities to log replication from DDIA

## follow_up_excellent - Depth Probe
**Question**: "How does CRAQ ensure idempotence when the upstream replica replays writes after failure?"
- **Looking for**: Write identifiers, duplicate suppression, compare to Raft index/term
- **Red flags**: Double-applying updates without guards

## follow_up_partial - Guided Probe  
**Question**: "If the replacement replica is very slow, what temporary routing policy keeps read latency acceptable?"
- **Hint embedded**: Route to tail, keep clean set authoritative
- **Concept testing**: Availability vs performance

## follow_up_weak - Foundation Check
**Question**: "When a baton runner drops out, how do the remaining runners keep the baton moving without skipping steps?"
- **Simplification**: Replay baton handoff
- **Building block**: Sequential transfer of state

## bar_raiser_question - L3→L4 Challenge
"Compare CRAQ's failure replay with the 'state transfer vs replicated state machine' approaches from DDIA. When would you switch to snapshot transfer instead of log replay?"

### bar_raiser_concepts
- Snapshot vs log catch-up cost
- State size considerations
- Recovery time vs operational simplicity
- Equivalent choices in Raft and primary-backup chapters

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 2-3 min answer + 3-4 min discussion
- **Common next topics**: Snapshotting, log compaction, monitoring recovery health
