---
id: farm-wal-placement
day: 5
tags: [wal, logging, memory, recovery]
---

# Write-Ahead Log Placement

## question
Where does FaRM store write-ahead logs (WAL), and what is the key advantage of this placement?

## options
- A) On disk for durability, accepting slower write performance
- B) In per-client message queues stored in server memory, enabling fast appends while maintaining recoverability
- C) In a centralized log server that all transactions write to
- D) In client-local storage only, with no server-side logging

## answer
B

## explanation
FaRM stores WAL entries in per-client message queues located in server non-volatile RAM. This provides fast memory-speed appends without disk I/O, while still enabling crash recovery since NVRAM contents survive power failures. The per-client structure creates N² communication channels across the cluster, avoiding centralized bottlenecks.

## hook
How does per-client log placement affect recovery when a server crashes?
