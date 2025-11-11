---
id: memcached-subjective-L5-001
type: subjective
level: L5
category: baseline
topic: memcached
subtopic: consistency-model
estimated_time: 10-12 minutes
---

# question_title - Consistency Model Analysis

## main_question - Core Question
"Analyze Facebook's consistency model for memcached. What level of consistency does the system provide, what mechanisms ensure read-your-writes, and why is this sufficient for social media despite accepting eventual consistency?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- System accepts eventual consistency with seconds of staleness
- Critical requirement: read-your-writes consistency
- Front-end deletes after writes guarantee users see their own updates
- Casual browsing (feeds, profiles) tolerates stale data
- Users rarely notice slightly outdated content
- Invalidation protocols and timeouts prevent indefinite staleness

### expected_keywords
- eventual consistency, read-your-writes, staleness, invalidation, user experience, tolerance

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Trade-off between performance and consistency
- Linearizability not required for most operations
- Dual invalidation (front-end + MySQL async) provides redundancy
- Replication across regions amplifies staleness but acceptable
- Timeout mechanisms as safety net

### bonus_keywords
- linearizability, CAP theorem, cross-region lag, redundant invalidation, cache expiration

## sample_excellent - Example Excellence
"Facebook's memcached system provides eventual consistency with seconds of staleness, prioritizing low-latency reads over strong consistency guarantees. The critical requirement is read-your-writes: after a user updates their data, they must immediately see the new value. This is guaranteed by front-end deletes issued immediately after database writes, forcing a cache miss that fetches fresh data. For other users viewing that data, eventual consistency suffices because social media content (news feeds, profiles, photos) rarely requires real-time accuracy—users don't notice if a post is a few seconds old. The system employs dual invalidation (front-end deletes plus MySQL async deletes) to ensure caches eventually refresh even if some deletes fail. Timeouts provide an additional safety net against indefinite staleness. This design trades global consistency for performance, which is appropriate given the read-heavy workload where 99%+ of requests hit cache."

## sample_acceptable - Minimum Acceptable
"Eventual consistency with read-your-writes. Front-end deletes ensure users see their updates. Other users can see stale data briefly. This works for social media because exact freshness isn't critical for browsing."

## common_mistakes - Watch Out For
- Not distinguishing read-your-writes from eventual consistency
- Missing the dual invalidation mechanism
- Not explaining why staleness is acceptable for social media specifically
- Confusing linearizability with read-your-writes

## follow_up_excellent - Depth Probe
**Question**: "Under what social media scenarios would eventual consistency become problematic? Design a consistency model for those edge cases."
- **Looking for**: Financial transactions (ads, payments), security-sensitive operations (password changes), real-time messaging—these need stronger guarantees

## follow_up_partial - Guided Probe
**Question**: "How does regional replication affect the consistency model?"
- **Hint embedded**: Amplifies staleness across regions, but within-region consistency still maintained

## follow_up_weak - Foundation Check
**Question**: "What is read-your-writes consistency and why does it matter?"
