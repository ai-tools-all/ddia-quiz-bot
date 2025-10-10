---
id: ch09-gfs-purpose
day: 1
tags: [gfs, distributed-systems, storage, scalability]
related_stories: []
---

# Purpose of GFS

## question
What was the primary motivation for Google to build GFS (Google File System)?

## options
- A) To replace traditional relational databases
- B) To handle large-scale data storage for web crawling and indexing with commodity hardware
- C) To provide real-time transaction processing
- D) To eliminate the need for data replication

## answer
B

## explanation
GFS was designed to handle Google's massive data storage needs for web crawling, indexing, and other large-scale data processing tasks. It was built to work with commodity hardware that fails frequently, making fault tolerance a key design goal while keeping costs low.

## hook
How does Google store billions of web pages reliably on cheap, failure-prone hardware?
