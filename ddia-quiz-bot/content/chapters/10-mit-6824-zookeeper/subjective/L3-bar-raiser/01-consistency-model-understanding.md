---
id: zookeeper-subjective-L3-bar-raiser-001
type: subjective
level: L3
category: bar-raiser
topic: zookeeper
subtopic: consistency-model
estimated_time: 7-10 minutes
---

# question_title - Zookeeper's Mixed Consistency Model

## main_question - Core Question
"Zookeeper provides linearizable writes but not linearizable reads by default. Walk through a scenario where this could cause confusion for a developer, and explain how Zookeeper provides tools to handle this situation."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Linearizable Writes**: All writes appear in a single global order
- **Non-Linearizable Reads**: Reads may return stale data for performance
- **Sync Operation**: Forces a client to see latest state
- **Practical Impact**: Can read old values even after writes complete

### expected_keywords
- Primary keywords: linearizable, writes, reads, stale, sync
- Technical terms: consistency model, performance trade-off, sync+read

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Performance Reasoning**: Why reads from local replicas improve throughput
- **FIFO Interaction**: How client ordering helps despite stale reads
- **Common Patterns**: Read-your-writes is guaranteed due to FIFO
- **Monotonic Reads**: Client won't go backwards in time
- **Leader Routing**: Writes go through leader, reads from any replica

### bonus_keywords
- Implementation: local replicas, read scalability, quorum
- Patterns: sync+read, write-through, eventual consistency window
- Trade-offs: consistency vs availability vs performance

## sample_excellent - Example Excellence
"Scenario: Service A writes a new configuration value to Zookeeper, then tells Service B via external message to read the config. Service B might read the old value because reads aren't linearizable - they can be served from any replica that might not have the latest write yet. This breaks the intuitive expectation that 'after a write completes, everyone sees it.'

Zookeeper handles this through the sync() operation. Service B should call sync() then read(), which guarantees seeing at least all writes that completed before sync() was called. The sync() doesn't return data but ensures the client's server has caught up with the leader.

This design is intentional: reads scale linearly with replicas since they don't need coordination, while writes go through the leader for ordering. For many use cases like reading slowly-changing configuration, slightly stale data is acceptable and the performance gain is significant. When fresh data is critical (like checking lock ownership), sync+read provides the stronger guarantee.

The FIFO guarantee helps here too - a single client always sees its own writes and won't go backward in time, which handles many common patterns without needing sync()."

## sample_acceptable - Minimum Acceptable
"A client might read stale data because reads come from local replicas while writes go through the leader. If Client A writes and Client B immediately reads, B might see old data. To fix this, Client B can use sync() before reading, which makes sure it sees all completed writes."

## common_mistakes - Watch Out For
- Not explaining WHY this design choice was made
- Forgetting about the sync() solution
- Confusing with eventual consistency (it's stronger)
- Not giving a concrete scenario

## follow_up_excellent - Depth Probe
**Question**: "How would you design an API that hides this complexity from developers while still maintaining good performance for both critical and non-critical reads?"
- **Looking for**: Critical flag parameter, automatic sync for certain patterns, cost awareness
- **Red flags**: Always using sync (defeats performance purpose)

## follow_up_partial - Guided Probe  
**Question**: "In your scenario, what if Service A needs to be absolutely sure Service B will see the new value? How could they coordinate?"
- **Hint embedded**: Service A could pass version number or timestamp
- **Concept testing**: Understanding coordination patterns

## follow_up_weak - Foundation Check
**Question**: "Let's think about a news website. Would it matter if different users see articles published a few seconds apart? When would it matter?"
- **Simplification**: Consistency requirements vary by use case
- **Building block**: Trade-off understanding

## bar_raiser_question - L3→L4 Challenge
"Design a leader election system using Zookeeper where followers need to immediately recognize a new leader. How do you ensure no follower ever acts on commands from an old leader due to stale reads? Consider both correctness and performance."

### bar_raiser_concepts
- Sequential znodes for leader election
- Watches for immediate notification
- Version numbers or epochs in commands
- Sync usage for critical decisions
- Performance optimization strategies

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 3-4 min answer + 4-6 min discussion
- **Common next topics**: CAP theorem, eventual consistency, read replicas, CDNs
