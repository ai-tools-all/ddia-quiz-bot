---
id: debug-server-same-subnet-unreachable
day: 2
tags: [networking, troubleshooting, arp, layer2, server]
related_stories:
  - network-fundamentals
  - arp-protocol
---

# Server Cannot Reach Another Server in Same Subnet

## question
Server A (192.168.1.10) cannot SSH to Server B (192.168.1.20) in the same /24 subnet, but Server A can reach the gateway (192.168.1.1) and other servers. Running `ip route` shows the correct route, and `iptables -L` shows ACCEPT policies. What should you check next?

## options
- A) Whether the default gateway has a route back to Server B
- B) Whether Server A has an ARP entry for Server B and if Server B is responding to ARP requests
- C) Whether NAT is properly configured between the two servers
- D) Whether the DNS server can resolve Server B's hostname

## answer
B

## explanation
Since both servers are in the same subnet (192.168.1.0/24), communication happens at Layer 2 using MAC addresses, not routing. Server A needs to resolve Server B's MAC address via ARP before sending frames. Common issues include: (1) incomplete ARP cache entry, (2) Server B not responding to ARP requests due to firewall rules blocking ARP or interface issues, (3) wrong subnet mask causing Server A to think B is on a different network. The gateway (A) is only involved for cross-subnet traffic. NAT (C) is not needed for same-subnet communication. DNS (D) only affects hostname resolution, not IP connectivity.

## hook
How would you use `tcpdump` or `arping` to diagnose whether ARP requests are being sent and if Server B is responding?
