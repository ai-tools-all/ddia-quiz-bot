---
id: craq-subjective-L6-003
type: subjective
level: L6
category: baseline
topic: craq
subtopic: consistency-monitoring
estimated_time: 9-12 minutes
---

# question_title - Building a CRAQ Consistency Monitoring Service

## main_question - Core Question
"Architect a consistency auditing service that continuously validates CRAQ's linearizability by sampling reads and comparing them against a secondary analytical store. Demonstrate how you would avoid false positives using DDIA's guidance on data pipelines and probabilistic checks." 

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Sampling Strategy**: Trace reads with metadata (epoch, version) to verify against authoritative tail snapshot
- **Analytical Store Sync**: Use tail commit index to drive incremental ETL without lag-induced skew
- **False Positive Control**: Account for propagation delay via watermarks
- **DDIA Link**: Borrow streaming/batch integration patterns and validation techniques

### expected_keywords
- Primary keywords: auditing, watermark, sampling, analytical store
- Technical terms: ETL pipeline, watermark, tail commit index, probabilistic checking

## peripheral_concepts - Nice to Have (40%)
- **Anomaly Classification**: Distinguish stale read vs pipeline lag vs data corruption
- **Feedback Loop**: Alert to reliability engineers, tie into error budgets
- **Security Considerations**: Access control to sensitive data snapshots
- **Operationalization**: Scheduled reconciliation jobs, streaming join window

### bonus_keywords
- Implementation: Kafka, Flink, batch snapshots, Bloom filters, change queries
- Scenarios: nightly reconciliation, on-demand audits, regulatory reports
- Trade-offs: audit accuracy vs cost, sampling vs exhaustive verification

## sample_excellent - Example Excellence
"We instrument CRAQ clients to attach epoch + tail index to each read. The auditing service subscribes to CDC and ingests tail-confirmed writes into an analytical store with watermarks. A streaming job compares sampled reads against this store, allowing for a bounded watermark delay equal to clean propagation time to avoid false positives. Deviations beyond that window indicate consistency issues. This mirrors DDIA's recommendation to combine streaming with batch snapshots to validate data quality." 

## sample_acceptable - Minimum Acceptable
"Track read metadata (epoch, tail index), compare it against a tail-driven analytical store using a watermark equal to propagation delay, and alert on mismatches—borrowing DDIA's streaming + batch validation ideas." 

## common_mistakes - Watch Out For
- Comparing against stale analytical data without watermarks
- Ignoring propagation delay leading to false alarms
- No sampling or metadata to trace individual reads
- Not referencing DDIA's data pipeline practices

## follow_up_excellent - Depth Probe
**Question**: "How do you tune the sampling rate to balance detection time vs overhead?"
- **Looking for**: Adaptive sampling, risk-based prioritization, cost-benefit analysis
- **Red flags**: Static arbitrary value without justification

## follow_up_partial - Guided Probe  
**Question**: "What additional metadata lets you distinguish between a stale read and a replayed write?"
- **Hint embedded**: Request ID, duplicate suppression counters
- **Concept testing**: Diagnostic depth

## follow_up_weak - Foundation Check
**Question**: "Why do auditors use random spot checks instead of redoing the entire ledger every day?"
- **Simplification**: Sampling rationale
- **Building block**: Efficiency vs certainty

## bar_raiser_question - L6→L7 Challenge
"Extend the auditing service to supply regulatory proof of linearizability with cryptographic attestations, combining DDIA's log provenance advice with CRAQ metadata." 

### bar_raiser_concepts
- Cryptographic commitments, Merkle trees
- Regulatory compliance, provable correctness
- Attestation pipeline integration
- Data lineage tracking

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 5-6 min answer + 4-6 min discussion
- **Common next topics**: Compliance, data validation, observability
