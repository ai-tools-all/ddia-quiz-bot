---
id: fault-tolerance-subjective-L3-003
type: subjective
level: L3
category: baseline
topic: fault-tolerance
subtopic: log-replication
estimated_time: 5-7 minutes
---

# question_title - Log Replication Basics

## main_question - Core Question
"Explain how Raft ensures that all servers have identical logs. What happens when a follower's log diverges from the leader's?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **AppendEntries RPC**: Leader sends log entries to followers
- **Log Matching Property**: If entries match at index, all previous match
- **Consistency Check**: Leader finds divergence point and overwrites

### expected_keywords
- Primary keywords: log, replication, AppendEntries, consistency
- Technical terms: log index, term number, commit

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Forced Overwrites**: Followers accept leader's log as truth
- **Commit Index**: Entries replicated on majority before committing
- **PrevLogIndex/Term**: Used to detect inconsistencies
- **Leader Completeness**: Committed entries preserved across leaders

### bonus_keywords
- Implementation: nextIndex, matchIndex tracking
- Related: state machine, durability
- Safety: never committed entries lost

## sample_excellent - Example Excellence
"Raft ensures identical logs through the AppendEntries RPC mechanism. The leader sends new log entries along with the index and term of the preceding entry. Followers check if their log matches at that position - if yes, they append the new entries; if no, they reject. When rejected, the leader decrements its nextIndex for that follower and retries with earlier entries, walking backward until finding the divergence point. The follower then deletes conflicting entries and replaces them with the leader's entries. This maintains the Log Matching Property: if two logs have the same entry at an index, all preceding entries are identical. Combined with only committing entries replicated on a majority, this ensures all committed entries eventually appear in every server's log in the same order."

## sample_acceptable - Minimum Acceptable
"The leader sends log entries to followers using AppendEntries messages. If a follower's log doesn't match the leader's, the leader finds where they diverged and overwrites the follower's log from that point forward. This ensures all servers end up with the same log."

## common_mistakes - Watch Out For
- Thinking followers can reject valid leader entries
- Not understanding leader's log is authoritative
- Missing the backward search for divergence
- Confusing replication with commitment

## follow_up_excellent - Depth Probe
**Question**: "Consider a follower that was disconnected and missed 100 log entries. Walk through how the leader brings it up to date efficiently."
- **Looking for**: Backward search optimization, batch sending
- **Red flags**: Thinking it needs 100 round trips

## follow_up_partial - Guided Probe  
**Question**: "You said the leader overwrites the follower's log. What if the follower had accepted some client commands while disconnected?"
- **Hint embedded**: Only leaders accept client commands
- **Concept testing**: Understanding leader authority

## follow_up_weak - Foundation Check
**Question**: "Imagine you're copying notes from the board and you missed a few lines. How would you figure out where you got lost and catch up?"
- **Simplification**: Finding sync point
- **Building block**: Sequential consistency

## bar_raiser_question - L3→L4 Challenge
"A follower has log entries [1:A, 2:B, 3:C] from term 2. The new leader (term 3) has [1:A, 2:B, 2:D, 3:E]. Walk through the AppendEntries consistency check and explain what happens to each server's log."

### bar_raiser_concepts
- PrevLogIndex=2 check fails (different terms)
- Leader retries with PrevLogIndex=1
- Follower deletes C, accepts D and E
- Final state: both have [1:A, 2:B, 2:D, 3:E]

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 2-3 min answer + 3-4 min discussion
- **Common next topics**: Commitment rules, safety properties, log compaction

## assistant_answer
The leader replicates entries via AppendEntries, including prevLogIndex/prevLogTerm to enforce the Log Matching Property. Followers reject on mismatch; the leader backs up nextIndex to find the divergence, the follower deletes conflicting suffix, then appends the leader’s entries. Majority replication precedes commitment, so committed entries appear in the same order on all servers.

## improvement_suggestions
- Include a concrete mismatch example (indices/terms) and require candidates to walk the retry/backoff sequence.
- Ask for discussion of batching/backpressure strategies and their impact on catch-up performance.
