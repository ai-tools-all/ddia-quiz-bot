---
id: zookeeper-subjective-L3-001
type: subjective
level: L3
category: baseline
topic: zookeeper
subtopic: linearizability
estimated_time: 5-7 minutes
---

# question_title - Understanding Linearizability in Zookeeper

## main_question - Core Question
"Explain what linearizability means in the context of Zookeeper. Why is it important for a coordination service to provide this guarantee?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Linearizability Definition**: Operations appear to happen atomically at some point between start and end
- **Strong Consistency**: All clients see the same order of operations
- **Coordination Service Need**: Ensures all nodes agree on critical state (configs, leaders)

### expected_keywords
- Primary keywords: linearizability, atomic, consistency, coordination
- Technical terms: total order, real-time ordering, write operations

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Write Operations Only**: Zookeeper provides linearizability only for writes, not reads
- **Performance Trade-off**: Linearizability comes at cost of latency/throughput
- **Comparison**: Stronger than eventual consistency, sequential consistency
- **Use Cases**: Leader election, configuration management need this guarantee

### bonus_keywords
- Implementation: consensus protocol, atomic broadcast
- Related: Raft, Paxos, ZAB (Zookeeper Atomic Broadcast)
- Trade-offs: availability vs consistency, CAP theorem

## sample_excellent - Example Excellence
"Linearizability in Zookeeper means that all write operations appear to happen atomically at a single point in time between when they start and complete. This provides strong consistency - all clients observe the same sequence of state changes in the same order. For a coordination service like Zookeeper, this is critical because distributed applications rely on it for agreement on fundamental state like who is the leader, what's the current configuration, or who holds a lock. Without linearizability, you could have split-brain scenarios where different parts of the system have conflicting views of critical state. Importantly, Zookeeper only guarantees linearizability for writes - reads may return stale data for performance reasons, though clients can force linearizable reads using sync() operations."

## sample_acceptable - Minimum Acceptable
"Linearizability means operations appear to happen in a single, well-defined order that all clients agree on. For Zookeeper as a coordination service, this ensures all nodes see the same sequence of configuration changes and leader elections, preventing inconsistencies."

## common_mistakes - Watch Out For
- Confusing linearizability with serializability (database concept)
- Thinking all operations are linearizable (reads aren't by default)
- Not understanding the 'point in time' aspect
- Missing why coordination services specifically need this

## follow_up_excellent - Depth Probe
**Question**: "You mentioned Zookeeper doesn't provide linearizable reads by default. Why would they make this design choice, and how can a client work around it when needed?"
- **Looking for**: Performance benefits, sync() operation, read-your-writes consistency
- **Red flags**: Not understanding read performance implications

## follow_up_partial - Guided Probe  
**Question**: "You said operations happen in order. What happens if two clients try to create the same znode at exactly the same time?"
- **Hint embedded**: Only one can succeed atomically
- **Concept testing**: Atomic operation understanding

## follow_up_weak - Foundation Check
**Question**: "Let's simplify - imagine a shared counter that multiple people are trying to increment. What problems could occur without proper coordination?"
- **Simplification**: Race conditions, lost updates
- **Building block**: Why atomicity matters

## bar_raiser_question - L3→L4 Challenge
"Consider a distributed configuration service using Zookeeper. A configuration update is written, and immediately after, multiple services read the config. Some see the old value, some see the new value. Is this violating linearizability? Explain your reasoning."

### bar_raiser_concepts
- Linearizability applies only to writes
- Read consistency models (eventual, monotonic, read-your-writes)
- Client-side solutions (sync+read pattern)
- Trade-offs between consistency and performance

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 2-3 min answer + 3-4 min discussion
- **Common next topics**: CAP theorem, consensus protocols, Raft vs ZAB
