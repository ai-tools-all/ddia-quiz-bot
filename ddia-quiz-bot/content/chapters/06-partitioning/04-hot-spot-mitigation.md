---
id: ch06-hot-spot-mitigation
day: 4
tags: [partitioning, hot-spots, load-balancing, performance]
related_stories: []
---

# Hot Spot Mitigation

## question
A social media application experiences hot spots when celebrity users post content (millions of followers reading the same data). Which technique would be most effective to mitigate this?

## options
- A) Switch from hash to range partitioning
- B) Add random prefix to keys and split reads across multiple partitions
- C) Increase the replication factor
- D) Use smaller partition sizes

## answer
B

## explanation
Adding a random prefix (like a two-digit number 00-99) to hot keys and splitting them across multiple partitions helps distribute the load. When reading, the application queries all prefixed versions and merges results. This technique, though it adds complexity, effectively spreads hot spot load across multiple nodes. Simply changing partitioning strategy or partition size won't help if one key is inherently popular.

## hook
How did Twitter solve the "Justin Bieber problem" in their timeline infrastructure?
