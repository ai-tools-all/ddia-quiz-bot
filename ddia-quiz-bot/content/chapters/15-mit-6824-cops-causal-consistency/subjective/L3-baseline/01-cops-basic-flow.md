---
id: cops-subjective-L3-001
type: subjective
level: L3
category: baseline
topic: cops
subtopic: basic-flow
estimated_time: 6-8 minutes
---

# question_title - COPS Basic Operation Flow

## main_question - Core Question
"Describe the basic flow of a put operation in COPS. What information does the client send with the put, and what happens at the local shard server?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Client maintains a context of key-version pairs from previous gets
- Put operation sends key, value, AND current context as dependency list
- Local shard server immediately acknowledges the put and assigns version number
- Server asynchronously replicates put with dependencies to remote data centers
- No waiting for cross-datacenter coordination on critical path

### expected_keywords
- context, key-version pairs, dependencies, local acknowledgment, asynchronous replication

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Client can immediately proceed after local acknowledgment
- Dependencies encode "this write depends on seeing X:v2, Y:v4"
- Version numbering (e.g., Lamport clocks)

### bonus_keywords
- Lamport clock, causal chain, non-blocking write

## sample_excellent - Example Excellence
"When a client issues a put, it sends the key-value pair along with its current context—a list of key-version pairs it has read. The local shard server immediately acknowledges this put and assigns it a new version number using a Lamport clock. The client can proceed without waiting. The server then asynchronously replicates this put with its attached dependency list to corresponding shard servers at remote data centers, ensuring causal ordering is preserved."

## sample_acceptable - Minimum Acceptable
"Client sends put with its context (dependencies from previous reads). Local server acks immediately and replicates asynchronously to other data centers with the dependency information."

## common_mistakes - Watch Out For
- Thinking the client waits for remote acknowledgments
- Not mentioning that dependencies are sent with the put
- Confusing with synchronous replication models

## follow_up_excellent - Depth Probe
**Question**: "Why doesn't the client need to wait for remote data centers to acknowledge the put?"
- **Looking for**: Local-only coordination, fault tolerance, low latency goal

## follow_up_partial - Guided Probe
**Question**: "What does the dependency list contain and why is it important?"
- **Hint embedded**: Connect to causal ordering preservation

## follow_up_weak - Foundation Check
**Question**: "What is the difference between COPS put and Strawman 2's sync barrier approach?"
