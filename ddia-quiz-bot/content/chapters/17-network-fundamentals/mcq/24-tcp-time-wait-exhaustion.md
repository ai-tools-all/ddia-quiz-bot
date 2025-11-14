---
id: network-tcp-time-wait-exhaustion
day: 24
tags: [networking, tcp, time-wait, connection-pooling, practical]
related_stories:
  - network-basics
  - tcp-fundamentals
---

# TCP TIME_WAIT Exhaustion

## question
Your client application makes thousands of short-lived HTTP requests per second to an API server. You start getting "Cannot assign requested address" errors. Checking `netstat`, you see thousands of connections in TIME_WAIT state. What's happening?

## options
- A) The server is refusing new connections because it's overloaded
- B) Your client has exhausted available local ports because TIME_WAIT sockets still occupy them
- C) Your network interface has run out of IP addresses
- D) The server has too many connections in TIME_WAIT state

## answer
B

## explanation
TCP TIME_WAIT exhaustion is a common problem with high-rate short-lived connections. When a TCP connection closes, the side that initiates the close (active close) enters TIME_WAIT state for 2*MSL (typically 60-120 seconds) to handle delayed packets. During this time, the local port remains reserved and can't be reused for new connections to the same server:port. Your client has ~64K ephemeral ports (1024-65535). At thousands of requests per second, you quickly exhaust available ports because ports in TIME_WAIT can't be reused. After 64K connections within 60 seconds, you get "Cannot assign requested address". Solutions: (1) Use connection pooling/keep-alive to reuse connections instead of creating new ones, (2) Enable SO_REUSEADDR/SO_REUSEPORT, (3) Tune net.ipv4.tcp_tw_reuse kernel parameter, (4) Reduce TIME_WAIT duration (risky), (5) Use multiple client IPs to increase port space.

## hook
Why does TCP require TIME_WAIT state, and what problems could occur if it didn't exist?
