---
id: fault-tolerance-subjective-L3-002
type: subjective
level: L3
category: baseline
topic: fault-tolerance
subtopic: leader-election
estimated_time: 5-7 minutes
---

# question_title - Leader Election in Raft

## main_question - Core Question
"Describe how leader election works in Raft. What triggers an election and how does a candidate become the leader?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Election Timeout**: Follower times out and becomes candidate
- **RequestVote RPC**: Candidate asks other servers for votes
- **Majority Requirement**: Need votes from majority to become leader

### expected_keywords
- Primary keywords: election, candidate, leader, follower, timeout
- Technical terms: term number, RequestVote, heartbeat

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Term Numbers**: Logical clock preventing old leaders
- **Randomized Timeouts**: Prevents simultaneous elections
- **Vote Persistence**: Each server votes only once per term
- **Split Vote Handling**: What happens when no majority achieved

### bonus_keywords
- Implementation: election timeout range (150-300ms typical)
- Related: heartbeat interval, AppendEntries
- State transitions: follower → candidate → leader

## sample_excellent - Example Excellence
"Leader election in Raft begins when a follower doesn't receive heartbeats from the current leader within its election timeout (typically 150-300ms, randomized to avoid collisions). The follower increments its term number, transitions to candidate state, votes for itself, and sends RequestVote RPCs to all other servers. Other servers grant their vote if they haven't voted in this term and the candidate's log is at least as up-to-date as theirs. A candidate becomes leader upon receiving votes from a majority of servers (including itself). If the election times out without a majority, a new election begins with a higher term. The randomized timeouts help ensure usually only one candidate emerges, avoiding split votes."

## sample_acceptable - Minimum Acceptable
"When a follower doesn't hear from the leader for a timeout period, it becomes a candidate and starts an election. It asks other servers to vote for it. If it gets votes from more than half the servers, it becomes the new leader."

## common_mistakes - Watch Out For
- Forgetting the candidate votes for itself
- Not mentioning term numbers
- Unclear about timeout triggering elections
- Missing the majority requirement

## follow_up_excellent - Depth Probe
**Question**: "What happens if two servers become candidates simultaneously and split the votes evenly? How does Raft resolve this?"
- **Looking for**: Election timeout, new term, randomized delays
- **Red flags**: Thinking deadlock is permanent

## follow_up_partial - Guided Probe  
**Question**: "You mentioned timeouts. Why are they randomized rather than fixed?"
- **Hint embedded**: Multiple candidates problematic
- **Concept testing**: Understanding collision avoidance

## follow_up_weak - Foundation Check
**Question**: "Think of a classroom where students need to pick a team leader. How would you ensure only one leader is chosen?"
- **Simplification**: Voting basics
- **Building block**: Single leader importance

## bar_raiser_question - L3→L4 Challenge
"A 5-server Raft cluster experiences a network partition: 3 servers in one partition, 2 in another. The original leader is in the 2-server partition. Walk through what happens in both partitions."

### bar_raiser_concepts
- 2-server partition: Leader steps down (no majority for heartbeats)
- 3-server partition: Election occurs, new leader emerges
- Term numbers ensure old leader can't interfere when partition heals
- Client requests in minority partition fail

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 2-3 min answer + 3-4 min discussion
- **Common next topics**: Log replication, term numbers, safety properties
