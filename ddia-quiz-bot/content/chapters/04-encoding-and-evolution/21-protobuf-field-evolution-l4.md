---
id: ch04-protobuf-field-evolution-l4
day: 21
level: L4
tags: [protobuf, field-evolution, breaking-changes, practical]
related_stories: []
---

# Protocol Buffers Field Type Changes

## question
You have a Protocol Buffer message with `int32 user_id = 1;` deployed in production. Your team wants to change it to `int64 user_id = 1;` to support larger IDs. What happens to existing encoded data?

## options
- A) Old data with int32 values can still be read correctly as int64
- B) All existing data becomes unreadable and must be migrated
- C) The data is automatically converted during deserialization
- D) It works only if you increment the schema version number

## answer
A

## explanation
Protocol Buffers allow certain type changes without breaking compatibility. Changing int32 to int64 is safe because the wire format is compatible - int32 values fit within int64's range. The smaller values are zero-extended. However, the reverse (int64 to int32) would truncate values and potentially lose data. This is a practical example of schema evolution in action.

## hook
What other field type changes are safe in Protobuf without breaking compatibility?
