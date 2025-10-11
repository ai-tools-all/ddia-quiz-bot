---
id: fault-tolerance-subjective-L4-001
type: subjective
level: L4
category: baseline
topic: fault-tolerance
subtopic: term-numbers
estimated_time: 6-8 minutes
---

# question_title - Term Numbers and Safety

## main_question - Core Question
"Explain the role of term numbers in Raft. How do they prevent 'stale' leaders from corrupting the system state after recovering from a partition?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Logical Clock**: Terms act as logical timestamps for leader epochs
- **Monotonic Increase**: Terms only increase, never decrease
- **Authority Establishment**: Higher terms override lower terms
- **Vote Tracking**: Each server votes once per term

### expected_keywords
- Primary keywords: term, epoch, logical clock, stale leader
- Technical terms: RequestVote, AppendEntries, current term

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Persistent Storage**: Terms must survive crashes
- **RPC Integration**: All RPCs include sender's term
- **Immediate Updates**: Servers update term on seeing higher one
- **Leader Step-down**: Leaders revert to follower on higher terms
- **Election Trigger**: Term increment starts new election

### bonus_keywords
- Implementation: persistent state, term comparison
- Related: Lamport clocks, vector clocks, causality
- Safety: prevents multiple leaders in same term

## sample_excellent - Example Excellence
"Term numbers in Raft serve as a logical clock that divides time into epochs, each with at most one leader. Every server maintains a current term that monotonically increases. When a follower times out and starts an election, it increments its term. All RPCs include the sender's term, and recipients immediately update to any higher term they see, stepping down if they were leader. This mechanism prevents stale leaders from corrupting state: imagine a leader gets partitioned, misses several elections, then reconnects. Its AppendEntries RPCs will carry its old term number. Followers with higher terms will reject these RPCs, and the stale leader will see the higher term in responses, immediately stepping down to follower. Since servers vote only once per term and a leader needs majority votes, there can't be two leaders with the same term number. This creates a total ordering of leadership periods, ensuring that an older leader can never override decisions made by newer leaders."

## sample_acceptable - Minimum Acceptable
"Term numbers are like version numbers that increase with each election. They prevent old leaders from interfering after network problems. When a server sees a higher term number in any message, it updates its term and stops being leader if necessary. This ensures only the most recent leader can make changes."

## common_mistakes - Watch Out For
- Not explaining monotonic increase property
- Missing the connection to voting once per term
- Unclear about automatic step-down behavior
- Forgetting terms must be persistent

## follow_up_excellent - Depth Probe
**Question**: "Consider a leader with term 5 that gets partitioned. While isolated, the rest of the cluster goes through terms 6, 7, and 8. What's the minimum number of messages needed for the old leader to discover it's stale when the partition heals?"
- **Looking for**: Just one RPC response with term 8
- **Red flags**: Thinking multiple rounds needed

## follow_up_partial - Guided Probe  
**Question**: "You mentioned servers vote once per term. How does this interact with term numbers to prevent split-brain?"
- **Hint embedded**: Can't have two majorities in same term
- **Concept testing**: Connecting voting to leader uniqueness

## follow_up_weak - Foundation Check
**Question**: "Think of term numbers like version numbers on a document. If you have version 3 and someone shows you version 5, what should you do with your changes?"
- **Simplification**: Accepting newer authority
- **Building block**: Version precedence

## bar_raiser_question - L4→L5 Challenge
"A cluster experiences rapid network instability, causing 10 elections in 30 seconds, advancing from term 1 to term 11. Server S was leader in term 3 and has been partitioned since. It has uncommitted entries from term 3 in its log. When S reconnects at term 11, could those term 3 entries ever be committed? Why or why not?"

### bar_raiser_concepts
- Term 3 entries can only be committed by term 3 leader
- S is no longer term 3 leader (term ended)
- New leader might have different entries at those indices
- Demonstrates commitment rules and term boundaries
- Shows why leaders can only commit from current term

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 3-4 min answer + 3-4 min discussion
- **Common next topics**: Commitment rules, log matching property, election safety
