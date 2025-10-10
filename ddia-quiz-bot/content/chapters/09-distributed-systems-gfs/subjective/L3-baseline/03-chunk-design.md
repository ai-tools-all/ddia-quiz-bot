---
id: gfs-subjective-L3-003
type: subjective
level: L3
category: baseline
topic: gfs
subtopic: chunk-architecture
estimated_time: 5-7 minutes
---

# question_title - GFS Chunk Size Decision

## main_question - Core Question
"GFS uses 64MB chunks while traditional file systems use 4KB blocks. Why such a huge difference? What are the implications?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Metadata Reduction**: Fewer chunks means less metadata for master
- **Network Efficiency**: Reduces overhead for large file operations
- **Client-Master Interaction**: Fewer requests to master

### expected_keywords
- Primary keywords: chunk size, metadata, master, efficiency
- Technical terms: 64MB chunks, block size

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Hot Spots**: Small files problem with large chunks
- **Internal Fragmentation**: Wasted space for small files
- **Persistent TCP Connections**: Reuse for large transfers
- **Seek Time Amortization**: Large sequential reads

### bonus_keywords
- Implementation: chunk handle, lazy allocation
- Trade-offs: space efficiency, load balancing
- Comparisons: HDFS block size, traditional FS

## sample_excellent - Example Excellence
"GFS chose 64MB chunks (vs 4KB traditional blocks) primarily to minimize metadata overhead. With billions of files, using small blocks would require the master to store enormous amounts of metadata in memory. Large chunks also reduce client-master interactions - a client reading a large file needs fewer chunk location requests. Additionally, 64MB chunks improve network efficiency by reducing overhead and allowing persistent TCP connections for transfers. However, this creates problems for small files which become hot spots if many clients access them, since all requests hit the same chunk servers."

## sample_acceptable - Minimum Acceptable
"The 64MB chunk size reduces the amount of metadata the master needs to track and decreases the number of times clients need to contact the master for chunk locations. This helps the single master design scale better."

## common_mistakes - Watch Out For
- Not understanding metadata implications
- Missing the single master scalability angle
- Forgetting about small file problems
- Confusing chunks with file size

## follow_up_excellent - Depth Probe
**Question**: "How would you handle a system with millions of 1KB configuration files using GFS's architecture?"
- **Looking for**: Higher replication for hot files, caching strategies, alternative storage
- **Red flags**: Suggesting to change chunk size globally

## follow_up_partial - Guided Probe
**Question**: "You mentioned metadata. Can you calculate roughly how much memory the master needs if we have 1 petabyte of data?"
- **Hint embedded**: Math with 64MB chunks
- **Concept testing**: Scale comprehension

## follow_up_weak - Foundation Check
**Question**: "Think about organizing a library. Would you prefer having many small boxes or fewer large boxes for books? What are the trade-offs?"
- **Simplification**: Physical analogy
- **Building block**: Chunking benefits

## bar_raiser_question - L3→L4 Challenge
"Design a modification to GFS that handles both large video files (10GB each) and millions of small thumbnail images (100KB each) efficiently. What changes would you make?"

### bar_raiser_concepts
- Variable chunk sizes or tiered storage
- Metadata optimization strategies
- Separate systems for different workloads
- Cost-benefit analysis of complexity

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 2-3 min answer + 3-4 min discussion
- **Common next topics**: HDFS comparison, hot spot mitigation, metadata scaling
