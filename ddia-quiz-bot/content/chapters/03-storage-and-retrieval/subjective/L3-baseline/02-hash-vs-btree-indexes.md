---
id: storage-subjective-L3-002
type: subjective
level: L3
category: baseline
topic: storage-and-retrieval
subtopic: hash-vs-btree-indexes
estimated_time: 5-7 minutes
---

# question_title - Hash Indexes vs B‑Tree Indexes

## main_question - Core Question
"Compare hash indexes and B‑tree indexes. When is each appropriate, and why are hash indexes unsuitable for range queries while B‑trees excel at them?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Hash Index**: Key→bucket mapping, O(1) point lookups
- **B‑Tree**: Ordered, balanced tree over key space
- **Range Queries**: Require order; hash destroys order
- **Locality**: B‑tree preserves key locality for scans

### expected_keywords
- Primary: hash, B‑tree, range, point lookup, order
- Technical: buckets, collisions, fan‑out, pages/blocks

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Write Patterns**: Random vs append‑friendly (LSM alternative)
- **Space/Cache**: Fill factor, page splits, rehashing
- **Workload Fit**: Equality heavy vs range heavy

### bonus_keywords
- Implementation: clustered vs secondary index, prefix compression
- Related: LSM‑tree, skip list, fractal tree

## sample_excellent - Example Excellence
"Hash indexes map keys to buckets, giving near‑constant‑time point lookups, but they destroy key order: neighbors in key space land in unrelated buckets. Thus you cannot iterate keys in sorted order or efficiently answer range queries with a hash alone. B‑trees store keys in sorted order across fixed‑size pages with high fan‑out, making point lookups logarithmic and range scans efficient via in‑order traversal and page locality. Use hashes for equality‑dominated workloads (e.g., caches, exact key/value stores) and B‑trees when you need ordering: prefix/range scans, ORDER BY, MIN/MAX."

## sample_acceptable - Minimum Acceptable
"Use hash indexes for exact lookups; B‑trees support both exact and range lookups because they keep keys ordered. Hashes don’t preserve order."

## common_mistakes - Watch Out For
- Claiming hash indexes can do range scans efficiently
- Ignoring page splits and rebalancing costs in B‑trees
- Assuming B‑trees are always faster than hashing for equality

## follow_up_excellent - Depth Probe
**Question**: "Describe how clustered vs non‑clustered B‑tree indexes affect range scan performance."
- **Looking for**: Clustered aligns data with index order → fewer random I/Os; non‑clustered needs lookups/rid fetches

## follow_up_partial - Guided Probe  
**Question**: "Why might a hash index outperform a B‑tree for exact lookups in memory, but not on SSD/HDD?"
- **Hint**: Cache locality, page/block I/O granularity

## follow_up_weak - Foundation Check
**Question**: "What does it mean that hashing ‘destroys order’?"
- **Simplification**: Nearby keys become unrelated numbers

## bar_raiser_question - L3→L4 Challenge
"Given a workload with 90% point reads, 10% range scans (100–1000 keys), and mixed read/write ratio 70/30 on SSD, pick an index strategy and justify it. Include maintenance/space and tail latency."

### bar_raiser_concepts
- B‑tree (clustered) or LSM with secondary range index
- Trade space/maintenance vs scan cost
- Consider compaction or page split overhead

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 2-3 min answer + 3-4 min discussion
- **Common next topics**: LSM vs B‑tree, secondary indexes, bloom filters

## assistant_answer
Hash indexes are ideal for equality lookups but unsuitable for ranges because hashing removes order; B‑trees keep keys ordered, enabling efficient scans while still supporting point lookups with logarithmic cost and good locality.

## improvement_suggestions
- Ask for a decision matrix mapping workload characteristics (range %, write %, device) to index choices.
- Include a brief discussion of clustered vs non‑clustered impacts on scans.

## improvement_exercises
### exercise_1 - Decision Matrix
**Question**: "Fill a simple matrix recommending hash vs B‑tree (or LSM) for combinations of range share {0%,10%,50%} and write share {10%,50%,90%}."

**Sample answer**: "0% range/10% writes → hash; 10% range/50% writes → B‑tree or LSM+sec idx; 50% range/90% writes → LSM with leveled compaction and range‑capable index."

### exercise_2 - Clustered Impact
**Question**: "Explain how a clustered B‑tree reduces random I/O for range queries vs non‑clustered."

**Sample answer**: "Data rows are stored in index order, so a range scan reads contiguous pages; non‑clustered requires lookups from index to data pages, increasing random reads."
