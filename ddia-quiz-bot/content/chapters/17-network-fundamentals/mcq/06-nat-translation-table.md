---
id: network-nat-translation
day: 6
tags: [networking, nat, snat, port-forwarding]
related_stories:
  - network-basics
---

# NAT Translation Table

## question
Your router receives a packet from the Internet with destination IP 100.72.5.12:30001. The router's NAT table shows port 30001 maps to 192.168.1.10:51100. What will the router do?

## options
- A) Drop the packet because it's from the Internet and unsolicited
- B) Forward the packet to all devices on the LAN using broadcast
- C) Translate the destination to 192.168.1.10:51100 and forward to that internal device
- D) Respond with an ICMP error message indicating NAT is blocking the connection

## answer
C

## explanation
This is SNAT (Source NAT) in action. The NAT table entry exists because device 192.168.1.10:51100 previously initiated an outbound connection, and the router mapped it to external port 30001. When the response arrives at the router's public IP:30001, the router looks up the NAT table, translates the destination back to the internal address 192.168.1.10:51100, and forwards the packet to that device.

## hook
What happens when a NAT table entry expires before the response arrives?
