---
id: network-mtu-mss-fragmentation-issues
day: 21
tags: [networking, mtu, mss, fragmentation, troubleshooting, practical]
related_stories:
  - network-basics
  - tcp-fundamentals
---

# MTU/MSS Issues and Packet Fragmentation

## question
Your application works fine over LAN but hangs when sending large payloads over a VPN (MTU 1400). Small requests work. TCP connections establish successfully. What's happening?

## options
- A) The VPN's bandwidth is too low for large payloads
- B) Large packets exceed the VPN's MTU, causing fragmentation or black-holing when "Don't Fragment" is set
- C) The VPN's firewall only allows packets smaller than 1400 bytes
- D) TCP's congestion control is limiting throughput over the VPN

## answer
B

## explanation
Classic MTU black hole scenario. Standard Ethernet MTU is 1500 bytes. VPN encapsulation (IPsec/OpenVPN headers) reduces effective MTU to ~1400. When your application sends large TCP segments, the IP layer tries to send packets larger than 1400 bytes. If the "Don't Fragment" (DF) bit is set (common for Path MTU Discovery), routers can't fragment and should send ICMP "fragmentation needed" messages back. However, many firewalls block ICMP, creating a black hole: large packets are silently dropped, no error returns, TCP retransmits the same large packet forever, connection hangs. Small packets (< 1400 bytes) work fine. Solution: lower TCP MSS (Maximum Segment Size) to account for reduced MTU. MSS = MTU - IP header (20) - TCP header (20). For MTU 1400: MSS should be 1360. This is often done via MSS clamping on VPN routers or by configuring endpoints.

## hook
How does Path MTU Discovery work, and why does ICMP blocking break it?
