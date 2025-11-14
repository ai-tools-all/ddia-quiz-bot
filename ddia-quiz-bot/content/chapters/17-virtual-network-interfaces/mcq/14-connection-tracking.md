---
id: virtual-network-connection-tracking
day: 14
tags: [networking, conntrack, nat, stateful-firewall]
related_stories:
  - docker-networking
---

# Connection Tracking (conntrack)

## question
Why does Linux maintain a connection tracking (conntrack) table for NAT operations?

## options
- A) To ensure return packets from external servers are correctly de-NATed and routed back to the original container
- B) To encrypt all outbound connections for security
- C) To balance load across multiple containers
- D) To cache DNS responses for faster lookups

## answer
A

## explanation
The conntrack table remembers active connections and their NAT transformations. When a container sends a packet that gets SNAT'd to the host's IP, conntrack records this mapping. When the reply comes back, conntrack uses this information to reverse the NAT (destination becomes container's IP) and route the packet correctly.

## hook
What happens to long-lived connections when the conntrack table becomes full?
