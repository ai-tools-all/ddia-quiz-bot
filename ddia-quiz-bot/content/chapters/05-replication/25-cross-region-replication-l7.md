---
id: ch05-cross-region-replication-l7
day: 25
level: L7
tags: [replication, architecture, geo-distribution, principal-engineer, trade-offs]
related_stories: []
---

# Cross-Region Replication Strategy

## question
Your company is expanding from US-only to global operations with strict data residency requirements in EU (GDPR) and China (local data laws). Current architecture uses single-leader PostgreSQL replication. As Principal Engineer, design a replication strategy that balances compliance, performance, consistency, and operational complexity. Consider that some data must stay within regions while other data needs global visibility.

## expected_concepts
- Data sovereignty and regulatory compliance
- Hybrid replication topology (mix of single/multi-leader)
- Logical data sharding by region/customer
- Selective replication and data classification
- Conflict resolution for shared global data
- Monitoring and operational complexity
- Disaster recovery across regions

## answer
The solution requires a hybrid approach: partition data into three categories - region-locked (PII, regulated data), globally replicated (product catalog, configurations), and region-preferred with global fallback (user sessions, caches).

Architecture: (1) Deploy regional clusters with local single-leader replication for region-locked data. Each region is authoritative for its users' regulated data. (2) Implement multi-leader replication for global data with designated regional conflict resolution priorities based on data ownership. (3) Use logical replication with row-level filters to control what crosses boundaries.

Key design decisions: Use application-level sharding to route users to their home region. Implement field-level encryption for any PII that must cross regions for disaster recovery. Deploy a global service mesh for intelligent request routing. Use CRDTs for naturally convergent global data like feature flags.

Operational considerations: Centralized schema management with region-specific migrations. Global observability with regional data isolation. Automated compliance auditing. Clear playbooks for region-specific failures vs global outages.

Second-order insight: The real challenge isn't technical but organizational - you need clear data governance, with each team understanding data classification and compliance implications of their design choices.

## hook
How do you maintain ACID properties when a transaction spans region-locked and global data?

## follow_up
Six months later, a critical bug in the application causes data corruption that propagates through multi-leader replication to all regions before detection. The corruption affects both financial records (requiring point-in-time consistency) and user-generated content (where latest version matters). How do you orchestrate recovery while minimizing downtime and data loss? Consider that regions have different backup schedules and some regions have already processed additional transactions on top of corrupted data.

## follow_up_answer
Recovery requires a coordinated but region-aware approach: (1) Immediately fence off writes globally to prevent further corruption. (2) Identify the corruption point-in-time using distributed tracing and audit logs - this timestamp becomes your global coordination point.

Recovery strategy by data type: For financial records, restore all regions to the last consistent snapshot before corruption, then replay validated transactions from the audit log. For user content, use a "merge recovery" - restore corrupted records from backup but preserve newer non-corrupted changes using timestamp comparison.

Execution sequence: (1) Create recovery environments parallel to production in each region. (2) Restore and validate data in recovery environments. (3) Build compensation transactions for financial inconsistencies. (4) Use blue-green switchover per region, starting with regions having least post-corruption activity. (5) Implement temporary eventual consistency mode during recovery with clear user communication.

Critical insight: Design your replication topology with recovery in mind from day one. This means maintaining separate audit logs from replication streams, region-specific backup strategies aligned with data criticality, and automated corruption detection that can freeze replication before propagation. The cost of these systems is justified by scenarios like this where hours matter.
