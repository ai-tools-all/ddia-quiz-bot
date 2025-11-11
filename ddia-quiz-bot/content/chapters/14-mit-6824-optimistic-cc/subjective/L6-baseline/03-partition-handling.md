---
id: farm-subjective-L6-003
type: subjective
level: L6
category: baseline
topic: fault-tolerance
subtopic: network-partition
estimated_time: 12-15 minutes
---

# question_title - Network Partition and Split-Brain Prevention

## main_question - Core Question
"Analyze how FaRM handles network partitions within a single data center. Consider a partition that splits the cluster such that primary and backup for a shard end up in different partitions, both reachable by different sets of clients. How does FaRM prevent split-brain? What role does the configuration manager (ZooKeeper) play? Design a test scenario that would violate consistency if split-brain prevention were broken, and explain how FaRM's mechanisms prevent it."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Configuration manager (ZooKeeper) maintains authoritative shard-to-server mapping
- Leases or epochs ensure only one primary per shard at a time
- Primary must have valid lease from configuration manager to serve writes
- Partition healing: failed components detected, new primary elected, clients updated
- Split-brain scenario: both sides think they're primary, accepting conflicting writes

### expected_keywords
- split-brain, network partition, configuration manager, ZooKeeper, lease, epoch, primary election, consistency

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Single datacenter makes partitions less likely than WAN scenarios
- RDMA networks have different failure modes than TCP/IP
- Fencing mechanisms prevent zombie primaries
- Version numbers provide additional safety (monotonic increase per primary)
- Trade-off: availability during partition vs consistency

### bonus_keywords
- fencing, zombie primary, failover, lease timeout, quorum, consistency vs availability, CAP theorem

## sample_excellent - Example Excellence
"FaRM uses ZooKeeper as configuration manager to prevent split-brain. Each primary holds a lease with epoch number; writes require valid lease. Test scenario: partition splits P1 (primary for shard X) from B1 (backup). Both partitions have clients. Without prevention: clients in P1's partition write X@v5→v6; clients in B1's partition promote B1 to primary, write X@v5→v7; when partition heals, conflicting states. FaRM's prevention: ZooKeeper maintains quorum in majority partition; P1 keeps lease if in majority, else lease expires; B1 can only become primary after getting lease from ZooKeeper; minority partition can't form quorum, blocks writes. Clients observe unavailability rather than inconsistency. Version monotonicity + epochs provide defense in depth."

## sample_acceptable - Minimum Acceptable
"FaRM uses ZooKeeper to track which server is primary. Primaries need a lease to serve writes. If there's a partition, only one side can get the lease, preventing both sides from acting as primary."

## common_mistakes - Watch Out For
- Not explaining the role of leases/epochs
- Missing the quorum requirement for configuration changes
- Thinking FaRM maintains availability during partitions (it sacrifices availability for consistency)
- Not providing a concrete split-brain scenario
- Forgetting about client-side effects

## follow_up_excellent - Depth Probe
**Question**: "Design a recovery protocol for when a partitioned primary heals and discovers a new primary was elected in its place. How do you ensure consistency?"
- **Looking for**: Epoch comparison, state reconciliation, rollback uncommitted transactions, replay logs from new primary, version vector comparison

## follow_up_partial - Guided Probe
**Question**: "Why is a single datacenter deployment less vulnerable to partitions than geo-distributed systems?"
- **Hint embedded**: Shared fate, better network reliability, same physical infrastructure

## follow_up_weak - Foundation Check
**Question**: "What is split-brain and why is it dangerous in distributed systems?"
