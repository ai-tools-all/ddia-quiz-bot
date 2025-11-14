---
id: virtual-network-packet-flow-inbound-complete
day: 20
tags: [networking, packet-flow, iptables, nat, advanced]
related_stories:
  - docker-networking
---

# Complete Inbound Packet Flow with NAT

## question
A client sends an HTTP request to your_host:8080 which is published to container port 80 at 172.17.0.2. The container responds. What happens to the packet's source and destination IPs throughout the round trip?

## options
- A) Request: Client→Host remains unchanged; Response: Container→Client remains unchanged
- B) Request: Client→Host gets DNAT to Client→172.17.0.2:80; Response: 172.17.0.2→Client gets SNAT to Host→Client
- C) Request: Client→Host gets SNAT to 172.17.0.2→Host; Response: Host→Client gets DNAT to Client→172.17.0.2
- D) Request and response both go through MASQUERADE without specific NAT transformations

## answer
B

## explanation
Inbound: (1) Client→Host:8080 arrives, (2) PREROUTING DNAT changes destination to 172.17.0.2:80, (3) packet reaches container. Outbound: (1) Container sends 172.17.0.2→Client, (2) conntrack remembers the original mapping, (3) source gets changed back to Host:8080, (4) Client receives Host:8080→Client. The conntrack table ensures the reverse transformation happens correctly.

## hook
What would happen if the conntrack table entry expired before the container's response arrived?
