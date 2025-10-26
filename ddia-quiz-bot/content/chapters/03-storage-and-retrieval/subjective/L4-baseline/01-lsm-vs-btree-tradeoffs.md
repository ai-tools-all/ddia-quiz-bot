---
id: storage-subjective-L4-001
type: subjective
level: L4
category: baseline
topic: storage-and-retrieval
subtopic: lsm-vs-btree
estimated_time: 6-8 minutes
---

# question_title - LSM Trees vs B‑Trees: Trade‑offs

## main_question - Core Question
"Compare LSM‑tree and B‑tree storage engines across write/read/space amplification, latency profiles, and operational behavior. When would you choose one over the other?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Write Amplification**: LSM compaction vs B‑tree page flush/splits
- **Read Amplification**: Multiple SSTables vs single index path
- **Space Amplification**: Stale versions in LSM; fill‑factor in B‑trees
- **Latency**: Tail behavior (compaction stalls vs page splits)

### expected_keywords
- Primary: amplification, compaction, page splits, locality
- Technical: memtable/WAL, leveled vs size‑tiered, cache

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Hardware Fit**: HDD vs SSD characteristics
- **Operational**: Backpressure, throttling, tuning knobs
- **Workload Fit**: Write‑heavy vs read‑heavy; range scans

### bonus_keywords
- Implementation: prefix compression, bloom filters, page size
- Related: Fractal trees, Bw‑trees, FASTER

## sample_excellent - Example Excellence
"LSMs turn random writes into sequential appends and background merges, yielding high write throughput and good SSD/HDD performance, at the cost of read amplification (multiple tables to check) and compaction‑induced CPU/IO. B‑trees perform in‑place updates and keep a single ordered index; point and range reads have predictable paths, but writes cause page splits and random IO, stressing HDDs and still incurring write amp. Choose LSM for write‑heavy workloads and large working sets with acceptable read amp (mitigated by Bloom filters), and B‑trees for read/range‑heavy OLTP with tight tail latencies and strong locality (especially with clustered indexes)."

## sample_acceptable - Minimum Acceptable
"LSMs are better for writes due to sequential IO but worse for reads; B‑trees are better for reads but have more random‑write costs."

## common_mistakes - Watch Out For
- Ignoring tail latency (compaction stalls) in LSMs
- Assuming B‑trees always beat LSMs for reads without cache considerations
- Overlooking space amplification differences

## follow_up_excellent - Depth Probe
**Question**: "Explain how leveled compaction changes the read/write/space amplification trade‑offs vs size‑tiered."
- **Looking for**: Leveled lowers read/space amp, increases write amp

## follow_up_partial - Guided Probe  
**Question**: "How do SSDs vs HDDs influence your choice?"
- **Hint**: Random IO penalty differs; sequential write advantages

## follow_up_weak - Foundation Check
**Question**: "What’s write amplification in one sentence?"
- **Simplification**: Extra bytes written per byte of user data

## bar_raiser_question - L4→L5 Challenge
"Design a switchable engine that toggles between LSM and B‑tree modes based on workload signals (write ratio, cache hit rate, tail latencies). What metrics and thresholds drive the change without thrashing?"

### bar_raiser_concepts
- Metrics (W/R ratio, compaction debt, split rate)
- Hysteresis to avoid flapping
- Data migration cost and dual‑write strategies

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 3-4 min answer + 3-4 min discussion
- **Common next topics**: Compaction strategy, cache/tiered storage, observability

## assistant_answer
LSMs favor write throughput (sequential appends + compaction) at the cost of higher read and space amp; B‑trees favor predictable read paths and range scans but pay with random‑write costs and page splits. Choose based on workload and latency SLOs.

## improvement_suggestions
- Require a small table comparing amplification metrics under leveled vs size‑tiered vs B‑tree.
- Ask for device‑aware tuning notes (SSD page size, queue depth).

## improvement_exercises
### exercise_1 - Amplification Table
**Question**: "Fill a qualitative table (Low/Med/High) for read/write/space amplification across size‑tiered LSM, leveled LSM, and B‑tree."

**Sample answer**: "Size‑tiered: write=Low, read=High, space=High; Leveled: write=Med/High, read=Low/Med, space=Low; B‑tree: write=Med, read=Low, space=Med."

### exercise_2 - Device Fit
**Question**: "Explain why HDDs favor LSM more than B‑trees, and how SSDs change the calculus."

**Sample answer**: "HDDs penalize random writes heavily; LSM sequential writes win. SSDs reduce random penalty, making B‑trees more attractive for read‑heavy workloads."
