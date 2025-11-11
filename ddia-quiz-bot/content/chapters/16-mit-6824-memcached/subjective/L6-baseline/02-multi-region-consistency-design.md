---
id: memcached-subjective-L6-002
type: subjective
level: L6
category: baseline
topic: memcached
subtopic: multi-region-design
estimated_time: 12-15 minutes
---

# question_title - Multi-Region Consistency Design

## main_question - Core Question
"Design a multi-region caching system for a global e-commerce platform where users expect to see consistent inventory and pricing across devices/sessions, but can tolerate some staleness for product reviews. How would you modify Facebook's regional replication approach? Consider: write routing, inventory deduction correctness, cross-region cache invalidation, and failure handling when the primary region is unavailable."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Different consistency requirements for different data types
- Inventory/pricing needs stronger consistency than reviews
- Write routing strategies: all-to-primary vs per-item home region
- Cross-region invalidation must reach all caches
- Failover when primary unavailable: promote secondary, accept temporary inconsistency, or queue writes

### expected_keywords
- tiered consistency, strongly consistent, eventual consistency, write routing, failover, inventory correctness

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Distributed transactions for inventory deduction
- Compensating transactions if inventory oversold
- Read-from-primary for critical reads (payment)
- Regional pools with different TTLs per data type
- Monitoring replication lag per data category
- User session affinity to reduce apparent inconsistency

### bonus_keywords
- distributed transaction, compensation, session affinity, lag monitoring, TTL tuning, quorum

## sample_excellent - Example Excellence
"E-commerce requires tiered consistency. For inventory and pricing, I'd use stronger guarantees: writes go to a designated primary region for each product SKU (partitioned by product ID for write scalability), with synchronous replication or quorum writes to ensure cross-region consistency within SLA. Reads for checkout/payment always query the primary region or require quorum reads. For reviews, use Facebook's eventual consistency model—local reads with async replication. Cross-region cache invalidation is critical: after inventory update, send deletes to all regional memcached clusters, with retries and monitoring to ensure propagation. For primary region failure, options include: (1) promote secondary with brief inconsistency window, (2) queue writes and serve reads-only from secondaries, (3) implement consensus protocol (Raft/Paxos) for coordinated failover. Additionally, use session affinity (sticky routing) so a user's requests hit the same region, reducing apparent inconsistency. Monitor replication lag per data category and alert when inventory lag exceeds threshold."

## sample_acceptable - Minimum Acceptable
"Inventory needs strong consistency, reviews can be eventual. Route writes to primary for inventory. Invalidate all regional caches on inventory change. Failover to secondary region if primary fails. Use different consistency for different data."

## common_mistakes - Watch Out For
- Applying same consistency model to all data types
- Not addressing inventory correctness specifically
- Missing cross-region invalidation complexity
- Ignoring failover scenario
- Not considering write scalability with all-to-primary routing

## follow_up_excellent - Depth Probe
**Question**: "Your inventory system temporarily oversold due to replication lag. Design a compensating transaction system to handle this gracefully."
- **Looking for**: Reservation system, hold inventory temporarily, rollback/cancel orders, customer communication, eventual convergence

## follow_up_partial - Guided Probe
**Question**: "Why is session affinity useful for reducing apparent inconsistency?"
- **Hint embedded**: User sees consistent view within their session even if other regions lag

## follow_up_weak - Foundation Check
**Question**: "What's the difference between eventual consistency and strong consistency?"
