---
id: primary-backup-subjective-L4-002
type: subjective
level: L4
category: baseline
topic: primary-backup
subtopic: output-commit
estimated_time: 6-8 minutes
---

# question_title - The Output Commit Problem

## main_question - Core Question
"Explain the output commit problem in primary-backup replication. Why must the primary delay sending output to external clients, and what happens if it doesn't? How does this impact system performance?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Output Visibility**: External outputs make state visible to clients
- **Backup Synchronization**: Must ensure backup has received log entries
- **Consistency After Failover**: Client shouldn't see state rollback
- **Performance Trade-off**: Latency added for every output

### expected_keywords
- Primary keywords: output commit, external visibility, synchronization
- Technical terms: network packets, disk writes, screen output, latency

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Output Rule**: Only one replica sends output (primary)
- **Buffering Mechanism**: Queue outputs until backup confirms
- **Network Optimizations**: Batching, pipelining possibilities
- **Split-Brain Prevention**: Ensures single consistent view

### bonus_keywords
- Implementation: output buffer, acknowledgment protocol, TCP responses
- Scenarios: client requests, database commits, user interfaces
- Performance: round-trip time, throughput impact, batching strategies

## sample_excellent - Example Excellence
"The output commit problem arises because external outputs (network packets to clients, disk writes, screen updates) make the VM's state visible outside the replication system. If the primary sends output immediately but crashes before the backup receives the corresponding log entries, the backup won't be able to reproduce that output after failover. From the client's perspective, they would see the system 'go backwards' - they received a response that the new primary doesn't know about. VMware FT solves this by delaying output until the backup acknowledges receiving all log entries up to that output point. This ensures the backup can reproduce the exact same state if it takes over. The cost is added latency on every output - potentially doubling response time since we wait for primary→backup→primary communication. This is especially painful for network-intensive applications with many small requests."

## sample_acceptable - Minimum Acceptable
"The output commit problem is that if the primary sends output to clients but crashes before the backup gets the log entries, the backup can't recreate that output. Clients would see responses disappear after failover. The primary must delay sending any output until it knows the backup has received all relevant log entries. This adds latency to every client response but ensures consistency."

## common_mistakes - Watch Out For
- Thinking both replicas send output
- Not understanding state visibility concept
- Missing the client perspective of inconsistency
- Underestimating performance impact

## follow_up_excellent - Depth Probe
**Question**: "Could you optimize VMware FT by batching multiple outputs together before waiting for backup acknowledgment? What are the trade-offs?"
- **Looking for**: Latency vs throughput, application semantics, buffering limits
- **Red flags**: Not considering individual request latency

## follow_up_partial - Guided Probe  
**Question**: "A client sends a 'transfer $100' request. The primary processes it and crashes right after sending 'success' but before backup gets the log. What happens from the client's view?"
- **Hint embedded**: State inconsistency visible to client
- **Concept testing**: Understanding rollback implications

## follow_up_weak - Foundation Check
**Question**: "Imagine texting a friend. You send 'yes' to their question, but your phone dies. They got the message but when your phone restarts, it shows the message was never sent. Is this confusing?"
- **Simplification**: Real-world analogy to output commit
- **Building block**: External visibility consequences

## bar_raiser_question - L4→L5 Challenge
"A replicated web server handles 1000 requests/second. Network latency between primary and backup is 1ms. Calculate the performance impact of output commit. How would you modify the design for a read-heavy vs write-heavy workload?"

### bar_raiser_concepts
- Quantitative performance analysis
- Read vs write optimization strategies
- Asynchronous possibilities for reads
- Consistency vs performance trade-offs
- Application-aware optimizations

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 3-4 min answer + 3-4 min discussion
- **Common next topics**: Performance optimization, async replication, eventual consistency
