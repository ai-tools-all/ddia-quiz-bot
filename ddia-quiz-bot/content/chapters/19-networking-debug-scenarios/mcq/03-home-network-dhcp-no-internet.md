---
id: debug-home-network-dhcp-no-internet
day: 3
tags: [networking, troubleshooting, dhcp, default-gateway, nat, home-network]
related_stories:
  - network-fundamentals
  - nat-basics
---

# Home Network Device Gets IP But No Internet

## question
A laptop on a home network successfully obtains an IP address (192.168.0.15/24) via DHCP and can ping the router (192.168.0.1), but cannot reach any external websites. The laptop's routing table shows 192.168.0.1 as the default gateway. What is the most likely cause?

## options
- A) The DHCP server didn't provide DNS server information
- B) The router's WAN interface is down or the ISP connection is broken
- C) The laptop's subnet mask is incorrect
- D) The laptop needs a static route to reach external networks

## answer
B

## explanation
Since the laptop can ping the router, the local network is functioning correctly (DHCP, IP configuration, default gateway, local routing). The issue is that traffic cannot go beyond the router to external networks. This indicates the router's WAN interface is down, the ISP connection is broken, or the router itself has issues routing to the Internet. DNS issues (A) would only affect domain name resolution, not connectivity to IPs like 8.8.8.8. The subnet mask (C) is correct if local pinging works. A static route (D) is unnecessary when the default gateway is properly configured.

## hook
What commands would you run to differentiate between a DNS issue versus a complete Internet connectivity failure?
