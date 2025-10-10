---
id: gfs-subjective-L5-bar
type: subjective
level: L5
category: bar-raiser
topic: gfs
subtopic: cross-system-integration
estimated_time: 12-15 minutes
---

# question_title - GFS in Modern Data Pipeline

## main_question - Core Question
"Design a complete data pipeline using GFS as the storage layer: real-time events from millions of IoT devices → stream processing → batch analytics → ML training → serving predictions. How do you handle the impedance mismatch between GFS's design and modern streaming/serving requirements?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Buffering Layer**: Bridge between streaming and batch storage
- **Format Optimization**: Different formats for different stages
- **Consistency Boundaries**: Where strong consistency is needed
- **Performance Isolation**: Preventing pipeline stages from interfering

### expected_keywords
- Primary keywords: streaming, batch, lambda architecture, data lake
- Technical terms: Kafka, Flink, Parquet, feature store

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Exactly-Once Semantics**: End-to-end delivery guarantees
- **Schema Evolution**: Handling changing data formats
- **Multi-Temperature Storage**: Hot/warm/cold data tiers
- **Backpressure Handling**: Flow control across systems
- **Data Lineage**: Tracking transformations
- **Cost Optimization**: Compute vs storage trade-offs

### bonus_keywords
- Technologies: Arrow, Iceberg, Delta Lake, Hudi
- Patterns: Kappa architecture, medallion architecture
- Optimization: Predicate pushdown, columnar formats

## sample_excellent - Example Excellence
"To bridge GFS's batch-optimized design with modern requirements: Architecture: 1) Ingestion: Kafka for streaming ingestion, buffering events until GFS's 64MB chunks are efficient. Write-ahead log pattern for durability. 2) Stream→Batch: Flink writes micro-batches to GFS in Parquet format every 5 minutes, balancing latency vs efficiency. 3) Storage organization: Partitioned data lake with bronze/silver/gold tiers. Bronze=raw, Silver=cleaned, Gold=aggregated. 4) ML Pipeline: Direct GFS reads for training, but feature store (Redis/DynamoDB) for serving to avoid GFS latency. 5) Serving: Predictions cached in CDN, only fallback to GFS for cold data. Key adaptations: GFS append-only for streaming writes, external metadata catalog (Hive/Iceberg) for ACID semantics where needed, separate low-latency store for serving. This accepts GFS's limitations rather than fighting them, using it where it excels (batch storage) and complementing with specialized systems elsewhere."

## sample_acceptable - Minimum Acceptable
"Use Kafka to buffer streaming data before writing to GFS in batches. Store raw data in GFS for batch processing and ML training. For serving, cache processed results in a faster database since GFS isn't optimized for low-latency random reads. Use different storage formats optimized for each use case."

## common_mistakes - Watch Out For
- Trying to use GFS for real-time serving
- Not addressing streaming-batch impedance
- Ignoring format optimization
- Missing consistency requirements

## follow_up_excellent - Depth Probe
**Question**: "Your pipeline processes PII data with GDPR requirements for deletion. How do you implement 'right to be forgotten' with GFS's append-only design?"
- **Looking for**: Tombstoning, compaction strategies, crypto-shredding
- **Red flags**: Saying it's impossible

## follow_up_partial - Guided Probe
**Question**: "You mentioned using Kafka for buffering. How do you handle Kafka failures without losing data or creating duplicates in GFS?"
- **Hint embedded**: Transactional writes, idempotency
- **Concept testing**: Distributed transaction understanding

## follow_up_weak - Foundation Check
**Question**: "Let's simplify: You have a temperature sensor sending readings every second. How would you efficiently store a year of data?"
- **Simplification**: Single data source
- **Building block**: Batching concept

## bar_raiser_question - L6→L7 Challenge
"Now make this pipeline multi-region with active-active processing, maintaining global consistency for financial transactions while optimizing for local performance. Design the architecture."

### bar_raiser_concepts
- Geo-distributed consensus
- Conflict-free replicated data types (CRDTs)
- Region-aware routing
- Cross-region consistency protocols
- Disaster recovery with RPO/RTO guarantees
- Regulatory compliance across jurisdictions

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 5-7 min answer + 7-8 min discussion
- **Common next topics**: Modern data platforms, streaming systems, ML infrastructure
