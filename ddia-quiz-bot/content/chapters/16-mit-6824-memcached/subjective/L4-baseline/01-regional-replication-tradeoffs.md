---
id: memcached-subjective-L4-001
type: subjective
level: L4
category: baseline
topic: memcached
subtopic: regional-replication
estimated_time: 8-10 minutes
---

# question_title - Regional Replication Trade-offs

## main_question - Core Question
"Explain Facebook's regional replication architecture. How does it balance latency and consistency, and why is this design acceptable for a social media workload?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Multiple regions (e.g., West Coast primary, East Coast secondary)
- All writes flow to primary region's MySQL master
- Asynchronous replication to secondary regions with seconds of lag
- Reads are local within each region (memcached + MySQL)
- Accepts stale reads for users far from primary

### expected_keywords
- regional replication, primary, secondary, asynchronous, lag, local reads, eventual consistency

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Read-heavy workload makes this viable
- Serving a single page requires hundreds of data items
- Cross-country latency would be prohibitive for synchronous reads
- Social media content rarely requires immediate global consistency

### bonus_keywords
- read-heavy, page composition, cross-datacenter latency, user-perceived staleness

## sample_excellent - Example Excellence
"Facebook deploys multiple regions, each with complete data replicas. All writes flow to the primary region's MySQL master for a single source of truth. Asynchronous MySQL replication propagates updates to secondary regions with seconds of lag. Critically, reads remain local—both memcached and MySQL queries stay within the user's region. This design prioritizes low-latency local reads over global consistency. For social media, this is acceptable because users rarely notice slightly outdated news feeds or profile data, and cross-region consistency isn't critical. Serving a single page often requires hundreds of data items, so local access is essential for performance. The read-heavy workload (99%+ reads) makes this trade-off worthwhile."

## sample_acceptable - Minimum Acceptable
"Multiple regions with primary for writes and secondary for reads. Async replication causes lag. Reads are local for low latency. Social media can tolerate slightly stale data because content isn't time-critical."

## common_mistakes - Watch Out For
- Assuming synchronous cross-region replication
- Not explaining why staleness is acceptable for social media
- Missing that both memcached AND MySQL reads are local

## follow_up_excellent - Depth Probe
**Question**: "What happens if the replication lag grows to minutes instead of seconds? What user-visible issues might arise?"
- **Looking for**: Users might see very outdated content, friend requests might appear inconsistent, notification counts might be wrong

## follow_up_partial - Guided Probe
**Question**: "Why serve reads locally instead of always reading from the primary region?"
- **Hint embedded**: Latency costs and page composition requirements

## follow_up_weak - Foundation Check
**Question**: "Which region handles writes and why?"
