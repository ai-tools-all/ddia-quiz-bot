---
id: ch06-secondary-indexes-document
day: 5
tags: [partitioning, secondary-indexes, document-partitioning]
related_stories: []
---

# Document-Based Secondary Indexes

## question
In document-based partitioning of secondary indexes (local indexes), what is the main drawback?

## options
- A) Indexes consume more storage space
- B) Write operations become slower
- C) Read queries must scatter to all partitions
- D) Index updates require distributed transactions

## answer
C

## explanation
With document-based partitioning, each partition maintains its own secondary indexes for the documents it stores. While this makes writes simple and fast (only update the local partition), read queries using secondary indexes must query all partitions and merge results (scatter-gather), which is expensive and increases latency. This is why it's sometimes called a "scatter-gather" approach.

## hook
Why does Elasticsearch default to document-based partitioning despite its read penalty?
