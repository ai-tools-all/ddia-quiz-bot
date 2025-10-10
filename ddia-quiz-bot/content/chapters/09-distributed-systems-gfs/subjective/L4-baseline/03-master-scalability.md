---
id: gfs-subjective-L4-003
type: subjective
level: L4
category: baseline
topic: gfs
subtopic: architecture-scalability
estimated_time: 7-10 minutes
---

# question_title - Single Master Architecture Analysis

## main_question - Core Question
"GFS uses a single master design while most distributed systems avoid single points of failure. Analyze why this worked for Google and when it would break down."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Simplicity Benefit**: Easier consistency and coordination
- **Metadata in Memory**: Fast operations, size limitations
- **Minimal Data Path**: Master not involved in data transfer
- **Scale Limits**: File count and operation rate boundaries

### expected_keywords
- Primary keywords: single master, metadata, scalability, bottleneck
- Technical terms: namespace, chunk mapping, operation log

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Shadow Masters**: Read scaling solution
- **64MB Chunks**: Reduces metadata pressure
- **Client Caching**: Reduces master load
- **Operational Benefits**: Easier debugging and reasoning
- **Evolution**: Colossus multi-master architecture

### bonus_keywords
- Implementation: B-tree namespace, checkpoint/replay
- Limits: 100 million files, thousands of operations/sec
- Hardware: 64GB RAM limitations circa 2003

## sample_excellent - Example Excellence
"The single master design was a brilliant simplification that worked because: 1) Simplicity - no distributed consensus needed for metadata operations, making the system easier to build and debug. 2) All metadata fits in memory (64GB could handle ~100 million files with 64MB chunks) enabling microsecond operations. 3) Master only handles metadata, not data flow, so bandwidth isn't bottlenecked. 4) Client caching and large chunks minimize master interactions. This breaks down when: file count exceeds memory capacity, operation rate exceeds single CPU capacity (~10K ops/sec), or when geographic distribution requires regional masters. Google hit these limits around 2010, leading to Colossus with distributed masters. The design traded future scalability for immediate simplicity and time-to-market - a good trade-off for 2003 Google."

## sample_acceptable - Minimum Acceptable
"Single master simplified the design and worked because metadata fits in memory and the master doesn't handle actual data. It breaks down when you have too many files to fit in memory or too many operations for one CPU to handle."

## common_mistakes - Watch Out For
- Calling it a "bad design" without context
- Not understanding memory vs CPU limitations
- Missing the simplicity benefits
- Ignoring the time period context (2003)

## follow_up_excellent - Depth Probe
**Question**: "How would you gradually migrate from GFS's single master to a multi-master system without downtime?"
- **Looking for**: Sharding strategy, migration phases, consistency during transition
- **Red flags**: Big-bang replacement approach

## follow_up_partial - Guided Probe
**Question**: "You mentioned memory limits. Let's calculate: if each file's metadata is 64 bytes and you have 64GB RAM, how many files can you support?"
- **Hint embedded**: Simple division with overhead
- **Concept testing**: Scale comprehension

## follow_up_weak - Foundation Check
**Question**: "Think about a library with one librarian vs multiple librarians. What are the trade-offs?"
- **Simplification**: Physical world analogy
- **Building block**: Centralization trade-offs

## bar_raiser_question - L4→L5 Challenge
"Design a hybrid architecture that keeps GFS's single master simplicity for small deployments but can scale to multiple masters as needed. Include migration strategy."

### bar_raiser_concepts
- Namespace sharding strategies
- Consistent hashing for file distribution
- Cross-shard operations handling
- Backwards compatibility requirements
- Operational complexity growth

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 3-4 min answer + 4-5 min discussion
- **Common next topics**: HDFS NameNode Federation, Ceph, distributed metadata systems
