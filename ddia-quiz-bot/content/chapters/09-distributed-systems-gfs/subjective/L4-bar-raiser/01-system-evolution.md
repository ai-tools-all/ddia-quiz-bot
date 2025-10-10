---
id: gfs-subjective-L4-bar
type: subjective
level: L4
category: bar-raiser
topic: gfs
subtopic: system-evolution
estimated_time: 10-12 minutes
---

# question_title - GFS to Modern Requirements

## main_question - Core Question
"You're tasked with building 'GFS 2.0' for 2025 workloads: mix of AI training (large files), microservices (small files), and real-time analytics (strong consistency). What would you change from the original GFS design and why?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Workload Diversity**: Different consistency/performance needs
- **Multi-Master**: Scale beyond single master limits
- **Consistency Tiers**: Strong for analytics, relaxed for training
- **Small File Handling**: Address original GFS weakness

### expected_keywords
- Primary keywords: multi-master, consistency levels, workload-aware
- Technical terms: sharding, consensus, tiered storage

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **ML Workload Optimization**: Checkpoint/restore, GPU direct storage
- **Container Integration**: CSI, Kubernetes native
- **Disaggregated Storage**: Separate compute and storage scaling
- **Multi-Region**: Geo-replication with locality
- **Cost Optimization**: Hot/warm/cold tiers, compression
- **Security**: Encryption, tenant isolation

### bonus_keywords
- Implementation: RDMA, NVMe, persistent memory
- Standards: S3 compatibility, POSIX where needed
- Modern: Object storage, erasure coding, deduplication

## sample_excellent - Example Excellence
"For GFS 2.0, I'd architect a multi-tier system: 1) Replace single master with sharded masters using consistent hashing for namespace distribution, with Raft consensus per shard. 2) Introduce storage classes: 'Strong' tier using chain replication for analytics needing consistency, 'Eventually consistent' tier for AI training data optimizing throughput, 'Small file' tier with smaller blocks and aggregation for microservices. 3) Add RDMA support for high-throughput AI workloads bypassing kernel. 4) Implement erasure coding (not just 3x replication) for cold data cost optimization. 5) Provide S3-compatible API alongside traditional file interface. 6) Add multi-tenancy with QoS guarantees. The key insight is that no single storage system can optimize for all workloads, so GFS 2.0 should be a platform providing different guarantees for different needs while sharing infrastructure."

## sample_acceptable - Minimum Acceptable
"I would replace the single master with multiple masters for scalability, add support for different consistency levels depending on workload needs, and improve small file handling with variable chunk sizes or a separate small file storage system."

## common_mistakes - Watch Out For
- Trying to solve everything with one approach
- Ignoring backward compatibility
- Over-engineering without justification
- Missing modern workload characteristics

## follow_up_excellent - Depth Probe
**Question**: "Your multi-tier design increases operational complexity significantly. How do you prevent this from becoming unmaintainable?"
- **Looking for**: Automation, standardization, clear tier boundaries, migration tools
- **Red flags**: Hand-waving complexity away

## follow_up_partial - Guided Probe
**Question**: "You mentioned different consistency levels. How would an application choose which level to use, and could it change dynamically?"
- **Hint embedded**: API design, performance implications
- **Concept testing**: Practical implementation thinking

## follow_up_weak - Foundation Check
**Question**: "Let's focus on just one aspect - AI training workloads. What makes their storage needs different from traditional web applications?"
- **Simplification**: Single workload type
- **Building block**: Workload characteristics

## bar_raiser_question - L5→L6 Challenge
"Now consider the business side: GFS 2.0 needs to be commercially viable as a cloud service competing with S3/Azure Storage. How does this change your technical design?"

### bar_raiser_concepts
- Multi-tenancy and isolation requirements
- SLA guarantees and monitoring
- Billing and metering infrastructure
- Compliance and data sovereignty
- Ecosystem and tooling compatibility
- Migration paths from competitors

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 4-5 min answer + 5-7 min discussion
- **Common next topics**: Cloud storage economics, SDS trends, AI infrastructure
