---
id: fault-tolerance-subjective-L4-003
type: subjective
level: L4
category: bar-raiser
topic: fault-tolerance
subtopic: election-restriction
estimated_time: 8-10 minutes
---

# question_title - Election Restriction and Log Completeness

## main_question - Core Question
"Explain Raft's election restriction: why does a candidate need an 'up-to-date' log to win an election? Define what 'up-to-date' means precisely and explain what safety problems this prevents."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Up-to-date Definition**: Higher term or same term with longer log
- **Voting Restriction**: Servers deny votes to out-of-date candidates
- **Preserves Committed Entries**: Ensures elected leader has all committed entries
- **Comparison Algorithm**: Term comparison first, then log length

### expected_keywords
- Primary keywords: election restriction, up-to-date, voting, log completeness
- Technical terms: lastLogTerm, lastLogIndex, RequestVote

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Leader Completeness Property**: Mathematical guarantee
- **Prevents Rollback**: Committed entries never lost
- **Majority Overlap**: At least one voter has committed entries
- **Safety Without Checking Commit Status**: Voter doesn't need to know what's committed

### bonus_keywords
- Implementation: RequestVote RPC fields
- Related: quorum intersection, safety properties
- Proof elements: induction, invariants

## sample_excellent - Example Excellence
"Raft's election restriction requires that a candidate's log must be at least as up-to-date as the voter's log to receive its vote. 'Up-to-date' is determined by: first comparing the term of the last log entry (higher term wins), then if terms are equal, comparing log length (longer wins). This lexicographic ordering ensures the Leader Completeness Property: if an entry is committed, all future leaders will have it. Here's why: a committed entry exists on a majority of servers. To win election, a candidate needs votes from a majority. These two majorities must overlap in at least one server that has the committed entry. That server won't vote for anyone whose log doesn't contain the entry (they'd fail the up-to-date check). Therefore, only candidates with all committed entries can win. This elegantly prevents rollback of committed entries without servers explicitly tracking what's committed - the voting restriction implicitly preserves committed state through the election process."

## sample_acceptable - Minimum Acceptable
"A candidate must have a log at least as up-to-date as the servers it's requesting votes from. Up-to-date means either having a last entry with a higher term, or same term but longer log. This ensures the new leader has all committed entries, preventing loss of acknowledged client operations. Servers won't vote for candidates with older or shorter logs."

## common_mistakes - Watch Out For
- Getting the comparison order wrong (term vs length priority)
- Not explaining why majority overlap matters
- Thinking servers need to know what's committed
- Missing that this prevents committed entry loss

## follow_up_excellent - Depth Probe
**Question**: "Consider a committed entry E at index 5, term 3, replicated on servers S1, S2, S3 out of 5 total. S4 and S5 have empty logs. Could S4 or S5 ever become leader? Walk through the voting."
- **Looking for**: No - need votes from {S1,S2,S3}, all will refuse
- **Red flags**: Not recognizing majority requirement blocks them

## follow_up_partial - Guided Probe  
**Question**: "You explained the comparison uses term then length. Why not just use length alone? What would break?"
- **Hint embedded**: Old long logs could win over recent short logs
- **Concept testing**: Understanding term importance

## follow_up_weak - Foundation Check
**Question**: "Imagine choosing a team leader based on experience. Would you pick someone who missed the last three important meetings? Why not?"
- **Simplification**: Missing critical information
- **Building block**: Information completeness for decisions

## bar_raiser_question - L4→L5 Challenge
"Five servers have these logs (showing [index:term]):
- S1: [1:1, 2:1, 3:2, 4:2, 5:3]
- S2: [1:1, 2:1, 3:2, 4:2, 5:3]
- S3: [1:1, 2:1, 3:2, 4:2, 5:3]
- S4: [1:1, 2:1, 3:4]
- S5: [1:1, 2:1]

Entries through index 4 are committed. S3 crashes. Can S4 become leader in term 5? Can S5? Explain the voting patterns."

### bar_raiser_concepts
- S4 can win: Its log (last entry term 4) beats S1,S2 (term 3)
- S4 gets votes from itself, S5, possibly S1 and S2
- S5 cannot win: Would need 3 votes, but S1,S2,S4 all refuse
- Shows committed entry at index 4 is safe despite S4 not having it
- Demonstrates subtle case where shorter log with higher term wins

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 4-5 min answer + 4-5 min discussion
- **Common next topics**: Safety proof, liveness properties, reconfiguration
