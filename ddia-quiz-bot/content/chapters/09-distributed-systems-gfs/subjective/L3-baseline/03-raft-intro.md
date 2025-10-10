---
id: raft-subjective-L3-001
type: subjective
level: L3
category: baseline
topic: raft-consensus
subtopic: introduction
estimated_time: 6-8 minutes
---

# question_title - Introduction to Raft Consensus

## main_question - Core Question
"What problem does the Raft consensus algorithm solve in distributed systems? Explain why we need consensus and what could go wrong without it."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Agreement Problem**: Multiple servers need to agree on values
- **Fault Tolerance**: System continues despite server failures
- **Split Brain Prevention**: Avoiding conflicting decisions
- **Single Source of Truth**: One consistent view across servers

### expected_keywords
- Primary keywords: consensus, agreement, consistency, fault tolerance
- Technical terms: replicated state machine, leader, majority

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Use Cases**: Configuration management, distributed locks, leader election
- **State Machine Replication**: Same commands in same order
- **Majority Requirement**: Why need more than half
- **Network Partitions**: Handling network splits
- **Comparison with Manual Coordination**: Why algorithms better than human operators
- **Real Systems**: etcd, Consul examples

### bonus_keywords
- Examples: Distributed database, service discovery, configuration store
- Problems solved: Distributed locking, leader election, ordered logs
- Industry usage: Kubernetes (etcd), Docker Swarm, HashiCorp Consul

## sample_excellent - Example Excellence
"Raft solves the fundamental problem of getting multiple servers to agree on shared state in a distributed system, even when some servers fail or the network has problems.

**The Core Problem**: Imagine multiple database servers trying to stay synchronized. Without consensus, each server might process operations in different orders, leading to inconsistent data. Worse, during network partitions, different parts might make conflicting decisions - the 'split brain' problem where two servers both think they're the leader.

**What Raft Provides**:
1. **Consistent Ordering**: All servers process the same operations in the same order
2. **Fault Tolerance**: System continues working even if some servers fail (up to minority)
3. **Split Brain Prevention**: Mathematical guarantee of at most one leader, preventing conflicting decisions
4. **Automatic Recovery**: Failed servers can rejoin and catch up automatically

**Without Consensus**: Systems would face data inconsistency, lost updates, and conflicting operations. Manual coordination doesn't scale and is error-prone. Example: Two bank servers both approving a withdrawal from the same account during network issues.

**Real Applications**: Kubernetes uses etcd (Raft-based) to store cluster configuration. All cluster decisions go through this consensus layer, ensuring all nodes see consistent configuration even during failures."

## sample_acceptable - Minimum Acceptable
"Raft helps multiple servers agree on the same data even when some servers fail. Without it, servers might have different data or make conflicting decisions (split brain). Raft ensures all servers process operations in the same order and prevents two servers from both thinking they're in charge."

## common_mistakes - Watch Out For
- Confusing consensus with simple replication
- Not mentioning fault tolerance
- Thinking it's just about leader election
- Missing the consistency aspect

## follow_up_excellent - Depth Probe
**Question**: "How many servers would you need to tolerate 2 server failures? Why?"
- **Looking for**: Need 5 servers (2f+1 formula), majority requirement explanation
- **Red flags**: Not understanding majority math

## follow_up_partial - Guided Probe
**Question**: "You mentioned split brain. Can you give a specific example of what could go wrong?"
- **Hint embedded**: Two leaders accepting different writes
- **Concept testing**: Understanding concurrent conflict scenarios

## follow_up_weak - Foundation Check
**Question**: "If you have 3 servers and 1 fails, can the system continue working? Why or why not?"
- **Simplification**: Basic majority concept
- **Building block**: Understanding minimum viable cluster

## bar_raiser_question - L3→L4 Challenge
"If Raft guarantees consistency, why do some systems choose to use eventually consistent databases instead? When would you pick each approach?"

### bar_raiser_concepts
- CAP theorem basics
- Consistency vs availability trade-offs
- Latency implications
- Use case analysis
- Cost considerations

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 2-3 min answer + 4-5 min discussion
- **Common next topics**: CAP theorem, database replication, distributed system basics
