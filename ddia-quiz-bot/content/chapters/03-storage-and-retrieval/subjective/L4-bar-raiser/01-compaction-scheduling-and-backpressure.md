---
id: storage-subjective-L4-003
type: subjective
level: L4
category: bar-raiser
topic: storage-and-retrieval
subtopic: compaction-scheduling
estimated_time: 8-10 minutes
---

# question_title - Compaction Scheduling and Backpressure

## main_question - Core Question
"Design a compaction scheduler for an LSM engine that sustains bursty ingest without violating tail read latency SLOs. What signals do you monitor, how do you throttle, and how do you prioritize compactions across levels?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Signals**: L0 file count, compaction debt, read latency, write stall indicators
- **Priorities**: L0→L1 first to control amplification and stalls
- **Throttling**: I/O bandwidth caps, CPU budgeting, admission control
- **Fairness**: Avoid starving lower levels and hot ranges

### expected_keywords
- Primary: backpressure, scheduling, priorities, SLOs
- Technical: per‑level quotas, token bucket, rate limiter

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Hot/Cold Separation**: Differentiate hot ranges for faster compaction
- **Adaptive Parallelism**: Scale threads by load
- **Integration**: Cache warming, bloom filter build cost

### bonus_keywords
- Implementation: picker heuristics, dynamic rebalancing
- Related: RocksDB triggers, L0 slowdown/stop triggers

## sample_excellent - Example Excellence
"Monitor L0 file count, total compaction debt, p99 read/commit latencies, and disk bandwidth. Prioritize L0→L1 compactions to cap amplification and reduce stalls; use per‑level quotas (token bucket) and global I/O caps. Apply admission control on ingest when L0 exceeds thresholds (slow/stop). Dynamically adjust compaction parallelism based on observed debt and SLOs; opportunistically compact hot ranges. Integrate with cache/Bloom build to prevent read‑path thrash."

## sample_acceptable - Minimum Acceptable
"Watch L0 growth and read p99; compact L0 first and throttle writes when thresholds are crossed."

## common_mistakes - Watch Out For
- Compaction everywhere at once causing cache churn
- Ignoring hot ranges that dominate tail latency
- No admission control leading to unbounded debt

## follow_up_excellent - Depth Probe
**Question**: "Propose concrete thresholds for L0 slowdown/stop and explain how they relate to device throughput."
- **Looking for**: Trigger values tied to MB/s and compaction catch‑up time

## follow_up_partial - Guided Probe  
**Question**: "How do you ensure compactions don’t evict hot pages from cache?"
- **Hint**: Throttle compactions, pin hot pages, warm cache post‑compaction

## follow_up_weak - Foundation Check
**Question**: "Why compact L0 first?"
- **Simplification**: Overlapping files cause worst read amp and stalls

## bar_raiser_question - L4→L5 Challenge
"Given SSD bandwidth 2GB/s and CPU budget 8 cores, target ingest 500MB/s with p99 read SLO 5ms. Design quotas and triggers to keep L0 stable and reads within SLO under a 2× burst."

### bar_raiser_concepts
- Budgeting, headroom for bursts
- Prioritization and hysteresis
- Observability and feedback loops

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 4-5 min answer + 4-5 min discussion
- **Common next topics**: I/O schedulers, admission control, QoS

## assistant_answer
Prioritize L0→L1 with per‑level quotas and a global rate limiter; throttle ingest via slowdown/stop triggers tied to compaction debt and L0 file counts; adapt parallelism and preserve cache locality to maintain read p99 while absorbing bursts.

## improvement_suggestions
- Ask for a numeric plan (MB/s per level, thread counts, thresholds) and hysteresis to avoid flapping.
- Include cache‑aware compaction scheduling.

## improvement_exercises
### exercise_1 - Trigger Tuning
**Question**: "Propose slowdown/stop thresholds for L0 and expected recovery times at 500MB/s ingest."

**Sample answer**: "Slowdown at 8–12 files, stop at 20–24 given compaction capacity; recovery modeled by compaction bandwidth headroom vs ingest."

### exercise_2 - Cache Preservation
**Question**: "Outline techniques to limit cache churn during heavy compaction."

**Sample answer**: "Pin hot pages, stagger compactions, warm filters/index blocks, throttle rebuild I/O."
