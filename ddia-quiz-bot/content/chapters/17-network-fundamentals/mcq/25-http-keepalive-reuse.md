---
id: network-http-keepalive-connection-reuse
day: 25
tags: [networking, http, keepalive, connection-pooling, performance, practical]
related_stories:
  - network-basics
  - http-fundamentals
---

# HTTP Keep-Alive and Connection Reuse

## question
Your application makes 100 sequential HTTP/1.1 requests to the same server using a library with keep-alive enabled. Each request takes 50ms. Disabling keep-alive, the same 100 requests take 15 seconds. What's causing the dramatic difference?

## options
- A) Keep-alive compresses the HTTP payload, reducing transmission time
- B) Keep-alive allows pipelining multiple requests simultaneously
- C) Without keep-alive, each request creates a new TCP connection with handshake overhead (~150ms)
- D) The server prioritizes keep-alive connections over new connections

## answer
C

## explanation
The math reveals the overhead: 100 requests × 150ms handshake = 15 seconds. Without keep-alive (HTTP/1.0 default), each request: (1) creates new TCP connection (SYN, SYN-ACK, ACK = ~1.5 RTTs), (2) performs TLS handshake if HTTPS (~2-3 RTTs), (3) sends HTTP request and receives response (~1 RTT), (4) closes connection (FIN handshake). With 50ms RTT, this totals ~150ms per request. With keep-alive (HTTP/1.1 default), the first request pays the handshake cost, but subsequent 99 requests reuse the same connection, paying only the request/response RTT (~50ms). Total time: 150ms + 99×50ms = ~5 seconds versus 100×150ms = 15 seconds. Option A is wrong - keep-alive doesn't compress. Option B is wrong - HTTP/1.1 pipelining exists but is rarely used and problematic. Keep-alive simply reuses the same TCP connection for multiple requests sequentially.

## hook
How does HTTP/2's multiplexing improve upon HTTP/1.1 keep-alive, and what problem does it solve?
