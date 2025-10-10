---
id: ch09-single-master-bottleneck
day: 1
tags: [gfs, scalability, master-server, bottleneck]
related_stories: []
---

# Single Master Bottleneck

## question
How does GFS avoid the single master becoming a performance bottleneck?

## options
- A) By replicating the master across multiple machines
- B) By having clients cache metadata and read/write data directly from chunk servers
- C) By sharding the master's responsibilities
- D) By limiting the number of clients

## answer
B

## explanation
GFS minimizes master involvement by having clients cache metadata and interact directly with chunk servers for data operations. The master is only involved in metadata operations and coordination, not in the data path for reads and writes.

## hook
If you have one traffic cop for a whole city, how do you keep traffic flowing?
