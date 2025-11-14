---
id: network-dnat-port-forwarding
day: 9
tags: [networking, nat, dnat, port-forwarding, server-hosting]
related_stories:
  - network-basics
---

# DNAT and Port Forwarding

## question
You're hosting a web server (192.168.1.20:80) behind your router (public IP: 100.72.5.12). You configured DNAT to forward external port 8080 to your server. An Internet user connects to 100.72.5.12:8080. What does your web server see as the incoming connection?

## options
- A) Source: Internet user's IP, Destination: 100.72.5.12:8080 (no translation occurs)
- B) Source: Internet user's IP, Destination: 192.168.1.20:80 (only destination translated)
- C) Source: 100.72.5.12, Destination: 192.168.1.20:80 (both addresses translated)
- D) Source: 192.168.1.1 (router's LAN IP), Destination: 192.168.1.20:80

## answer
B

## explanation
DNAT (Destination NAT) translates only the destination IP and port. The router receives a packet with destination 100.72.5.12:8080, applies the DNAT rule to translate it to 192.168.1.20:80, and forwards it to the internal server. The source IP remains the Internet user's actual IP address (unchanged), which allows the server to see where requests come from and respond directly. The router then performs reverse NAT on the response, translating the source back to 100.72.5.12:8080 before sending it to the Internet user.

## hook
Why is DNAT typically combined with SNAT in home router scenarios, and what's the difference?
