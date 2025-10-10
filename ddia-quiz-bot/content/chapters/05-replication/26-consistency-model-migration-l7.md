---
id: ch05-consistency-model-migration-l7
day: 26
level: L7
tags: [replication, consistency, migration, architecture, principal-engineer]
related_stories: []
---

# Consistency Model Migration at Scale

## question
Your social media platform with 500M users currently uses single-leader replication with strong consistency for all operations. This causes scaling issues and poor user experience in distant regions. Design a migration strategy to selective consistency - strong consistency for critical paths (payments, privacy settings) and eventual consistency for others (likes, comments). The migration must be zero-downtime with ability to rollback. How do you approach this fundamental architecture change?

## expected_concepts
- Gradual migration with feature flags
- Dual-write patterns during transition
- Consistency SLAs and monitoring
- Client-side handling of eventual consistency
- Rollback strategies and circuit breakers
- A/B testing infrastructure changes
- Data reconciliation during migration

## answer
Phase 1 - Classification and Instrumentation: Categorize every API endpoint and data model by consistency requirements. Instrument current system to measure actual consistency lag when simulated under eventual consistency. Add client-side retry logic and conflict resolution before migration starts.

Phase 2 - Dual-Mode Architecture: Build parallel eventual consistency path alongside existing strong consistency. Implement consistency router that can dynamically choose path based on operation type and feature flags. Deploy read replicas but initially route 0% traffic to them. Add comprehensive monitoring for consistency SLA violations.

Phase 3 - Gradual Migration: Start with read-heavy, non-critical operations (viewing posts, browsing feeds). Use percentage-based feature flags per operation type, not per user, to ensure consistent experience. Implement "consistency budget" - automatically fall back to strong consistency if eventual consistency SLA breaches threshold.

Phase 4 - Critical Decision Points: For mixed operations (read profile, then update), implement session stickiness to ensure read-your-writes. For social features, accept temporary inconsistency but implement "healing" - background jobs that detect and fix consistency violations. Maintain strong consistency escape hatch for customer support operations.

Rollback Strategy: Each phase is independently reversible via feature flags. Maintain dual-write to both systems during migration. If issues arise, flip back to strong consistency within seconds. Keep rollback capability for 2x the expected migration duration.

Key insight: The hardest part isn't the technical migration but changing the organization's mental model. Engineers must learn to design for eventual consistency, customer support needs new tools, and product managers must accept temporary inconsistencies as normal.

## hook
How do you handle user complaints about "missing" data that's actually just eventual consistency lag?

## follow_up
Three months into migration, you discover that 15% of users are "super users" who generate 60% of the write traffic and have followers across all regions. Their actions cause cascade effects where eventual consistency leads to confusing experiences for millions of followers (seeing responses before original posts, likes appearing/disappearing). How do you handle this bimodal distribution without compromising the benefits of eventual consistency for the remaining 85% of users?

## follow_up_answer
Implement a dynamic consistency model based on user influence and action impact: (1) User Classification: Calculate "influence score" based on follower count, cross-region distribution, and interaction rates. Classify users into tiers with different consistency guarantees. (2) Adaptive Consistency: Super users' posts use "bounded staleness" - maximum 1-second lag globally via priority replication channels. Their interactions (likes, comments) on others' content use regional consistency with faster anti-entropy cycles.

Technical Implementation: Deploy "fast path" replication for high-influence accounts using dedicated Kafka topics with guaranteed ordering. Implement "influence-aware" caching where super users' content has shorter TTLs and aggressive invalidation. Use predictive pre-warming - when a super user starts composing, prepare infrastructure in all regions.

Cascade Management: Implement "causal consistency tokens" for related actions - responses carry parent post version, ensuring ordering. Deploy "influence amplification detection" - if an action might affect >1M users, temporarily upgrade to stronger consistency. Use "eventual notification" pattern - show content optimistically but delay notifications until consistency is achieved.

Economic Trade-off: This complexity is justified because super users drive engagement and revenue disproportionately. The infrastructure cost of special handling for 15% of users is offset by improved experience for their millions of followers. Long-term solution: investigate edge computing to bring computation closer to these users' followers, reducing the consistency/latency trade-off.
