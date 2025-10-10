---
id: ch04-forward-compatibility
day: 7
tags: [schema-evolution, forward-compatibility, compatibility]
related_stories: []
---

# Forward Compatibility

## question
What does forward compatibility mean in the context of schema evolution?

## options
- A) New code can read data written by old code
- B) Old code can read data written by new code
- C) The schema cannot be changed once deployed
- D) All services must be updated simultaneously

## answer
B

## explanation
Forward compatibility means old code can read data written by new code. This is crucial when you need to roll back a deployment - the old version of the code should still be able to process data created by the newer version. This typically requires that old code can ignore fields it doesn't recognize.

## hook
What happens when you roll back your service after deploying a schema change?
