---
id: storage-subjective-L3-003
type: subjective
level: L3
category: baseline
topic: storage-and-retrieval
subtopic: bloom-filters-lsm
estimated_time: 5-7 minutes
---

# question_title - Bloom Filters in LSM Trees

## main_question - Core Question
"What is a Bloom filter and how does it reduce read amplification in LSM‑tree storage? Explain false positives/negatives, sizing trade‑offs, and where filters are placed."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Bloom Filter**: Probabilistic set membership
- **No False Negatives**: May say ‘maybe’; never misses existing
- **False Positives**: Trade space for lower FP rate
- **Placement**: Per‑SSTable (and sometimes per‑block) to skip disk reads

### expected_keywords
- Primary: false positive rate (FPR), hash functions, bit array
- Technical: k hashes, m bits, n items, FPR≈(1−e^{−kn/m})^k

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Tiered vs Leveled**: Different number of tables to probe
- **Memory Budget**: Filters live in RAM; cache efficiency
- **Point vs Range**: Filters help point lookups more

### bonus_keywords
- Implementation: per‑file metadata, block filters, prefix bloom
- Related: Cuckoo filters (deletions), quotient filters

## sample_excellent - Example Excellence
"A Bloom filter is a compact bitset with k hash functions used to test set membership: it can return false positives but never false negatives. In LSM stores, each SSTable has a Bloom filter for its key set; on a point lookup, if the filter says ‘definitely not present’, the engine skips opening that SSTable—eliminating disk I/Os and reducing read amplification. Tuning the bit‑per‑key budget sets the FPR: more bits and well‑chosen k reduce FPR. Filters are typically stored in table metadata and kept hot in memory; some engines add per‑block or prefix filters to further prune reads."

## sample_acceptable - Minimum Acceptable
"Bloom filters quickly tell you if a key is definitely not in an SSTable, so you can skip reading that file. They can be wrong in saying ‘maybe’, but won’t miss real keys."

## common_mistakes - Watch Out For
- Claiming Bloom filters can have false negatives
- Forgetting the memory budget trade‑off
- Assuming filters help range scans as much as point lookups

## follow_up_excellent - Depth Probe
**Question**: "Given a 1% FPR and 10 levels, estimate how many extra SSTables you’ll read on a miss. How does leveled vs tiered compaction change this?"
- **Looking for**: Roughly 0.1 tables per level on average; tiered has more tables than leveled

## follow_up_partial - Guided Probe  
**Question**: "Why keep filters in memory and how large should they be?"
- **Hint**: RAM avoids extra I/Os; budget bits per key vs FPR

## follow_up_weak - Foundation Check
**Question**: "If a filter says ‘no’, what does that mean?"
- **Simplification**: Definitely not in that file

## bar_raiser_question - L3→L4 Challenge
"You have 1GB RAM for filters and 1B keys across all SSTables. Propose a bit‑per‑key budget, estimate FPR, and explain the impact on tail latency for point misses."

### bar_raiser_concepts
- Back‑of‑envelope FPR math and latency reasoning
- Sensitivity to level count and device IOPS

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 2-3 min answer + 3-4 min discussion
- **Common next topics**: Compaction strategy, cache design, I/O scheduling

## assistant_answer
Bloom filters are per‑SSTable probabilistic indexes that rule out absent keys with no false negatives, cutting pointless disk probes; tuning bits per key and k hashes balances memory vs FPR and thus tail read latency.

## improvement_suggestions
- Ask for concrete FPR calculations and a simple latency model (tables probed × device latency).
- Include discussion of per‑block filters and prefix filters for range‑like lookups.

## improvement_exercises
### exercise_1 - FPR Math
**Question**: "For m/n=10 bits per key and k≈(m/n)ln2≈7, estimate FPR using (1−e^{−kn/m})^k."

**Sample answer**: "FPR≈(1−e^{−7/10})^7≈(1−e^{−0.7})^7≈(1−0.4966)^7≈(0.5034)^7≈0.0078 (≈0.78%)."

### exercise_2 - Memory Budgeting
**Question**: "With 1GB for filters and 1B keys, what’s bits/key and approximate FPR?"

**Sample answer**: "≈8 bits/key → around 1–2% FPR with optimal k; tune to workload/tail goals."
