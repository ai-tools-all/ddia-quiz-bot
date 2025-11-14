---
id: debug-wrong-source-ip
day: 6
tags: [debugging, routing, source-ip, multiple-interfaces]
related_stories:
  - network-debugging
---

# Debugging Wrong Source IP Selection

## question
A server with two interfaces (eth0: 192.168.1.10, eth1: 10.0.0.10) sends packets to the internet with source IP 10.0.0.10 instead of the expected 192.168.1.10. What determines which source IP is used?

## options
- A) The interface that the packet exits through, determined by the routing table
- B) The interface with the lowest MAC address
- C) The interface that was configured first during boot
- D) Random selection between all available interfaces

## answer
A

## explanation
The routing table determines which interface a packet uses based on the destination IP. By default, the kernel selects the source IP from the interface the packet will exit through. If the default route is via eth1, packets will have source IP 10.0.0.10. To force a specific source IP, you need either policy-based routing or explicit source IP binding in the application using `bind()`.

## hook
How would you configure the system to always use 192.168.1.10 as the source IP for internet traffic?
