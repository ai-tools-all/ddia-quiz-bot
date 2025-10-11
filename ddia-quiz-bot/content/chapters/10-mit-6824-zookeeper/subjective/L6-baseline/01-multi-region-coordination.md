---
id: zookeeper-subjective-L6-001
type: subjective
level: L6
category: baseline
topic: zookeeper
subtopic: multi-region
estimated_time: 12-15 minutes
---

# question_title - Multi-Region Coordination Architecture

## main_question - Core Question
"Design a global coordination system using Zookeeper that spans 5 regions with strict consistency requirements for some operations and eventual consistency for others. Address latency optimization, failure domains, and the trade-offs between consistency models. How do you handle region failures and network partitions?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Hierarchical Architecture**: Global and regional clusters with clear responsibilities
- **Consistency Zones**: Different consistency guarantees for different data
- **Latency Optimization**: Regional reads, global writes where needed
- **Partition Handling**: Degraded mode operations during splits

### expected_keywords
- Primary keywords: multi-region, consistency zones, partition tolerance, latency
- Technical terms: observer nodes, hierarchical clusters, consensus domains

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Conflict Resolution**: Strategies for concurrent updates across regions
- **Data Sovereignty**: Compliance with regional data requirements
- **Cost Optimization**: Balancing performance with infrastructure costs
- **Capacity Planning**: Regional growth patterns and scaling
- **Observability**: Cross-region tracing and debugging
- **DR Strategy**: Region failure recovery procedures
- **Client Routing**: Smart client libraries for optimal routing

### bonus_keywords
- Patterns: CRDT integration, vector clocks, causal consistency
- Architecture: edge clusters, CDC patterns, global load balancing
- Operations: chaos testing, game days, runbook automation

## sample_excellent - Example Excellence
"I'll design a hierarchical multi-region coordination system with differentiated consistency:

Architecture Overview:
```
Global Tier (Strict Consistency):
├── Global Primary Cluster (US-East)
│   ├── 5 voting nodes across 3 AZs
│   └── Handles: global configs, leader election, critical state
│
├── Global Secondary Cluster (EU-Central)
│   ├── 5 voting nodes (failover ready)
│   └── Async replication from primary
│
Regional Tier (Regional Consistency):
├── US-West Regional Cluster
│   ├── 3 voting + 5 observer nodes
│   └── Handles: regional configs, local coordination
│
├── EU-West Regional Cluster
├── APAC Regional Cluster
├── US-Central Regional Cluster
└── SA-East Regional Cluster

Edge Tier (Eventual Consistency):
└── Observer nodes in 20+ PoPs for read scaling
```

Consistency Zones:
```yaml
Strict Consistency (via Global Tier):
- Global leader election
- Cross-region service registry
- Critical configuration flags
- Distributed lock coordination
- Write latency: 50-200ms based on region

Regional Consistency (via Regional Tier):
- Regional service discovery
- Local rate limits
- Regional feature flags
- Session management
- Write latency: 5-20ms within region

Eventual Consistency (via Edge Tier):
- Read-heavy configs
- Static metadata
- Public keys/certificates
- Read latency: <5ms from edge
```

Client Routing Logic:
```python
class MultiRegionZKClient:
    def __init__(self, region, consistency_requirements):
        self.global_cluster = connect_global()
        self.regional_cluster = connect_regional(region)
        self.edge_cache = connect_edge()
    
    def read(path, consistency='eventual'):
        if consistency == 'strict':
            return self.global_cluster.sync_read(path)
        elif consistency == 'regional':
            return self.regional_cluster.read(path)
        else:
            return self.edge_cache.read(path) or self.regional_cluster.read(path)
    
    def write(path, data, scope='regional'):
        if scope == 'global' or is_critical_path(path):
            return self.global_cluster.write(path, data)
        else:
            result = self.regional_cluster.write(path, data)
            self.async_propagate(path, data)  # Best effort global sync
            return result
```

Partition Handling:
1. **Region Isolation**: Regional cluster continues serving regional data
2. **Degraded Global**: Cache last known global state, allow local decisions
3. **Healing**: Conflict resolution using vector clocks post-partition
4. **Circuit Breakers**: Prevent cascading failures across regions

Region Failure Scenarios:
```
Scenario: US-East (Global Primary) fails
1. EU-Central promotes to global primary (automatic, 30s)
2. Regional clusters detect via health checks
3. Clients failover to EU endpoints
4. When US-East recovers, becomes secondary
```

Latency Optimization:
- **Sticky Routing**: Clients prefer regional clusters
- **Predictive Caching**: Pre-warm edge nodes based on access patterns
- **Batching**: Aggregate cross-region updates
- **Compression**: Reduce bandwidth for large updates

Cost Optimization:
- Reserved instances for voting nodes (predictable)
- Spot/preemptible for observers (can lose)
- Auto-scaling observers based on load
- Cross-region traffic via private network (AWS PrivateLink, GCP Private Service Connect)

Monitoring & Observability:
```yaml
Per-Region Metrics:
- Latency percentiles (p50, p99, p999)
- Cross-region replication lag
- Consistency violations detected
- Partition events and duration

Global Dashboard:
- Region health matrix
- Global consistency score
- Cost per transaction by region
- Client distribution heat map
```

Trade-offs Made:
1. **Complexity vs Flexibility**: Hierarchical model adds complexity but enables consistency zones
2. **Cost vs Performance**: Edge nodes increase cost but critical for global latency
3. **Consistency vs Availability**: Degraded mode trades consistency for availability
4. **Automation vs Control**: Automated failover with manual override capability

This architecture provides <10ms reads globally for eventual consistency, <50ms for regional consistency, and maintains strict consistency for critical operations while handling region failures gracefully."

## sample_acceptable - Minimum Acceptable
"Deploy separate Zookeeper clusters per region for regional data and a global cluster for cross-region coordination. Use observer nodes in each region connected to the global cluster for read scaling. Implement different consistency levels by routing requests to appropriate clusters. Handle region failures by promoting regional clusters and network partitions by operating in degraded mode with cached data."

## common_mistakes - Watch Out For
- Single global cluster for everything (latency issues)
- Not differentiating consistency requirements
- No partition handling strategy
- Ignoring cost implications
- Missing conflict resolution

## follow_up_excellent - Depth Probe
**Question**: "How would you handle a scenario where regulatory requirements mandate that certain data never leaves specific regions?"
- **Looking for**: Data sovereignty patterns, encryption strategies, compliant replication
- **Red flags**: Trying to replicate everything globally

## follow_up_partial - Guided Probe  
**Question**: "What happens when two regions make conflicting updates during a network partition?"
- **Hint embedded**: Need conflict resolution strategy
- **Concept testing**: Understanding distributed consistency challenges

## follow_up_weak - Foundation Check
**Question**: "If you had offices worldwide, how would you share some information globally while keeping other information local?"
- **Simplification**: Basic data locality concepts
- **Building block**: Understanding scope and distribution

## bar_raiser_question - L6→L7 Challenge
"Design a system that provides a unified API over Zookeeper, etcd, and Consul deployments across different regions, where each region uses a different coordination service based on local requirements. Include automatic workload migration between backends during incidents."

### bar_raiser_concepts
- Abstraction layer design
- Protocol translation
- State migration strategies
- Consistency model mapping
- Automated failover orchestration
- Multi-vendor management

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 5-6 min answer + 7-9 min discussion
- **Common next topics**: Global infrastructure, edge computing, CDN architectures
