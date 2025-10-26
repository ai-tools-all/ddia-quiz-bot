---
id: fault-tolerance-subjective-L3-001
type: subjective
level: L3
category: baseline
topic: fault-tolerance
subtopic: split-brain
estimated_time: 5-7 minutes
---

# question_title - Understanding Split-Brain Problem

## main_question - Core Question
"Explain what the split-brain problem is in distributed systems. Why is it dangerous and how does majority voting help prevent it?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Split-Brain Definition**: Multiple servers think they're the leader/primary simultaneously
- **Data Inconsistency Risk**: Different parts of system may have conflicting state
- **Majority Voting Solution**: Requires more than half of servers to agree

### expected_keywords
- Primary keywords: split-brain, leader, primary, majority, quorum
- Technical terms: network partition, consensus, voting

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Odd Number Requirement**: Why we typically use 3, 5, 7 servers
- **Network Partitions**: How network failures trigger split-brain scenarios  
- **Client Confusion**: Multiple primaries accepting different writes
- **Real-world Examples**: Database clusters, distributed locks

### bonus_keywords
- Implementation: Raft, Paxos, ZAB
- Related: CAP theorem, partition tolerance
- Trade-offs: availability vs consistency

## sample_excellent - Example Excellence
"Split-brain occurs when a distributed system partitions and multiple nodes believe they are the leader or primary, leading to divergent state. This is dangerous because clients might write to different 'primaries', causing permanent data inconsistency. For example, in a replicated database, one primary might accept an account debit while another accepts a credit, violating invariants. Majority voting prevents this by requiring more than half the servers (a quorum) to elect a leader. Since there can only be one majority in any partition, at most one side can form a functioning cluster. This is why we use odd numbers like 3 or 5 servers - it ensures clear majorities and prevents ties."

## sample_acceptable - Minimum Acceptable
"Split-brain is when multiple servers think they're in charge at the same time, usually due to network problems. This causes data inconsistency as different parts accept different updates. Majority voting fixes this by requiring more than half the servers to agree on who's leader."

## common_mistakes - Watch Out For
- Confusing split-brain with simple leader failure
- Not understanding why majority prevents multiple leaders
- Thinking any voting scheme works (not just majority)
- Missing the connection to network partitions

## follow_up_excellent - Depth Probe
**Question**: "Consider a 5-server cluster that partitions into groups of 3 and 2. Walk through what happens when both sides try to elect a leader."
- **Looking for**: Only 3-server side can form quorum, 2-server side blocks
- **Red flags**: Thinking both sides could operate

## follow_up_partial - Guided Probe  
**Question**: "You mentioned majority voting. Why specifically more than half? What would happen if we just required any two servers to agree in a 3-server system?"
- **Hint embedded**: Multiple groups of 2 possible
- **Concept testing**: Understanding quorum uniqueness

## follow_up_weak - Foundation Check
**Question**: "Imagine you have two bank ATMs that can't talk to each other. What problems could occur if both process withdrawals?"
- **Simplification**: Overdraft scenario
- **Building block**: Why coordination matters

## bar_raiser_question - L3→L4 Challenge
"A company has 6 servers in their cluster. During a network partition, it splits into two groups of 3 servers each. Can either side elect a leader? What are the implications for availability?"

### bar_raiser_concepts
- Even number of servers creates tie scenarios
- Neither partition has majority (>3 servers needed)
- System becomes unavailable despite having running servers
- Trade-off between fault tolerance and availability

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 2-3 min answer + 3-4 min discussion
- **Common next topics**: Leader election, consensus protocols, Byzantine failures

## assistant_answer
Split-brain is when multiple nodes simultaneously believe they are leader (often due to a network partition), leading to conflicting writes and inconsistent state. Majority quorums ensure only one partition can elect a leader; the minority cannot reach quorum and must reject writes, preserving safety.

## improvement_suggestions
- Add a concrete 5-node partition scenario and require explicit quorum math in the answer.
- Contrast majority voting with “two nodes agree” in a 3-node cluster to highlight quorum uniqueness and why majority is required.
