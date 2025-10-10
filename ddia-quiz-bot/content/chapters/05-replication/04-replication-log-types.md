---
id: ch05-replication-log-types
day: 4
tags: [replication, WAL, logical-log, statement-based]
related_stories: []
---

# Types of Replication Logs

## question
Which replication log approach is most portable across different storage engines and database versions?

## options
- A) Statement-based replication
- B) Write-ahead log (WAL) shipping
- C) Logical (row-based) log replication
- D) Trigger-based replication

## answer
C

## explanation
Logical log replication decouples the replication log from the storage engine internals. It logs changes at the row level (inserts, updates, deletes) in a format independent of the storage engine, making it easier to replicate between different database versions or even different database systems. WAL shipping is tightly coupled to the storage engine.

## hook
How does logical log replication enable zero-downtime database upgrades?
