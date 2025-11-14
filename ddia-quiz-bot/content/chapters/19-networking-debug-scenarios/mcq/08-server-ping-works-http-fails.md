---
id: debug-server-ping-works-http-fails
day: 8
tags: [networking, troubleshooting, iptables, firewall, connection-tracking]
related_stories:
  - iptables-basics
  - connection-tracking
---

# Server Responds to Ping But Not HTTP

## question
You can ping a web server at 203.0.113.10, but HTTP requests on port 80 time out. The web server process is running and listening on 0.0.0.0:80. What is the most likely cause?

## options
- A) The server's default gateway is misconfigured
- B) The server's iptables firewall has rules blocking incoming TCP port 80 or connection tracking is broken
- C) ICMP and HTTP use different routing tables
- D) The server's MAC address is not in your ARP cache

## answer
B

## explanation
ICMP (ping) and TCP (HTTP) are different protocols. The server likely has iptables rules that allow ICMP but block TCP port 80, either with explicit DROP/REJECT rules in the INPUT chain or by having a default DROP policy without rules allowing port 80. Connection tracking issues can also cause TCP to fail while ICMP works. Default gateway (A) affects outbound traffic, not inbound connections. Linux doesn't use different routing tables for protocols by default (C). If MAC address (D) were an issue, ping wouldn't work either.

## hook
What iptables commands would you use to check the INPUT chain rules and verify if TCP port 80 is being blocked?
