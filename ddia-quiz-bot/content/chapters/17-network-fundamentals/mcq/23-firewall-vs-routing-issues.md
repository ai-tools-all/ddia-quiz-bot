---
id: network-firewall-vs-routing-diagnosis
day: 23
tags: [networking, firewall, routing, troubleshooting, practical]
related_stories:
  - network-basics
  - troubleshooting
---

# Distinguishing Firewall vs Routing Issues

## question
You can't reach service at 10.0.50.100:8080 from your application server. Running `traceroute 10.0.50.100` shows the path reaches the destination host. Running `telnet 10.0.50.100 8080` times out. What does this tell you?

## options
- A) There's a routing problem - traffic can't reach 10.0.50.100
- B) The host 10.0.50.100 is down
- C) A firewall is blocking port 8080, but ICMP (traceroute) is allowed
- D) The service on port 8080 is not running

## answer
C

## explanation
This is a classic firewall blocking scenario. Traceroute uses ICMP (or UDP) and successfully reaches 10.0.50.100, proving: (1) routing works, (2) the host is up and responding to ICMP. However, telnet to port 8080 (TCP) times out, indicating TCP SYN packets to port 8080 are being dropped. This is characteristic of firewall behavior - allowing ICMP for operational purposes (ping, traceroute) but blocking specific TCP/UDP ports. A routing problem would affect all traffic to 10.0.50.100, including traceroute. If the host was down, traceroute would timeout at the last hop. If the service wasn't running, you'd get "Connection refused" (RST response) not a timeout. Timeout = packets dropped silently = firewall. To confirm: check firewall rules on intermediate devices and the destination host itself (iptables, security groups, etc.).

## hook
How can you determine which hop in the network path is dropping your packets?
