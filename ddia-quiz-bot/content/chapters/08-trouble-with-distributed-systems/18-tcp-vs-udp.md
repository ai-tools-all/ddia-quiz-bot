---
id: ch08-tcp-reliability
day: 18
tags: [tcp, udp, network-protocols]
related_stories: []
---

# TCP Reliability Guarantees

## question
What does TCP guarantee about message delivery in distributed systems?

## options
- A) Messages arrive within a bounded time
- B) Messages arrive in order and without duplication, if they arrive at all
- C) Messages always arrive eventually
- D) Messages are delivered exactly once

## answer
B

## explanation
TCP provides reliable, ordered, duplicate-free delivery IF the connection stays alive. It handles packet loss via retransmission, ensures ordering via sequence numbers, and prevents duplication. However, TCP cannot guarantee delivery if the network is partitioned, the remote host crashes, or the connection times out. When a TCP connection fails, you don't know if the last messages were delivered. This is why application-level acknowledgments are still necessary.

## hook
If TCP is "reliable", why do distributed systems still lose messages?
