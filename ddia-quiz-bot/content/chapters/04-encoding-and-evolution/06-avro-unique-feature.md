---
id: ch04-avro-unique-feature
day: 6
tags: [avro, schema, writer-schema, reader-schema]
related_stories: []
---

# Avro's Unique Schema Approach

## question
What makes Avro different from Thrift and Protocol Buffers in terms of schema handling?

## options
- A) Avro doesn't support schema evolution
- B) Avro has no tag numbers in the schema; fields are identified by name
- C) Avro requires the schema to be embedded in every record
- D) Avro only works with JSON encoding

## answer
B

## explanation
Unlike Thrift and Protocol Buffers which use field tags/numbers, Avro identifies fields by name alone. This makes Avro schemas more friendly to dynamically generated schemas. Avro resolves differences between writer's schema and reader's schema by matching field names.

## hook
How does Avro encode data without field tags and still maintain compact size?
