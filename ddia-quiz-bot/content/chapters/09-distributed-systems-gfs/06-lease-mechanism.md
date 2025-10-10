---
id: ch09-lease-mechanism
day: 1
tags: [gfs, lease, distributed-coordination]
related_stories: []
---

# GFS Lease Mechanism

## question
What is the typical duration of a primary lease in GFS, and why?

## options
- A) 1 second - for quick failure detection
- B) 60 seconds - to balance failure detection with coordination overhead
- C) 1 hour - to minimize master involvement
- D) 24 hours - for stable long-running operations

## answer
B

## explanation
GFS uses a 60-second lease for primary replicas. This duration balances quick failure detection (not waiting too long if a primary fails) with minimizing the overhead of lease renewals. The primary can request extensions if needed for ongoing operations.

## hook
How long should you wait before declaring a server "dead" in a distributed system?
