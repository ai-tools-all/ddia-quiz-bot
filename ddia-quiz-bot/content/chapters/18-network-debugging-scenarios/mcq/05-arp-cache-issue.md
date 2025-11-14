---
id: debug-arp-cache-timeout
day: 5
tags: [debugging, arp, layer2, cache, intermittent]
related_stories:
  - network-debugging
---

# Debugging Intermittent Connectivity

## question
A server experiences intermittent connectivity issues where it can't reach another server on the same subnet for exactly 1-2 seconds, then connectivity resumes. Running `tcpdump` shows ARP requests being broadcast during the outage. What's happening?

## options
- A) The switch is periodically failing and recovering
- B) The ARP cache entry expired and needs to be refreshed
- C) The default gateway is becoming unreachable
- D) Network packets are being randomly dropped due to congestion

## answer
B

## explanation
ARP cache entries have a timeout (typically 60-300 seconds). When an entry expires, the system must send a new ARP request to discover the target's MAC address, which causes a brief delay (1-2 seconds) while waiting for the ARP reply. This creates intermittent connectivity issues with a predictable pattern. Switch failures would last longer, gateway issues don't affect same-subnet communication, and congestion would show in packet loss metrics, not regular ARP broadcasts.

## hook
How can you prevent these intermittent delays in production systems with static network topology?
