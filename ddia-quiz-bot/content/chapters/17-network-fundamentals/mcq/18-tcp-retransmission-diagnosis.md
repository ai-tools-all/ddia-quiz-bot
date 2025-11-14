---
id: network-tcp-retransmission-packet-loss
day: 18
tags: [networking, tcp, retransmission, packet-loss, troubleshooting, practical]
related_stories:
  - network-basics
  - tcp-fundamentals
---

# TCP Retransmissions and Packet Loss

## question
You notice your application's TCP connections show 5% retransmission rate. Response times are slow but connections don't fail. What is most likely happening?

## options
- A) The network is congested or lossy, causing packets to drop and TCP to retransmit them
- B) Your application is sending malformed packets that the receiver rejects
- C) A firewall is randomly blocking packets
- D) The CPU on your server is overloaded

## answer
A

## explanation
A 5% retransmission rate indicates packet loss somewhere in the network path. TCP automatically detects missing acknowledgments and retransmits lost segments. This works transparently but hurts performance: each retransmit adds latency (typically waiting for RTO timeout or 3 duplicate ACKs for fast retransmit). Connections don't fail because TCP keeps retrying, but response times suffer. Common causes: network congestion (buffers full, packets dropped), faulty network hardware (bad cables, failing switches), WiFi interference, or misconfigured QoS policies. Option B would cause connection failures, not retransmits. Option C would likely show timeouts. Option D affects processing, not packet delivery. To diagnose: use tcpdump/Wireshark to see where packets are lost, check network metrics (queue drops, CRC errors), test different network paths.

## hook
How can you distinguish between packet loss due to congestion versus faulty hardware?
