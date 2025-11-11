---
id: memcached-subjective-L7-001
type: subjective
level: L7
category: baseline
topic: memcached
subtopic: architectural-evolution
estimated_time: 15-20 minutes
---

# question_title - Architectural Evolution Philosophy

## main_question - Core Question
"Facebook's memcached architecture evolved from single-server setups through sharded databases to aggressive caching absorbing 99%+ of reads. Analyze this evolution as a case study in architectural decision-making at scale. What second-order effects emerge from pushing consistency to the application layer? How would you architect a caching layer for a greenfield system today, and what lessons from Facebook's approach would you keep versus discard given modern distributed systems primitives (CRDT, Raft, distributed transactions)? Consider operational complexity, team cognitive load, and long-term maintainability."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Evolution driven by read/write asymmetry and economic constraints
- Pushing consistency to application layer increases flexibility but complexity
- Second-order effects: front-end bugs cause inconsistency, hard to debug, operational burden
- Modern primitives (CRDT, Raft, distributed transactions) change trade-offs
- Operational complexity vs performance gains
- Team expertise and maintainability considerations

### expected_keywords
- application-layer consistency, second-order effects, operational complexity, modern primitives, maintainability, cognitive load, trade-off evolution

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- When to use strongly consistent systems vs eventual consistency
- Cost of distributed transactions at scale
- CRDT applicability to social media data models
- Industry shift toward managed services vs building primitives
- Developer productivity impact
- Observability and debugging challenges
- Alternative architectures: CQRS, event sourcing, change data capture
- Coordination avoidance techniques

### bonus_keywords
- CRDT, CQRS, event sourcing, CDC, coordination avoidance, managed services, developer experience, observability

## sample_excellent - Example Excellence
"Facebook's evolution reflects classic distributed systems trade-offs: read/write asymmetry (99%+ reads) economically justified aggressive caching and eventual consistency, accepting complexity in exchange for scale. Pushing consistency to the application layer (front-end manages cache-database relationship) maximizes flexibility—caching arbitrary transformations, custom invalidation logic—but creates second-order effects: application bugs cause inconsistency hard to debug; every engineering team must understand distributed systems semantics; operational burden increases (monitoring, alerting, incident response). Today, I'd reconsider this architecture. Modern primitives shift trade-offs: Raft/Paxos-based systems (etcd, Spanner) provide strong consistency at acceptable scale; CRDTs enable coordination-free eventual consistency with guaranteed convergence; distributed transactions (Spanner, CockroachDB) simplify application logic. For greenfield systems, I'd recommend: (1) Use managed strongly consistent stores (Spanner-class) for critical data, accepting latency cost. (2) Apply CRDT where applicable (counters, sets, collaborative editing). (3) Reserve application-layer consistency only where necessary flexibility justifies complexity. (4) Invest heavily in observability—distributed tracing, consistency monitoring, anomaly detection. Keep from Facebook: regional replication for latency, tiered consistency per data type, aggressive caching for read-heavy workloads. Discard: application-managed cache coherence (too error-prone), dual invalidation complexity. Prioritize developer productivity and cognitive load—most teams shouldn't need deep distributed systems expertise. The industry trend toward managed services reflects this: operational complexity eventually dominates performance gains."

## sample_acceptable - Minimum Acceptable
"Facebook's approach made sense for read-heavy workload and their scale, but pushing consistency to apps is complex. Modern systems like Spanner offer strong consistency with acceptable performance. I'd use managed services for critical data and keep app-level consistency only where needed. Operational complexity should drive decisions today."

## common_mistakes - Watch Out For
- Not addressing second-order effects of architectural choices
- Missing modern alternatives to Facebook's approach
- Ignoring operational/maintainability dimensions
- Not considering team expertise and cognitive load
- Treating performance as the only optimization target
- Assuming Facebook's constraints apply universally

## follow_up_excellent - Depth Probe
**Question**: "You're advising a startup building a social platform. They cite Facebook's architecture as a model. What questions would you ask to determine if this approach suits their context, and what alternative architectures would you propose?"
- **Looking for**: Scale assessment, team expertise, consistency requirements, read/write ratio, budget constraints, time-to-market, build vs buy decision framework, evolutionary architecture approach

## follow_up_partial - Guided Probe
**Question**: "How do modern consistency models (CRDT, Raft, distributed transactions) change the trade-off space compared to when Facebook built their system?"
- **Hint embedded**: Technology advancement makes strong consistency more accessible

## follow_up_weak - Foundation Check
**Question**: "What are second-order effects in system design?"
