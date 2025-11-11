---
id: memcached-incast-congestion-detail
day: 1
tags: [memcached, incast, network, parallel-requests]
---

# Incast Congestion Details

## question
A Facebook front-end rendering a user's homepage needs to fetch 500 different data items from memcached in parallel. What network problem does this create, and how does clustering help?

## options
- A) Too many TCP connections overwhelm memcached servers; clustering reduces total connections
- B) Responses arriving simultaneously (incast) overwhelm the front-end's network buffer, causing packet loss; clustering limits the fan-out per cluster
- C) Requests are serialized causing high latency; clustering enables better parallelism
- D) Network bandwidth is saturated; clustering distributes load across more network links

## answer
B

## explanation
When a front-end issues 500 parallel requests, the responses can arrive nearly simultaneously, creating an incast traffic pattern that overwhelms the front-end's network receive buffer and switch buffers, causing packet loss and retransmissions. Clustering limits this problem by reducing the fan-out: instead of connecting to all memcached servers globally, a front-end connects only to servers in its cluster. While this doesn't eliminate incast entirely, it reduces the burst size and makes it more manageable. Additionally, clusters can be sized to keep the parallel request count within network capacity.

## hook
What application-level techniques could further mitigate incast beyond clustering?
