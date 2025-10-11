---
id: craq-subjective-L5-006
type: subjective
level: L5
category: baseline
topic: craq
subtopic: hotspot-mitigation
estimated_time: 8-10 minutes
---

# question_title - Mitigating Hotspots in CRAQ

## main_question - Core Question
"A single key in CRAQ experiences 10× the normal write volume due to a celebrity event. Devise a mitigation plan that preserves linearizability and leverages DDIA's guidance on partitioning and caching."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Hotspot Isolation**: Move key to dedicated chain or shard to prevent queue buildup
- **Write Coalescing**: Batch updates or use command pattern to reduce tail load
- **Cache Strategy**: Use read-through caches with invalidation tied to tail commit
- **DDIA Link**: Recognize impact of uneven partitioning and apply mitigation patterns

### expected_keywords
- Primary keywords: hotspot, partitioning, batching, cache invalidation
- Technical terms: key relocation, hot shard, write amplification, event throttling

## peripheral_concepts - Nice to Have (40%)
- **Backpressure**: Rate-limit writes to protect propagation budget
- **Feature Flags**: Temporarily degrade less critical features for the hot key
- **Observability**: Monitor dirty duration spikes as detection signal
- **Customer Comms**: Align with operational best practices in DDIA

### bonus_keywords
- Implementation: consistent hashing adjustment, dedicated chain, command queue, TTL cache
- Scenarios: flash sale, trending topic, denial of service attempt
- Trade-offs: added complexity, possible stale cache windows, cost of dedicated hardware

## sample_excellent - Example Excellence
"We'd isolate the celebrity key into its own chain to prevent other data from inheriting the dirty backlog, a classic DDIA partitioning tactic. Writes would be batched or converted into incremental commands (e.g., 'increment by N') to reduce tail pressure. Reads would hit a cache invalidated only after tail commit, maintaining linearizability. If volume still exceeds capacity, we throttle low-priority write attempts and surface status updates, aligning with DDIA's advice on hotspot mitigation and graceful degradation." 

## sample_acceptable - Minimum Acceptable
"Move the hot key into its own CRAQ chain, batch its writes, and use caches that invalidate on tail commit. This keeps linearizability intact while applying DDIA's hotspot mitigation advice." 

## common_mistakes - Watch Out For
- Sharding key without respecting ordering metadata
- Serving stale cache entries before tail commit
- Ignoring user-facing degradation strategies
- Not tying approach to DDIA recommendations

## follow_up_excellent - Depth Probe
**Question**: "How would you detect the hotspot quickly and trigger the mitigation automatically?"
- **Looking for**: Dirty duration anomalies, tail queue length, adaptive sharding triggers
- **Red flags**: Manual detection only

## follow_up_partial - Guided Probe  
**Question**: "What metrics confirm your mitigation succeeded?"
- **Hint embedded**: Reduced dirty duration, stabilized latency, lower duplicate suppression
- **Concept testing**: Feedback measurement

## follow_up_weak - Foundation Check
**Question**: "When one checkout line explodes with customers, why open an express lane?"
- **Simplification**: Hotspot isolation analogy
- **Building block**: Load redistribution

## bar_raiser_question - L5→L6 Challenge
"Design a predictive model that shifts CRAQ chains ahead of expected hotspot events (e.g., live-stream drops), integrating DDIA's batch analytics guidance." 

### bar_raiser_concepts
- Predictive analytics, batch + streaming synergy
- Pre-splitting shards, warming caches
- Feedback loops with operational data
- Proactive vs reactive mitigation

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 4-5 min answer + 4-5 min discussion
- **Common next topics**: Predictive autoscaling, event-driven architecture, chaos rehearsal
