---
id: ch04-protobuf-default-values-l5
day: 24
level: L5
tags: [protobuf, default-values, api-design, semantics]
related_stories: []
---

# Protocol Buffers Default Values Problem

## question
Your Protobuf message has `bool is_premium = 5;` to indicate premium users. In the encoded message, you cannot distinguish between a user explicitly set to `false` and a user where the field was never set (defaults to false). How does this impact API design?

## options
- A) This is not a problem; false is false regardless of how it got there
- B) Use an optional wrapper or separate field for "has been set" semantics when the distinction matters
- C) Always use string fields instead of bool to avoid this issue
- D) Set premium users to false and non-premium to true to invert the logic

## answer
B

## explanation
Protocol Buffers' default values create a semantic ambiguity: you can't distinguish "explicitly set to default" from "never set". For booleans, 0 for numbers, or empty strings, this matters when the presence/absence of a field has meaning different from its default value. Solutions: use wrapper types (google.protobuf.BoolValue), a separate has_premium field, or restructure to make the presence meaningful (e.g., optional premium_tier string). This tests linking Protobuf mechanics with API design decisions.

## hook
How would you represent "user hasn't chosen a preference yet" vs "user explicitly chose false"?
