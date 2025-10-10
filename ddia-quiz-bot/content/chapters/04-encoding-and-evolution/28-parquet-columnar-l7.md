---
id: ch04-parquet-columnar-l7
day: 28
level: L7
tags: [parquet, columnar-storage, analytics, data-warehouse, principal-engineer]
related_stories: []
---

# Parquet and the Row vs Column Storage Decision

## question
Your company stores petabytes of event data in S3 as JSON files (one event per line). The data analytics team runs hundreds of Spark jobs daily, but query costs and times are spiraling. The proposed solution is migrating to Parquet columnar format. Beyond the obvious compression benefits, what architectural insights about query patterns and data access justify columnar storage, and when would columnar formats be the wrong choice?

## expected_concepts
- Columnar vs row-oriented storage trade-offs
- Query selectivity and projection pushdown
- Compression efficiency on homogeneous data
- Predicate pushdown and min/max statistics
- OLTP vs OLAP workload characteristics
- Parquet's schema evolution model
- Write amplification vs read optimization

## answer
The fundamental insight: Columnar storage optimizes for analytical queries that read many rows but few columns, while row storage optimizes for transactional queries that read/write all columns for few rows.

Why Parquet wins for analytics: (1) Column projection - reading only needed columns scans less data (if query needs 3 of 50 fields, read 3 columns not 50), (2) Homogeneous compression - columns have similar types/values, enabling better compression ratios (dictionary encoding, run-length encoding), (3) Predicate pushdown - column statistics (min/max, bloom filters) enable skipping irrelevant row groups without decompression, (4) SIMD/vectorized processing - columns fit in CPU cache for batch operations. For analytics on S3, this reduces both data scanned (lower costs) and query time.

When columnar is wrong: (1) Record-oriented access patterns - if queries need all columns for each row (e.g., serving user profiles), row format is better, (2) High write throughput with low latency requirements - columnar formats require buffering and organizing into column chunks, adding write latency, (3) Frequent updates to individual records - columnar storage doesn't support efficient in-place updates, (4) Schema volatility with many sparse columns - columnar format has overhead for each column.

Architectural decision: Use Parquet for immutable, append-only analytics (data lake, warehouse staging), use row formats (JSON, Avro) for operational data and frequent updates. Hybrid architectures are common: Lambda architecture with row format in speed layer, columnar in batch layer.

## hook
Why can Parquet achieve 10-100x compression compared to JSON beyond just gzip compression?

## follow_up
After migrating to Parquet, the data science team loves the query performance, but the real-time ML feature serving team complains that reading a single user's features (1 row with 500 columns) from Parquet is slower than from the old JSON format. They want to revert. How do you architect a solution that serves both workloads without maintaining two complete copies of the data?

## follow_up_answer
This is a classic case of access pattern mismatch requiring a polyglot persistence strategy: (1) Parquet in S3 remains the source of truth for analytical queries (batch feature engineering, historical analysis), (2) Implement a serving layer that materializes point-read access patterns - options include key-value stores (DynamoDB, Redis) indexed by user_id with row-oriented storage, or specialized feature stores (Feast, Tecton) that handle this pattern, (3) Use change data capture (CDC) or batch jobs to sync Parquet data to serving layer, optimizing for your freshness requirements.

Architectural insight: The same data may need multiple physical representations optimized for different access patterns - this is exactly what DDIA Chapter 3 teaches about derived data. Don't force one storage format to serve all workloads.

Alternative for cost-constrained scenarios: Use Apache ORC instead of Parquet - ORC has "row groups" with local indexing that makes point queries somewhat faster while maintaining columnar benefits for analytics. Or partition Parquet by user_id cohorts so reading one user only scans a small file. But these are compromises - proper solution is acknowledging different workloads need different storage.

The meta-principle: Storage format is not just about compression - it's about optimizing for your query patterns. Columnar for "many rows, few columns" (OLAP), row-oriented for "few rows, all columns" (OLTP), key-value for "one row, by key" (serving).
