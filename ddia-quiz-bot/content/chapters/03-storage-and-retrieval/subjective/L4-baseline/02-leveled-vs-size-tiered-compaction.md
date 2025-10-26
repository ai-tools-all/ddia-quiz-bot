---
id: storage-subjective-L4-002
type: subjective
level: L4
category: baseline
topic: storage-and-retrieval
subtopic: leveled-vs-size-tiered
estimated_time: 6-8 minutes
---

# question_title - Leveled vs Size‑Tiered Compaction

## main_question - Core Question
"Contrast leveled and size‑tiered compaction strategies in LSM engines. Explain their impact on read/write/space amplification, write stalls, and tuning knobs." 

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Size‑Tiered**: Merge similarly sized runs; fewer merges, more overlap
- **Leveled**: Strict level sizes, non‑overlapping ranges per level
- **Amplifications**: Leveled lowers read/space amp; raises write amp
- **Operational**: Backpressure, compaction debt, throttling

### expected_keywords
- Primary: overlap, levels, fanout, compaction debt
- Technical: file count, SSTable size, I/O budget

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Hot Keys**: Rewrites in leveled increase write amp
- **Cold Data**: Leveled stores once at Lmax
- **Tuning**: Target file size, level fanout, parallelism

### bonus_keywords
- Implementation: compaction picker, tombstone purge horizon
- Related: tiered‑leveled hybrids

## sample_excellent - Example Excellence
"Size‑tiered batches similarly sized SSTables and merges them occasionally, resulting in low write amplification but many overlapping files to probe on reads (higher read and space amplification). Leveled organizes data into levels with strict size ratios and non‑overlapping ranges per level, which reduces read and space amplification—at the cost of higher write amplification because keys rewrite multiple times while cascading through levels. Operationally, leveled needs careful throttling to avoid stalls under write bursts (compaction debt), while size‑tiered often delivers smoother writes but worse tail read latency."

## sample_acceptable - Minimum Acceptable
"Leveled compaction improves reads and space at the cost of more writes; size‑tiered does the opposite."

## common_mistakes - Watch Out For
- Not mentioning overlapping ranges in size‑tiered
- Ignoring the space amplification differences
- Forgetting about backpressure or I/O budgeting

## follow_up_excellent - Depth Probe
**Question**: "Explain how increasing target file size and level fanout affects write amplification and read amplification in leveled compaction."
- **Looking for**: Larger files/fanout → fewer levels → lower write amp but larger read working sets

## follow_up_partial - Guided Probe  
**Question**: "When would you pick a hybrid strategy?"
- **Hint**: Balance bursty ingest with acceptable read tails

## follow_up_weak - Foundation Check
**Question**: "What does ‘overlap’ mean in this context?"
- **Simplification**: Same key range appearing in multiple files you must check

## bar_raiser_question - L4→L5 Challenge
"Given ingest 100MB/s sustained and a compaction budget of 400MB/s, propose leveled settings (file size, fanout, parallelism) to avoid compaction debt. Analyze read amp at steady state."

### bar_raiser_concepts
- Balance ingest vs compaction throughput
- Calculate levels and rewrite rate
- Tail read cost model

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 3-4 min answer + 3-4 min discussion
- **Common next topics**: Bloom filter sizing, cache warmup, throttling

## assistant_answer
Leveled uses non‑overlapping ranges per level to cut read/space amp at the expense of more rewrites; size‑tiered reduces write amp but increases file overlap and read amp. Tuning fanout and file size balances compaction rewrite rate vs read working set.

## improvement_suggestions
- Ask for a numeric example estimating write amplification under different fanouts.
- Include an I/O budget exercise to reason about compaction debt and throttling.

## improvement_exercises
### exercise_1 - Write Amp Estimate
**Question**: "Estimate write amplification for leveled compaction with fanout 10 across 5 levels."

**Sample answer**: "Ballpark ~sum over levels of rewrite factor; often 10–20× total depending on overlap and tombstones; show reasoning for chosen model."

### exercise_2 - I/O Budgeting
**Question**: "With ingest 100MB/s and compaction capacity 400MB/s, what maximum fanout keeps debt bounded?"

**Sample answer**: "Choose fanout/levels so rewrite rate < 300MB/s slack; e.g., fanout 10, L0 cap and parallel compactions to maintain steady‑state without backlog."
