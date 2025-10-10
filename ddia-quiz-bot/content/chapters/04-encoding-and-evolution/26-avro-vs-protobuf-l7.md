---
id: ch04-avro-vs-protobuf-l7
day: 26
level: L7
tags: [avro, protobuf, architectural-choice, data-platform, principal-engineer]
related_stories: []
---

# Avro vs Protobuf: When Dynamic Schemas Matter

## question
LinkedIn chose Avro over Protocol Buffers for Kafka message serialization, while Google uses Protobuf extensively. Both companies have sophisticated engineering teams. What architectural insight drives this divergence, and when would you advocate for each at an organizational level?

## expected_concepts
- Dynamic schema generation and schema evolution
- Schema-on-write vs schema-on-read
- Code generation requirements and language polyglotism
- Data lake/warehouse integration patterns
- Organizational structure and Conway's Law
- Streaming vs RPC workload characteristics

## answer
The key insight: Avro excels when schemas are dynamically generated or frequently evolve without code changes, while Protobuf excels when schemas are statically defined with strong typing across compiled languages.

LinkedIn/Kafka use case: Data pipelines involve hundreds of data sources with varying schemas, often generated from database schemas automatically. Avro's lack of field tags and name-based resolution makes dynamic schema generation trivial. Schema evolution is handled at read time (schema-on-read), perfect for data lakes where consumers determine interpretation. No code generation required for schema changes - critical for data engineering workflows.

Google/gRPC use case: Service-to-service RPC needs strong typing and code generation for type safety across multiple languages (C++, Java, Go). Field tags provide stability and compact encoding. Schema changes are deliberate API changes requiring code review and deployment coordination - the ceremony is a feature, not a bug.

Architectural decision criteria: Choose Avro for data-centric, schema-flexible workloads (data lakes, event logs, ETL pipelines). Choose Protobuf for service-centric, type-safe RPC (microservices, mobile APIs). Hybrid approach: Protobuf for service APIs, Avro for event streams - acknowledges different needs.

## hook
Why is Avro's lack of field tags a strength for data engineering but a weakness for RPC?

## follow_up
Your organization has standardized on Protobuf for all internal communication. A data engineering team now wants to build a real-time feature store that ingests from 50+ microservices and serves features to ML models. They argue that maintaining 50+ Protobuf schemas and regenerating code on every schema change is unsustainable. How do you resolve this architecturally without abandoning your Protobuf standard?

## follow_up_answer
This is a classic case where architectural boundaries need schema translation layers: (1) Keep Protobuf for service-to-service RPC (maintains type safety, explicit contracts), (2) Implement a schema translation layer at the feature store boundary that converts Protobuf to Avro for internal storage and processing, (3) Use Protobuf reflection or descriptors to automate the conversion without manual code generation per schema change, (4) Feature store consumers (ML models) use schema-on-read with Avro, which suits their dynamic needs. Key insight: Don't force one encoding format across all use cases - honor the boundary between operational systems (need consistency, type safety) and analytical systems (need flexibility, schema evolution). This reflects the broader principle that batch/streaming analytics have fundamentally different requirements than transactional services.
