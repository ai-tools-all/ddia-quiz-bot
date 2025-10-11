---
id: fault-tolerance-subjective-L4-002
type: subjective
level: L4
category: baseline
topic: fault-tolerance
subtopic: commitment-rules
estimated_time: 6-8 minutes
---

# question_title - Log Commitment and Safety

## main_question - Core Question
"Explain Raft's commitment rules. Why can a leader only commit entries from its current term, and what safety problems could arise if this restriction didn't exist?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Majority Replication**: Entry committed when replicated on majority
- **Current Term Only**: Leader commits only its own term's entries
- **Indirect Commitment**: Previous terms' entries committed indirectly
- **State Machine Application**: Only committed entries applied

### expected_keywords
- Primary keywords: commit, majority, replication, safety
- Technical terms: commit index, applied, state machine

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Leader Completeness**: Committed entries in all future leaders' logs
- **Figure 8 Scenario**: Classic example of commitment confusion
- **Commitment Notification**: How followers learn of commits
- **Durability Guarantee**: Committed means permanent

### bonus_keywords
- Implementation: leaderCommit, lastApplied indices
- Related: consensus, linearizability
- Safety properties: never reverting commits

## sample_excellent - Example Excellence
"Raft's commitment rule states that a leader can only commit entries from its current term, not from previous terms. An entry is committed when it's replicated on a majority of servers AND at least one entry from the current term is also replicated on a majority. This prevents a subtle safety violation: imagine a leader from term 2 replicates an entry to a majority but crashes before committing. A new leader in term 3 is elected without this entry. If term 3's leader could directly commit the term 2 entry, and then crash, a term 4 leader might be elected with different entries at that index. The term 2 entry would be 'uncommitted' retroactively, violating safety. By only committing current-term entries, which indirectly commits all preceding entries, Raft ensures that once committed, an entry appears in all future leaders' logs. This is the Leader Completeness Property - fundamental to Raft's correctness."

## sample_acceptable - Minimum Acceptable
"A leader can only commit log entries from its own term, not previous terms. An entry is committed when stored on a majority of servers. This restriction prevents situations where an entry could be considered committed but then disappear if a new leader is elected. Previous entries get committed indirectly when a current-term entry is committed."

## common_mistakes - Watch Out For
- Thinking any majority replication means commitment
- Not understanding indirect commitment mechanism
- Missing the safety violation without this rule
- Confusing replication with commitment

## follow_up_excellent - Depth Probe
**Question**: "Walk through a specific scenario where a term 2 entry is replicated to 3 of 5 servers but not committed. Show how committing it in term 4 could lead to it being lost."
- **Looking for**: Clear scenario construction showing violation
- **Red flags**: Not seeing the retroactive uncommit problem

## follow_up_partial - Guided Probe  
**Question**: "You mentioned indirect commitment. If a leader has entries from terms 1, 2, and 3, and commits a term 3 entry, what happens to the term 1 and 2 entries?"
- **Hint embedded**: All get committed together
- **Concept testing**: Understanding commitment cascade

## follow_up_weak - Foundation Check
**Question**: "Think of commitment like finalizing a decision. Why might it be dangerous to finalize someone else's pending decision without adding your own?"
- **Simplification**: Authority and decision ownership
- **Building block**: Permanence of commitment

## bar_raiser_question - L4→L5 Challenge
"Consider this log state across 5 servers:
- S1 (leader, term 4): [1:X, 2:Y, 2:Z, 4:W]
- S2: [1:X, 2:Y, 2:Z]
- S3: [1:X, 2:Y]
- S4: [1:X]
- S5: [1:X, 3:A]

What is the highest index that can be committed? Walk through S1's process of achieving commitment."

### bar_raiser_concepts
- S1 must first replicate 4:W to achieve majority
- Once 4:W is on 3 servers, can commit up to index 4
- This indirectly commits entries at indices 1, 2, 3
- S5's conflicting entry (3:A) will be overwritten
- Demonstrates real commitment process

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 3-4 min answer + 3-4 min discussion
- **Common next topics**: Leader completeness, election restriction, safety proof
