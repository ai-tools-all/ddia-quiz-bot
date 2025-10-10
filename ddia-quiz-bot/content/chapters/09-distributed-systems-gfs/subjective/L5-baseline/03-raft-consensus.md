---
id: raft-subjective-L5-001
type: subjective
level: L5
category: baseline
topic: raft-consensus
subtopic: leader-election
estimated_time: 10-12 minutes
---

# question_title - Raft Leader Election and Split Brain Prevention

## main_question - Core Question
"In Raft, how does the leader election mechanism prevent split brain scenarios? Walk me through what happens when a network partition occurs and how the system maintains consistency."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Majority Voting**: Need majority (n/2 + 1) to elect leader
- **Term Numbers**: Monotonically increasing logical clock
- **Vote Restriction**: One vote per server per term
- **Network Partition Handling**: Minority partition cannot elect leader

### expected_keywords
- Primary keywords: majority, quorum, term, vote, candidate, follower
- Technical terms: split brain, partition tolerance, election timeout

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Randomized Timeouts**: Reduces election conflicts
- **Pre-vote Optimization**: Prevents disruption from isolated nodes
- **RequestVote RPC**: How votes are requested and granted
- **Leader Lease**: Implicit leadership duration
- **Step-down Mechanism**: Higher term causes immediate step-down
- **Log Completeness**: Vote restriction based on log comparison

### bonus_keywords
- Related concepts: Paxos comparison, Byzantine failures, CAP theorem
- Optimizations: Pre-vote, leadership transfer, learner nodes
- Implementation details: RPC retries, persistent state, volatile state

## sample_excellent - Example Excellence
"Raft prevents split brain through majority-based leader election. Key mechanism: A candidate needs votes from a majority (n/2 + 1) of servers to become leader. In a 5-server cluster, need 3 votes. 

During network partition: If network splits 3-2, only the 3-node partition can elect a leader since it has majority. The 2-node partition's candidates will never get 3 votes, preventing split brain. 

Term numbers act as logical clocks - each election attempt increments the term. Servers reject messages from lower terms and step down if they see higher terms. This ensures at most one leader per term.

The voting process: Each server votes once per term (persisted to disk). Candidates vote for themselves, then request votes from others. The RequestVote RPC includes candidate's term and log information. Voters check: 1) Haven't voted in this term already, 2) Candidate's log is at least as up-to-date.

This combination - majority requirement, unique terms, single vote per term - mathematically guarantees at most one leader per term, completely preventing split brain even during arbitrary network failures."

## sample_acceptable - Minimum Acceptable
"Raft prevents split brain by requiring a majority of servers to elect a leader. In a 5-server cluster, a candidate needs 3 votes to become leader. Since there's only one majority possible, only one partition can have a leader during network splits. Term numbers ensure we know which leader is most recent - higher terms override lower terms."

## common_mistakes - Watch Out For
- Confusing majority with any quorum
- Not mentioning term numbers' role
- Ignoring vote persistence requirement
- Missing the "at most one leader per term" guarantee

## follow_up_excellent - Depth Probe
**Question**: "How would Raft handle a scenario where a leader gets temporarily disconnected, a new leader is elected, then the old leader reconnects?"
- **Looking for**: Term comparison, automatic step-down, log reconciliation
- **Red flags**: Not understanding term-based conflict resolution

## follow_up_partial - Guided Probe
**Question**: "You mentioned majority voting. What happens if we have an even number of servers, like 4?"
- **Hint embedded**: Still need 3 votes, reduced fault tolerance
- **Concept testing**: Understanding majority math and fault tolerance

## follow_up_weak - Foundation Check
**Question**: "If we have 3 servers and one fails, can the remaining 2 elect a leader?"
- **Simplification**: Basic majority calculation
- **Building block**: Understanding quorum requirements

## bar_raiser_question - L5→L6 Challenge
"Design a modification to Raft that allows read operations to be served by followers while still guaranteeing linearizability. What additional mechanisms would you need?"

### bar_raiser_concepts
- Read leases and clock synchronization
- Leader heartbeat tracking
- Linearizable read options
- Performance vs consistency trade-offs
- Follower read indices

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 4-5 min answer + 5-7 min discussion
- **Common next topics**: Log replication, Byzantine consensus, Paxos comparison, Multi-Raft
