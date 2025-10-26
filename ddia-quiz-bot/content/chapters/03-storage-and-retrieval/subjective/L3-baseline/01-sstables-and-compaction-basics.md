---
id: storage-subjective-L3-001
type: subjective
level: L3
category: baseline
topic: storage-and-retrieval
subtopic: sstables-compaction-basics
estimated_time: 5-7 minutes
---

# question_title - SSTables and Compaction Basics

## main_question - Core Question
"Explain what SSTables are and how compaction works in an LSM-based store. Why do we keep only the most recent value per key after compaction, and what invariants hold?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **SSTables**: Immutable, sorted segments on disk
- **Merge Compaction**: Merge multiple sorted runs into fewer
- **Latest-Write-Wins**: Deduplicate duplicates by key; keep newest
- **Tombstones**: Deletes represented and purged after grace/TTL

### expected_keywords
- Primary: SSTable, compaction, merge, deduplication
- Technical: key ordering, sequence number/timestamp, tombstone

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Read/Write Amplification**: Compaction trades write amp for read perf
- **Levels/Tiers**: Size-tiered vs leveled compaction
- **Range Queries**: Benefit from sorted order
- **Bloom Filters**: Skip non-existent keys fast

### bonus_keywords
- Implementation: memtable flush, manifest, index blocks
- Related: LSM-tree, fence pointers

## sample_excellent - Example Excellence
"SSTables are immutable, sorted files produced by flushing a memtable. Compaction merges several SSTables by key order, emitting only the newest value per key (using sequence numbers) and honoring tombstones to remove deleted keys once older versions are no longer visible. The result is fewer, larger sorted runs, reducing read amplification. The invariants are: within any SSTable, keys are strictly ordered; across a level (in leveled compaction), ranges don’t overlap; for any key, at most one live version exists after compaction."

## sample_acceptable - Minimum Acceptable
"SSTables are sorted files. Compaction merges them and keeps the latest value per key, removing older duplicates and deletes after a while. This reduces the number of files to search."

## common_mistakes - Watch Out For
- Thinking compaction preserves all history by default
- Ignoring tombstone grace periods
- Confusing compaction (reorganization) with compression (encoding)
- Forgetting that files are immutable once written

## follow_up_excellent - Depth Probe
**Question**: "Contrast size-tiered vs leveled compaction in terms of write/read amplification and space amplification."
- **Looking for**: Size-tiered → lower write amp, higher read/space; Leveled → higher write amp, lower read/space

## follow_up_partial - Guided Probe  
**Question**: "How do tombstones get purged safely?"
- **Hint**: Wait until all older SSTables that might contain older versions are compacted

## follow_up_weak - Foundation Check
**Question**: "Why does sorting help with range scans?"
- **Simplification**: Adjacent keys stored together

## bar_raiser_question - L3→L4 Challenge
"You have three SSTables with overlapping key ranges and per-entry sequence numbers. Design the merge procedure to guarantee newest-value wins and safe tombstone purge. Detail tie-breaking and file ordering."

### bar_raiser_concepts
- Stable merge by key, prefer higher sequence number
- Track tombstone horizons; purge when older files no longer contribute
- Maintain monotonic manifest of active files

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 2-3 min answer + 3-4 min discussion
- **Common next topics**: Compaction strategies, Bloom filters, caching

## assistant_answer
SSTables are immutable, sorted files; compaction merges runs by key and keeps only the newest version (by sequence/timestamp), applying tombstones to delete keys once older versions can’t reappear. This reduces read amplification and preserves the invariant of at most one live value per key.

## improvement_suggestions
- Require the candidate to enumerate invariants (ordering, non-overlap, single live version) explicitly.
- Add a small worked example with three files and per-key sequence numbers to trace the merge.

## improvement_exercises
### exercise_1 - Merge Walkthrough
**Question**: "Given SSTables S1:{(a,1),(b,3T),(c,2)}, S2:{(a,2T),(b,2),(d,1)}, S3:{(b,4),(e,T)} where numbers are seq and T marks tombstone, produce the merged output."

**Sample answer**: "Keep highest seq per key, apply tombstones: a→2T (tombstone newest), b→4, c→2, d→1, e→T. Depending on purge policy, keys with T may be retained until safe to drop."

### exercise_2 - Strategy Choice
**Question**: "When would you choose leveled over size-tiered compaction for an OLTP workload?"

**Sample answer**: "When reads dominate or low tail read latency is needed: leveled reduces read and space amplification at the cost of higher write amplification."
