---
id: ch08-lease-expiry
day: 15
level: L5
tags: [leases, distributed-locking, timing]
related_stories: []
---

# Lease Expiry and Clock Issues

## question
A node holds a lease that expires at wall-clock time T. Its clock jumps forward due to NTP adjustment. What safety issue can occur?

## options
- A) The lease becomes permanent
- B) The node thinks the lease expired early and stops using it prematurely
- C) The node continues using the lease after it has actually expired
- D) Other nodes can't acquire new leases

## answer
C

## explanation
If a node's clock jumps backward (or runs slow), it might think its lease is still valid when it has actually expired according to real time. Other nodes, with correct clocks, might have already granted a new lease to someone else. This creates a situation where two nodes think they hold the same lease. This is why lease implementations should use monotonic clocks for duration checking, not wall-clock time, and why fencing tokens provide additional safety.

## hook
Can your lease expire while your clock says you still have time?
