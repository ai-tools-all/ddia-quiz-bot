---
id: ch04-schema-evolution-rules
day: 9
tags: [schema-evolution, optional-fields, default-values]
related_stories: []
---

# Schema Evolution Rules

## question
To maintain backward and forward compatibility when evolving a schema, what must be true about new fields you add?

## options
- A) New fields must be required with no default value
- B) New fields must be optional or have a default value
- C) New fields must have the same type as existing fields
- D) New fields cannot be added, only modified

## answer
B

## explanation
New fields must be optional or have a default value to maintain compatibility. For backward compatibility, old code (which doesn't know about the new field) needs to be able to read new data - it simply ignores the new field. For forward compatibility, new code must be able to read old data - the default value is used when the field is missing.

## hook
What's your strategy for adding required fields to a deployed schema?
