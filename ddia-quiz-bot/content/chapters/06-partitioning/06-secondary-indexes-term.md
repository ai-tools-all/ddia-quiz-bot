---
id: ch06-secondary-indexes-term
day: 6
tags: [partitioning, secondary-indexes, term-partitioning, global-indexes]
related_stories: []
---

# Term-Based Secondary Indexes

## question
What is the primary advantage of term-based partitioning for secondary indexes (global indexes)?

## options
- A) Simpler implementation than document-based indexes
- B) Read queries can be served from a single partition
- C) Writes are always faster
- D) No need for index maintenance

## answer
B

## explanation
Term-based partitioning creates a global index where each term (index entry) is assigned to a specific partition based on the term itself. This means a read query for a specific indexed value only needs to contact the partition(s) responsible for that term, making reads efficient. However, writes become more complex as they may need to update multiple partitions if the document has multiple indexed fields.

## hook
How do global secondary indexes in DynamoDB handle the write amplification problem?
