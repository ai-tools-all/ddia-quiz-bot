---
id: primary-backup-subjective-L4-BR-001
type: subjective
level: L4
category: bar-raiser
topic: primary-backup
subtopic: performance-analysis
estimated_time: 8-10 minutes
---

# question_title - Debugging Performance in Primary-Backup Systems

## main_question - Core Question
"A VMware FT protected VM is experiencing 3x slower response times than an unprotected VM running the same application. The application is a Redis cache handling 50,000 small GET/SET operations per second. Diagnose the likely causes and propose optimizations while maintaining fault tolerance guarantees."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Output Commit Overhead**: Each response waits for backup acknowledgment
- **Network RTT Impact**: Latency between primary and backup
- **Small Operation Amplification**: Overhead dominates for tiny operations
- **Logging Bandwidth**: High operation rate generates significant log traffic

### expected_keywords
- Primary keywords: output commit, network latency, logging overhead, synchronization
- Technical terms: RTT, bandwidth saturation, interrupt storms, batching

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **CPU Cache Effects**: Logging disrupts cache locality
- **Interrupt Coalescence**: High packet rate optimization
- **NUMA Effects**: Memory locality in modern systems
- **Network Offload**: TCP segmentation, checksum offloading
- **Application Patterns**: Read vs write ratio implications

### bonus_keywords
- Optimizations: batching, pipelining, async reads, log compression
- Monitoring: network utilization, CPU overhead, disk I/O
- Architecture: placement strategies, dedicated network, SSD logging

## sample_excellent - Example Excellence
"The 3x slowdown in Redis with VMware FT is primarily due to output commit overhead amplified by small, frequent operations. Each of the 50,000 ops/second requires: (1) Output commit delay - waiting for backup acknowledgment adds network RTT to every response, potentially 1-2ms per operation, (2) Logging overhead - each operation generates log entries, consuming bandwidth and CPU for 50K logs/second, (3) Interrupt handling - network packets for logging trigger interrupts, causing CPU overhead and cache pollution. 

Optimizations while maintaining fault tolerance: (1) Batch operations at application level - combine multiple Redis operations to amortize output commit overhead, (2) Use dedicated high-speed network for primary-backup communication to minimize RTT, (3) Place primary and backup in same rack/switch to reduce latency, (4) Enable interrupt coalescing to handle packet storms efficiently, (5) Consider async replication for read-heavy workloads where slight staleness acceptable, (6) Use SSD for log buffer to prevent disk bottleneck, (7) Implement log compression since Redis commands are text-based and highly compressible. 

The fundamental issue is that VMware FT's strong consistency guarantees impose per-operation overhead that becomes dominant for high-frequency small operations."

## sample_acceptable - Minimum Acceptable
"The slowdown is because VMware FT adds output commit delay to every Redis operation. With 50,000 ops/second, each one waits for the backup to acknowledge receiving log entries before responding. If network latency is 1ms, that adds 1ms to every operation. Also, logging overhead for 50,000 operations creates significant network traffic and CPU usage. To optimize: reduce network latency between primary and backup, batch multiple operations together, use faster network connections, and consider if all operations need synchronous replication."

## common_mistakes - Watch Out For
- Missing output commit as primary cause
- Not recognizing small operation amplification
- Suggesting to remove fault tolerance entirely
- Ignoring network bandwidth limits
- Not considering application-level optimizations

## follow_up_excellent - Depth Probe
**Question**: "Redis supports pipelining where clients send multiple commands without waiting for responses. How would this affect VMware FT's overhead, and what new challenges might arise?"
- **Looking for**: Batching benefits, output ordering, memory pressure
- **Red flags**: Not understanding pipeline vs individual operation handling

## follow_up_partial - Guided Probe  
**Question**: "If we measure that network RTT is 0.5ms and Redis normally responds in 0.1ms, what's the theoretical minimum response time with FT?"
- **Hint embedded**: Output commit adds full RTT
- **Concept testing**: Understanding synchronous overhead

## follow_up_weak - Foundation Check
**Question**: "If you had to wait for your friend to write down every text message you send before sending another, how would this affect your texting speed?"
- **Simplification**: Synchronous acknowledgment analogy
- **Building block**: Understanding blocking operations

## bar_raiser_question - L4→L5 Challenge
"Design a hybrid replication strategy for Redis where critical operations (financial transactions) get full VMware FT protection, while cache warming reads use asynchronous replication. Detail the implementation challenges and consistency guarantees."

### bar_raiser_concepts
- Operation classification mechanisms
- Mixed consistency models
- State divergence handling
- Client API changes
- Failover complexity
- Performance monitoring strategies

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 4-5 min answer + 4-5 min discussion
- **Common next topics**: Application-aware replication, eventual consistency, CRDT
