---
id: ch05-replication-purpose
day: 1
tags: [replication, distributed-systems, availability, scalability]
related_stories: []
---

# Purpose of Replication

## question
What is the primary purpose of data replication in distributed systems?

## options
- A) To reduce storage costs by compression
- B) To keep data close to users, tolerate failures, and scale read throughput
- C) To simplify database schema design
- D) To eliminate the need for backups

## answer
B

## explanation
Replication serves three main purposes: (1) keeping data geographically close to users to reduce latency, (2) allowing the system to continue working even if some parts fail (fault tolerance), and (3) scaling out the number of machines that can serve read queries to increase read throughput.

## hook
How does replication help when your database server is on another continent?
