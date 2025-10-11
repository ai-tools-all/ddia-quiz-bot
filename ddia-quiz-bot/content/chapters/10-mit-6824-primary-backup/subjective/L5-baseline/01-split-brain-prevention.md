---
id: primary-backup-subjective-L5-001
type: subjective
level: L5
category: baseline
topic: primary-backup
subtopic: split-brain
estimated_time: 7-9 minutes
---

# question_title - Split-Brain Prevention in Primary-Backup Systems

## main_question - Core Question
"How can split-brain occur in primary-backup replication when network partitions happen? Explain the mechanisms VMware FT uses to prevent split-brain, and discuss the trade-offs involved in these solutions."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Split-Brain Definition**: Both replicas think they're primary
- **Network Partition**: Primary and backup can't communicate
- **Test-and-Set Service**: Atomic arbitration mechanism
- **Availability vs Consistency**: CAP theorem implications

### expected_keywords
- Primary keywords: split-brain, network partition, test-and-set, arbitration
- Technical terms: atomic operation, shared storage, majority voting

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Shared Storage Requirement**: VMFS for test-and-set operation
- **Lease Mechanism**: Time-bounded primary role
- **Quorum Systems**: Alternative approach with odd replicas
- **Client Routing**: How clients find current primary

### bonus_keywords
- Implementation: storage heartbeat, SCSI reservations, witness server
- Alternatives: Paxos, Raft, majority voting, odd replica count
- Edge cases: slow networks vs failures, Byzantine scenarios

## sample_excellent - Example Excellence
"Split-brain occurs when network partition separates primary and backup, but both remain operational. Each might independently decide to become primary, leading to divergent state and inconsistent client views. This violates the fundamental guarantee of single-copy consistency. VMware FT prevents this using an atomic test-and-set service on shared storage (VMFS). When replicas lose contact, both attempt to acquire a lock via test-and-set. Only one succeeds and continues as primary; the other halts or restarts as backup. This requires reliable shared storage accessible to both replicas - a potential single point of failure. The trade-off is choosing consistency over availability during partitions (per CAP theorem). Alternative approaches include using an odd number of replicas with majority voting (like Raft), or external arbitration services. The key insight is that preventing split-brain requires some form of agreement mechanism outside the two replicas themselves, whether shared storage, majority quorum, or external coordinator."

## sample_acceptable - Minimum Acceptable
"Split-brain happens when network partition makes primary and backup unable to communicate, but both stay alive. Both might become primary, serving different client requests and diverging. VMware FT prevents this using atomic test-and-set on shared storage - when they lose contact, both race to acquire a lock, winner becomes/stays primary, loser stops. This sacrifices availability during partition for consistency."

## common_mistakes - Watch Out For
- Thinking heartbeat loss always means failure
- Not recognizing need for external arbitration
- Assuming network partition means node failure
- Missing availability trade-off

## follow_up_excellent - Depth Probe
**Question**: "Compare VMware FT's shared storage approach to Raft's majority voting. What are the minimum resource requirements and failure scenarios each can handle?"
- **Looking for**: 2 nodes + storage vs 3+ nodes, different failure modes
- **Red flags**: Not understanding fundamental arbitration requirements

## follow_up_partial - Guided Probe  
**Question**: "What happens if the shared storage itself becomes partitioned from both replicas?"
- **Hint embedded**: Neither can acquire lock
- **Concept testing**: Understanding dependency on arbitration mechanism

## follow_up_weak - Foundation Check
**Question**: "Two generals need to attack together but their messenger might be captured. How is this similar to the split-brain problem?"
- **Simplification**: Classic distributed systems problem
- **Building block**: Impossibility of consensus without reliable communication

## bar_raiser_question - L5→L6 Challenge
"Design a primary-backup system for a globally distributed database with replicas in US, Europe, and Asia. Network partitions between continents are common. How do you prevent split-brain while maximizing availability? Consider client locality and latency requirements."

### bar_raiser_concepts
- Geographic distribution challenges
- Multiple partition scenarios
- Regional failover strategies
- Consistency models per region
- Client affinity and routing
- Witness placement strategies

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 4-5 min answer + 3-4 min discussion
- **Common next topics**: Consensus protocols, geo-replication, eventual consistency
