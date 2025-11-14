---
id: network-performance-slow-diagnosis
day: 22
tags: [networking, performance, latency, bandwidth, troubleshooting, practical]
related_stories:
  - network-basics
  - performance
---

# Diagnosing Slow Network Performance

## question
Your application's API calls to a remote service take 500ms each. Ping to the service shows 10ms RTT. Running multiple requests in parallel completes much faster (50ms per request when running 10 concurrent). What's the bottleneck?

## options
- A) Network bandwidth is saturated
- B) The remote service is processing requests slowly
- C) Your application is making requests sequentially without connection reuse, paying TCP handshake cost each time
- D) Packet loss is causing retransmissions

## answer
C

## explanation
The key clues: (1) ping RTT is 10ms (network is fast), (2) single request takes 500ms (50x higher than RTT), (3) concurrent requests are fast (50ms each). This pattern indicates sequential requests with connection overhead. If your app opens a new TCP connection for each request: TCP handshake (3 RTTs = 30ms) + TLS handshake (2-3 RTTs = 20-30ms) + actual request/response (1 RTT = 10ms) + connection teardown, plus any DNS lookups. This overhead dominates. When running concurrent requests, connections overlap, amortizing the overhead. If option B (slow service) were true, concurrent requests would also be slow. If option A (bandwidth), throughput would matter more than RTT. If option D (packet loss), you'd see retransmissions in tcpdump. Solution: use HTTP keep-alive and connection pooling to reuse TCP connections across requests. This eliminates repeated handshake overhead.

## hook
What are the performance differences between HTTP/1.1 keep-alive, HTTP/2 multiplexing, and HTTP/3 with QUIC?
