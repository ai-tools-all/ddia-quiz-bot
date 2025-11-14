---
id: debug-home-router-port-forwarding-fails
day: 10
tags: [networking, troubleshooting, nat, dnat, port-forwarding, home-network]
related_stories:
  - nat-basics
  - dnat-port-forwarding
---

# Home Router Port Forwarding Not Working

## question
You configured port forwarding on your home router to forward external port 8080 to internal server 192.168.1.100:80, but external clients cannot reach the service. The internal server is reachable from other local devices at 192.168.1.100:80. What should you verify?

## options
- A) Whether the router's WAN interface has a public IP address or is behind carrier-grade NAT (CGNAT)
- B) Whether the internal server has the correct subnet mask
- C) Whether the router's DHCP server is functioning
- D) Whether DNS is correctly configured on the internal server

## answer
A

## explanation
For port forwarding to work from external clients, the router must have a publicly routable WAN IP address. Many ISPs now use CGNAT (Carrier-Grade NAT), where multiple customers share a public IP behind the ISP's NAT. In this case, external clients connect to the ISP's NAT device, not your router, so port forwarding rules on your router are ineffective. You need to check if your WAN IP (shown on router) matches your public IP (from sites like ifconfig.me). If different, you're behind CGNAT. Subnet mask (B) doesn't affect inbound routing. DHCP (C) and DNS (D) don't affect inbound port forwarding.

## hook
How would you determine if you're behind CGNAT, and what alternatives exist for exposing internal services when port forwarding is unavailable?
