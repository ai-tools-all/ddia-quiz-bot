---
id: ch04-backward-compatibility
day: 8
tags: [schema-evolution, backward-compatibility, compatibility]
related_stories: []
---

# Backward Compatibility

## question
What does backward compatibility mean in the context of schema evolution?

## options
- A) Old code can read data written by new code
- B) New code can read data written by old code
- C) Data must be migrated before deploying new code
- D) The schema must include a version number

## answer
B

## explanation
Backward compatibility means new code can read data written by old code. This is essential for gradual deployments where old and new versions of the code run simultaneously, and new code must process data created by the still-running old code. Generally easier to achieve than forward compatibility.

## hook
How do you deploy schema changes without breaking currently running services?
