---
id: ch04-schema-advantages
day: 3
tags: [schema, binary-encoding, thrift, protobuf, avro]
related_stories: []
---

# Schema-Based Binary Encoding Advantages

## question
What is a key advantage of schema-based binary encoding formats (Thrift, Protocol Buffers, Avro) over schema-less formats?

## options
- A) They don't require any documentation
- B) Field names can be omitted from encoded data, using field tags instead
- C) They work with any programming language without code generation
- D) They are human-readable in hex editors

## answer
B

## explanation
Schema-based binary formats use field tags (numbers) instead of field names in the encoded data, making the binary representation much more compact. The schema defines the mapping between field tags and names, which is kept separate from the data itself.

## hook
How much space does your JSON API waste by repeating field names in every record?
