---
id: network-arp-cache
day: 7
tags: [networking, arp, cache, broadcast-domain]
related_stories:
  - network-basics
---

# ARP Cache and Broadcast Domain

## question
Why do operating systems maintain an ARP cache instead of sending an ARP request for every packet?

## options
- A) ARP broadcasts can only be sent once per minute due to protocol limitations
- B) ARP requests consume bandwidth and broadcast to all devices in the network, so caching reduces unnecessary traffic
- C) ARP caching is required by the IEEE 802.3 Ethernet standard
- D) Without caching, routers would drop ARP requests as invalid

## answer
B

## explanation
ARP requests are broadcast to all devices in the broadcast domain (FF:FF:FF:FF:FF:FF), which creates network overhead. If every packet required an ARP request, the network would be flooded with broadcasts. The ARP cache stores IP-to-MAC mappings for a period (typically ~20 minutes), so subsequent packets to the same destination can be sent directly without additional ARP broadcasts. This significantly reduces network traffic and improves performance.

## hook
What security vulnerability arises from devices trusting ARP cache entries?
