---
id: ch05-read-scaling
day: 18
tags: [replication, scalability, read-replicas, architecture]
related_stories: []
---

# Read Scaling with Replicas

## question
Your application has a 100:1 read-to-write ratio. You add read replicas but don't see expected performance improvement. What's the most likely cause?

## options
- A) The replicas are too far from users geographically
- B) The application is reading from the primary for consistency guarantees
- C) The replicas have insufficient CPU
- D) The replication protocol is inefficient

## answer
B

## explanation
If the application requires strong consistency or read-after-write guarantees, it may still be routing most reads to the primary despite having replicas available. This is common when developers default to strong consistency without considering if eventual consistency would suffice for certain read operations. The solution is to identify which reads can tolerate stale data and route only those to replicas.

## hook
How would you identify which queries can safely use read replicas?
