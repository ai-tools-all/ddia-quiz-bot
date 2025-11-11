---
id: memcached-cold-cluster-warmup
day: 1
tags: [memcached, cold-start, warmup, cache-miss]
---

# Cold Cluster Warmup

## question
When Facebook brings up a new memcached cluster, it starts "cold" with an empty cache. What is the primary challenge and strategy for warming it up?

## options
- A) The cluster has no data, so all requests miss and overwhelm the database; strategy is to gradually redirect traffic while allowing the cluster to populate naturally through misses
- B) Network connections aren't established; strategy is to pre-warm all TCP connections before serving traffic
- C) The cluster has no knowledge of hot keys; strategy is to pre-load popular keys from other clusters
- D) Replication lag prevents immediate use; strategy is to wait for full database replication before serving

## answer
A

## explanation
A cold cluster starts empty, so initially 100% of requests are cache misses that must query the database. If full production traffic hits the cold cluster immediately, the database would be overwhelmed. Facebook gradually increases the traffic percentage directed to the new cluster, allowing it to naturally populate its cache through normal miss-and-populate flows. As the cache hit rate increases, more traffic can be safely directed to it. This gradual warmup prevents database overload while building up the working set in the new cluster's memory.

## hook
What metrics would you monitor to determine the safe rate of traffic increase during cluster warmup?
