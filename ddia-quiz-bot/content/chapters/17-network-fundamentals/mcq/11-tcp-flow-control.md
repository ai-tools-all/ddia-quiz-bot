---
id: network-tcp-flow-congestion-control
day: 11
tags: [networking, tcp, flow-control, congestion-control, performance]
related_stories:
  - network-basics
  - tcp-fundamentals
---

# TCP Flow Control vs Congestion Control

## question
A TCP sender has 64 KB of data to send. The receiver's advertised window is 32 KB, and the sender's congestion window (cwnd) is 16 KB. How much data can the sender transmit?

## options
- A) 64 KB - the full amount available
- B) 32 KB - limited by receiver's advertised window
- C) 16 KB - limited by congestion window
- D) 48 KB - sum of both windows

## answer
C

## explanation
The effective TCP sending window is the minimum of the receiver's advertised window (flow control) and the sender's congestion window (congestion control). Flow control prevents overwhelming the receiver's buffer (32 KB available), while congestion control prevents overwhelming the network (16 KB safe to send). TCP always uses the more restrictive limit: min(32 KB, 16 KB) = 16 KB. This dual mechanism ensures both endpoint and network capacity are respected.

## hook
How does TCP's congestion window (cwnd) grow over time, and what triggers it to shrink?
