---
id: ch09-replication-factor
day: 1
tags: [gfs, replication, fault-tolerance]
related_stories: []
---

# GFS Replication Factor

## question
How many replicas does GFS typically maintain for each chunk?

## options
- A) 1 (no replication)
- B) 2 replicas
- C) 3 replicas
- D) 5 replicas

## answer
C

## explanation
GFS typically maintains 3 replicas of each chunk on different chunk servers. This provides good fault tolerance against machine failures while keeping storage overhead reasonable. The replicas are spread across different machines and preferably across different racks for better availability.

## hook
If hard drives fail all the time, how many copies of your data do you need to sleep well at night?
