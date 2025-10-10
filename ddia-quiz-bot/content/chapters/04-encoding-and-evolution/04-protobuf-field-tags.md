---
id: ch04-protobuf-field-tags
day: 4
tags: [protobuf, field-tags, schema, evolution]
related_stories: []
---

# Protocol Buffers Field Tags

## question
In Protocol Buffers, what happens if you change a field tag number in the schema?

## options
- A) The field is automatically migrated to the new tag
- B) Old data becomes unreadable or incorrectly interpreted
- C) Protocol Buffers automatically creates an alias
- D) The data is re-encoded on the fly

## answer
B

## explanation
Field tags are critical in Protocol Buffers - they identify fields in the encoded binary data. If you change a field tag, old encoded data will either become unreadable or the field will be interpreted as a different field. Field tags must remain stable to maintain compatibility.

## hook
Why are field tags in Protobuf more important than field names?
