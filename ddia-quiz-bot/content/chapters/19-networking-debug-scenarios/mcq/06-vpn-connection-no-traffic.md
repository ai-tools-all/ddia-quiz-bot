---
id: debug-vpn-connection-no-traffic
day: 6
tags: [networking, troubleshooting, vpn, routing, ip-forwarding]
related_stories:
  - network-fundamentals
  - routing-basics
---

# VPN Connection Established But No Traffic Flows

## question
A VPN connection successfully establishes (client shows connected status), but you cannot reach any resources on the remote network (10.0.0.0/8). The VPN created a tun0 interface with IP 10.8.0.2. Running `ip route` shows a route: "10.0.0.0/8 via 10.8.0.1 dev tun0". What should you investigate?

## options
- A) Whether the VPN server has IP forwarding enabled and routes traffic between the VPN subnet and the internal network
- B) Whether your local machine has the correct DNS servers configured
- C) Whether the VPN server certificate is valid
- D) Whether ARP is functioning correctly on the tun0 interface

## answer
A

## explanation
The VPN tunnel is established (authentication worked, tun0 interface exists, routes are in place on the client). The issue is that traffic is not flowing through to the remote network. The VPN server must: (1) have IP forwarding enabled (sysctl net.ipv4.ip_forward=1) to route packets between networks, (2) have proper routes to reach the internal 10.0.0.0/8 network, and (3) potentially have iptables rules to allow forwarding between the VPN subnet (10.8.0.0/24) and internal network. DNS (B) doesn't affect IP connectivity. Certificate validation (C) happens during connection establishment. TUN/TAP devices (D) work at Layer 3 and don't use ARP.

## hook
What commands would you use to verify IP forwarding is enabled on the VPN server and check if packets are being forwarded between interfaces?
