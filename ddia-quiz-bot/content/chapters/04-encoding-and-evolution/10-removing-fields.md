---
id: ch04-removing-fields
day: 10
tags: [schema-evolution, field-removal, compatibility]
related_stories: []
---

# Removing Fields from Schema

## question
What rule must you follow when removing a field from a schema?

## options
- A) You can never remove fields, only deprecate them
- B) You can only remove optional fields, and you must never reuse the field tag/name
- C) You must remove all data containing that field first
- D) You can remove any field as long as you increment the schema version

## answer
B

## explanation
You can only remove optional fields (required fields can never be removed). After removing a field, you must never reuse that field tag number (in Protobuf/Thrift) or name (in Avro), because you may still have data stored somewhere that includes the old field, and reusing the tag/name would cause misinterpretation.

## hook
Why can't you reuse field tag 42 after removing the field that had that tag?
