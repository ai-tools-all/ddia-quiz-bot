---
id: craq-subjective-L6-004
type: subjective
level: L6
category: baseline
topic: craq
subtopic: hybrid-consistency
estimated_time: 9-12 minutes
---

# question_title - Designing a Hybrid Consistency Tier with CRAQ

## main_question - Core Question
"Propose a hybrid system that uses CRAQ for a strongly consistent core and a Dynamo-style eventually consistent edge cache to serve bursty workloads. Detail how you keep the two layers coherent using DDIA's guidance on consistency levels and conflict resolution." 

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Separation of Concerns**: CRAQ handles writes/critical reads; edge cache serves speculative reads
- **Invalidation/Repair**: Tail commit triggers cache invalidation or update, preventing stale responses
- **Conflict Resolution Strategy**: Use CRAQ as source of truth; edges reconcile via timestamps or version vectors
- **Operational Flow**: Describe read routing, fallback behaviour, and failure handling

### expected_keywords
- Primary keywords: hybrid consistency, invalidation, reconciliation, edge cache
- Technical terms: TTL, version vector, conflict resolution, monotonic reads

## peripheral_concepts - Nice to Have (40%)
- **Traffic Segmentation**: Identify operations safe for eventual layer
- **Monitoring**: Track divergence rate, reconciliation backlog
- **Cost Control**: Manage cache replication vs CRAQ chain count
- **Resilience**: Behaviour during cache partition or CRAQ degradation

### bonus_keywords
- Implementation: CDN edge, global load balancer, change feed, reconciliation worker
- Scenarios: read-heavy analytics, personalized timelines, partial connectivity
- Trade-offs: latency vs correctness, complexity of dual path

## sample_excellent - Example Excellence
"We route critical writes and read-after-write traffic directly to CRAQ, ensuring linearizability. For high-volume reads, we populate an eventually consistent edge cache fed by CRAQ's tail commit stream; each entry carries the CRAQ epoch and version. Edges serve cached data within SLA but invalidate on tail updates. If a cache serves stale data, reconciliation compares version vectors and refreshes from CRAQ. This fusion aligns with DDIA's advice to combine strong cores with eventual edges when workloads demand low latency with tolerable temporary staleness." 

## sample_acceptable - Minimum Acceptable
"Use CRAQ as the strong core and an eventually consistent cache for less critical reads. Invalidate cache entries on tail commit using version metadata so cached data doesn't drift too far, following DDIA's hybrid consistency suggestions." 

## common_mistakes - Watch Out For
- Letting cache writes bypass CRAQ without reconciliation
- No invalidation triggered by tail commits
- Ignoring conflict resolution details
- Not referencing DDIA's hybrid strategies

## follow_up_excellent - Depth Probe
**Question**: "How do you keep read-after-write semantics for a single user when using the edge cache?"
- **Looking for**: Session-aware routing, cache bypass for fresh writes, TTL adjustments
- **Red flags**: Relying on eventual consistency alone

## follow_up_partial - Guided Probe  
**Question**: "What metrics show the cache is drifting too far from CRAQ?"
- **Hint embedded**: Version delta, stale-hit rate, reconciliation backlog
- **Concept testing**: Monitoring comprehension

## follow_up_weak - Foundation Check
**Question**: "Why keep an official record book even if you post quick updates on social media?"
- **Simplification**: Strong core vs eventual edge
- **Building block**: Authoritative source of truth

## bar_raiser_question - L6→L7 Challenge
"Design a dynamic policy engine that routes requests between CRAQ and the eventual layer based on user tier, SLA, and detected divergence." 

### bar_raiser_concepts
- Policy-driven routing
- SLA tiering
- Divergence-aware control loops
- Multi-layer consistency governance

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 5-6 min answer + 4-6 min discussion
- **Common next topics**: Policy engines, consistency tiering, dynamic routing
