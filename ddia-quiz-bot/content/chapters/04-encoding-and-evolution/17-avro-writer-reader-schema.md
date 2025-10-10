---
id: ch04-avro-writer-reader-schema
day: 17
tags: [avro, writer-schema, reader-schema, schema-resolution]
related_stories: []
---

# Avro Writer and Reader Schema

## question
How does Avro handle the scenario where the writer's schema and reader's schema are different versions?

## options
- A) It throws an error and refuses to read the data
- B) It uses schema resolution to translate data from writer's schema to reader's schema
- C) It always uses the newer schema version
- D) It requires a separate migration tool

## answer
B

## explanation
Avro uses schema resolution to translate data between writer's and reader's schemas. The reader looks at both schemas and translates field-by-field. Fields that appear in the writer's schema but not the reader's are ignored. Fields in the reader's schema but not the writer's are filled with default values. Fields must have matching types.

## hook
How does your system handle data written two years ago with an old schema version?
