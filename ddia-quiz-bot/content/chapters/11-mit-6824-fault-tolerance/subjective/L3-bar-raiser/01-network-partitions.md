---
id: fault-tolerance-subjective-L3-004
type: subjective
level: L3
category: bar-raiser
topic: fault-tolerance
subtopic: network-partitions
estimated_time: 7-10 minutes
---

# question_title - Handling Network Partitions

## main_question - Core Question
"A 5-server Raft cluster gets partitioned: servers A, B, C in one partition and D, E in another. The original leader was server A. Describe what happens in each partition and what occurs when the network heals."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Majority Partition Continues**: A, B, C can still form quorum
- **Minority Partition Stops**: D, E cannot elect leader or serve requests
- **Term Number Resolution**: Higher term wins when partition heals

### expected_keywords
- Primary keywords: partition, majority, minority, quorum
- Technical terms: term number, election, unavailable

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Client Experience**: Requests to D, E fail/timeout
- **Log Divergence**: D, E logs may be stale after healing
- **Automatic Recovery**: No manual intervention needed
- **Safety Preserved**: No committed data lost

### bonus_keywords
- Implementation: heartbeat failures, election timeouts
- Related: CAP theorem, availability trade-offs
- Recovery: log reconciliation, catching up

## sample_excellent - Example Excellence
"When the partition occurs, the majority side (A, B, C) continues operating normally. Server A remains leader as it can still reach a majority for heartbeats and log replication. Client requests to this partition succeed. The minority side (D, E) quickly discovers they've lost the leader when heartbeats timeout. They'll attempt elections but fail to gather majority votes (need 3, can only get 2), so they remain in candidate/follower state, unable to process client requests. This preserves safety over availability. When the network heals, D and E will receive AppendEntries from A (still leader with likely higher term). They'll recognize A's authority, revert to followers, and A will repair their logs by finding the divergence point and overwriting any conflicting entries. The system automatically recovers to a consistent state with all 5 servers having identical committed logs."

## sample_acceptable - Minimum Acceptable
"The 3-server partition with the leader continues working since it has a majority. The 2-server partition cannot elect a new leader since they need 3 votes. When the network reconnects, the servers in the smaller partition rejoin and accept the leader from the majority partition, updating their logs to match."

## common_mistakes - Watch Out For
- Thinking minority partition can operate in degraded mode
- Not mentioning term numbers in resolution
- Assuming manual intervention needed
- Missing that committed entries are preserved

## follow_up_excellent - Depth Probe
**Question**: "What if during the partition, servers D and E had somehow elected a leader with term 10, while A only advanced to term 8? What happens during healing?"
- **Looking for**: D/E leader steps down, but scenario impossible due to majority requirement
- **Red flags**: Not recognizing this violates Raft's safety

## follow_up_partial - Guided Probe  
**Question**: "You mentioned the minority partition can't serve requests. Why is this preferable to allowing both partitions to operate?"
- **Hint embedded**: Data consistency importance
- **Concept testing**: CAP theorem understanding

## follow_up_weak - Foundation Check
**Question**: "Imagine a classroom splits into two groups that can't communicate. If both groups made different decisions about a project, what problems would arise when they reunite?"
- **Simplification**: Conflicting decisions
- **Building block**: Why single authority matters

## bar_raiser_question - L3→L4 Challenge
"Consider a 7-server cluster that partitions into three groups: [A,B,C], [D,E], and [F,G]. The original leader was in the first group. What is the availability of the system? What if it instead partitioned into [A,B,C,D] and [E,F,G]?"

### bar_raiser_concepts
- First scenario: Only [A,B,C] lacks majority (3 < 4), system unavailable
- Second scenario: [A,B,C,D] has majority, remains available
- Demonstrates partition size vs majority calculation
- Shows how partition patterns affect availability

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 3-4 min answer + 4-6 min discussion
- **Common next topics**: Byzantine failures, consensus impossibility, partition tolerance strategies
