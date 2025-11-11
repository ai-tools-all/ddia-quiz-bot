---
id: cops-visibility-rules
day: 1
tags: [cops, replica, visibility, dependencies]
---

# Visibility Rules at Replicas

## question
How do remote shard servers in COPS handle incoming puts with dependencies?

## options
- A) They immediately make the put visible to maximize availability
- B) They wait until all dependencies are satisfied before making the put visible to local gets
- C) They reject the put if any dependency is missing and request a retry
- D) They use a voting mechanism with other replicas to decide when to apply the put

## answer
B

## explanation
Remote shard servers receiving a put do not immediately make it visible to local gets. Instead, they defer visibility until all dependencies in the attached list are satisfied (i.e., those versions exist and are visible locally). This ensures that any client reading the new version will also see (or have already seen) all causally prior writes, preserving the application's expected causal ordering. Gets return the highest visible version.

## hook
What happens to system availability if a dependency is delayed or lost due to network issues?
