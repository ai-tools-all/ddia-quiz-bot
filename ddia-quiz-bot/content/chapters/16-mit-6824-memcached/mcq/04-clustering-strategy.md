---
id: memcached-intra-region-clustering
day: 1
tags: [memcached, clustering, scalability, network]
---

# Intra-Region Clustering Strategy

## question
Why does Facebook organize front-ends and memcached servers into independent clusters within each region?

## options
- A) To improve cache hit rates through data locality
- B) To limit N-squared TCP connection overhead and reduce incast congestion from hundreds of simultaneous responses
- C) To enable cross-cluster failover during server failures
- D) To partition users geographically within a region

## answer
B

## explanation
Within each region, front-ends and memcached servers group into independent clusters primarily to limit TCP connection overhead. Without clustering, every front-end would need connections to every memcached server (N-squared connections). Clusters also reduce incast congestion—when a front-end requests hundreds of items in parallel, responses arriving simultaneously can overwhelm network buffers. Clustering limits the blast radius of these parallel requests and keeps connection counts manageable.

## hook
What are the trade-offs between having many small clusters versus few large clusters in terms of cache efficiency and network overhead?
