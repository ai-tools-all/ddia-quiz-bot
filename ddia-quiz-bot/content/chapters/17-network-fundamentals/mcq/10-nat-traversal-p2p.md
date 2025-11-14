---
id: network-nat-traversal
day: 10
tags: [networking, nat, p2p, nat-traversal, connectivity]
related_stories:
  - network-basics
---

# NAT Traversal and P2P

## question
Two friends want to establish a peer-to-peer video call. Friend A is behind NAT (192.168.1.10 → 100.72.5.12) and Friend B is also behind NAT (10.0.0.15 → 203.0.113.5). Without any NAT traversal techniques, what problem will they encounter?

## options
- A) The connection will work normally because NAT is transparent to applications
- B) Neither can initiate a direct connection to the other because NAT blocks unsolicited incoming connections
- C) Only Friend A can call Friend B, but not vice versa
- D) The connection will work but with degraded video quality due to NAT overhead

## answer
B

## explanation
NAT creates a fundamental problem for peer-to-peer connectivity. When Friend A tries to connect to Friend B's public IP (203.0.113.5), Friend B's router has no NAT table entry for this unsolicited connection and will drop it. The same happens in reverse. Neither router has a mapping because neither internal host initiated an outbound connection first. This is why P2P applications use NAT traversal techniques like STUN (to discover public IPs), TURN (relay servers), and ICE (tries multiple connection methods). Without these, P2P connections between two NATed hosts fail, requiring a relay server or port forwarding configuration.

## hook
How do protocols like WebRTC establish P2P connections despite both peers being behind NAT?
