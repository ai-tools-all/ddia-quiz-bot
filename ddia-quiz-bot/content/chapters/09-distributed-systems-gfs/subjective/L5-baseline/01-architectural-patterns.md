---
id: gfs-subjective-L5-001
type: subjective
level: L5
category: baseline
topic: gfs
subtopic: architectural-patterns
estimated_time: 10-12 minutes
---

# question_title - GFS Architectural Patterns and Trade-offs

## main_question - Core Question
"GFS separates control plane (master) from data plane (chunk servers). This pattern appears in many distributed systems. Analyze the benefits, drawbacks, and alternatives to this architectural decision."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Separation of Concerns**: Different scaling requirements
- **Performance Isolation**: Control operations don't affect data throughput
- **Independent Failure Domains**: Master failure doesn't stop data access
- **Complexity Trade-off**: Coordination between planes

### expected_keywords
- Primary keywords: control plane, data plane, separation, scaling
- Technical terms: metadata management, coordination, failure isolation

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Industry Examples**: Kubernetes, SDN, cloud services
- **Alternative Patterns**: Peer-to-peer, symmetric systems
- **Operational Benefits**: Different upgrade cycles, specialized hardware
- **Security Implications**: Smaller attack surface for control
- **Economic Factors**: Different resource requirements
- **Evolution Patterns**: Starting simple, separating over time

### bonus_keywords
- Related systems: HDFS, Ceph, Kubernetes, OpenStack
- Patterns: Microservices, service mesh, control loops
- Trade-offs: Latency, consistency, operational complexity

## sample_excellent - Example Excellence
"Control/data plane separation is fundamental to scalable distributed systems. Benefits: 1) Independent scaling - GFS master needs CPU/memory for metadata, chunk servers need disk/network for data. 2) Performance isolation - metadata operations don't compete with data transfers. 3) Failure isolation - temporary master unavailability doesn't stop ongoing data operations. 4) Specialized optimization - master uses B-trees in memory, chunk servers optimize for sequential disk I/O. 5) Security - smaller, auditable control surface. Drawbacks: 1) Coordination complexity - keeping planes synchronized. 2) Additional latency - extra hop for metadata. 3) Consistency challenges - cache invalidation between planes. Alternatives like Ceph's distributed metadata or peer-to-peer systems (Cassandra) avoid single control points but add complexity. The pattern works best when control operations are much less frequent than data operations, which fits GFS's large file workload perfectly."

## sample_acceptable - Minimum Acceptable
"Separating control and data planes allows them to scale independently - the master scales for metadata operations while chunk servers scale for storage capacity and bandwidth. This isolation means data can continue flowing even if the control plane has issues. The downside is added complexity in keeping them coordinated."

## common_mistakes - Watch Out For
- Not recognizing this as a general pattern
- Missing the coordination complexity
- Ignoring failure isolation benefits
- Not considering alternatives

## follow_up_excellent - Depth Probe
**Question**: "How would you apply this pattern to a distributed database system? What would be in each plane?"
- **Looking for**: Query routing vs storage, schema vs data, coordinator roles
- **Red flags**: Force-fitting the pattern where it doesn't belong

## follow_up_partial - Guided Probe
**Question**: "You mentioned coordination complexity. What specific problems arise when control and data planes get out of sync?"
- **Hint embedded**: Stale metadata, invalid caches
- **Concept testing**: Consistency challenges understanding

## follow_up_weak - Foundation Check
**Question**: "Think about a restaurant with a host and servers. How is this similar to control/data plane separation?"
- **Simplification**: Real-world service analogy
- **Building block**: Role separation concept

## bar_raiser_question - L5→L6 Challenge
"Design a system that dynamically adjusts the boundary between control and data planes based on workload. When would you move responsibilities between planes?"

### bar_raiser_concepts
- Dynamic responsibility migration
- Workload-aware architecture adaptation
- Hybrid models with local control
- Edge computing implications
- Autonomous operation capabilities

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 4-5 min answer + 5-7 min discussion
- **Common next topics**: Service mesh, SDN, edge computing, serverless architectures
