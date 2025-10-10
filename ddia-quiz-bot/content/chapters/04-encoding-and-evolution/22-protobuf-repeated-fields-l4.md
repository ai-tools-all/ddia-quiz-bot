---
id: ch04-protobuf-repeated-fields-l4
day: 22
level: L4
tags: [protobuf, repeated-fields, optional-fields, compatibility]
related_stories: []
---

# Protocol Buffers Repeated vs Optional Fields

## question
You have `optional string email = 3;` in your Protobuf schema and want to change it to `repeated string email = 3;` to support multiple emails per user. Is this change backward compatible?

## options
- A) Yes, old data with single email is read as a list with one element
- B) No, this is a breaking change requiring a new field tag
- C) Yes, but only if you add a migration script
- D) It depends on the Protobuf version being used

## answer
A

## explanation
Changing optional to repeated is backward compatible in Protocol Buffers. Old data with a single value is interpreted as a repeated field with one element. However, changing repeated to optional is NOT safe - if old data has multiple values, only the last one would be kept. This demonstrates understanding of Protobuf's wire format and evolution rules.

## hook
What happens if new code writes multiple emails and old code reads the message?
