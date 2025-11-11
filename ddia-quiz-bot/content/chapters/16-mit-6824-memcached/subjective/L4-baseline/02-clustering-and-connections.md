---
id: memcached-subjective-L4-002
type: subjective
level: L4
category: baseline
topic: memcached
subtopic: intra-region-clustering
estimated_time: 8-10 minutes
---

# question_title - Clustering and Connection Management

## main_question - Core Question
"Describe Facebook's intra-region clustering strategy for memcached. What problems does clustering solve, and how does it balance hot-key handling with connection overhead?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Front-ends and memcached servers grouped into independent clusters
- Limits N-squared TCP connection overhead
- Reduces incast congestion from parallel requests
- Popular keys replicated across clusters for hot-key handling
- Regional pool for infrequently accessed items

### expected_keywords
- clustering, N-squared connections, incast, hot keys, replication, regional pool

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Without clustering: every front-end connects to every memcached server
- Page requests often need hundreds of items fetched in parallel
- Simultaneous responses can overwhelm network buffers
- Cold data shouldn't waste RAM across all clusters

### bonus_keywords
- connection explosion, parallel fetch, network buffers, memory efficiency

## sample_excellent - Example Excellence
"Within each region, front-ends and memcached servers organize into independent clusters to manage two problems: N-squared connection overhead and incast congestion. Without clustering, every front-end would need TCP connections to every memcached server—millions of connections. Clusters partition this, limiting connection counts. Additionally, when a front-end requests hundreds of items in parallel (common for rendering a page), responses arriving simultaneously cause incast congestion, overwhelming network buffers. Clusters reduce the blast radius. For hot keys that would bottleneck in a single partition, Facebook replicates them across clusters, multiplying serving capacity. Conversely, cold data goes to a regional pool shared across clusters, avoiding wasteful replication and dedicating per-cluster RAM to hot items."

## sample_acceptable - Minimum Acceptable
"Clusters limit connections between front-ends and memcached servers. They also reduce incast when many responses arrive at once. Hot keys are replicated across clusters, while cold data uses a shared regional pool."

## common_mistakes - Watch Out For
- Not explaining what N-squared connection overhead means
- Missing the incast congestion problem
- Confusing within-cluster partitioning with across-cluster replication

## follow_up_excellent - Depth Probe
**Question**: "Why is incast congestion worse when serving page requests compared to single-item lookups?"
- **Looking for**: Page rendering requires hundreds of parallel fetches, creating synchronized response bursts

## follow_up_partial - Guided Probe
**Question**: "What's the downside of replicating all data across all clusters?"
- **Hint embedded**: Memory waste for infrequently accessed data

## follow_up_weak - Foundation Check
**Question**: "What is N-squared connection overhead?"
