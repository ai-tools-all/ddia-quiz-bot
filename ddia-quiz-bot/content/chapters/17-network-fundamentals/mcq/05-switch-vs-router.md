---
id: network-switch-vs-router
day: 5
tags: [networking, switch, router, layer2, layer3]
related_stories:
  - network-basics
---

# Switch vs Router Operation

## question
What is the fundamental difference between how a switch and a router forward packets?

## options
- A) A switch forwards using MAC addresses (Layer 2) within a broadcast domain, while a router forwards using IP addresses (Layer 3) between different networks
- B) A switch forwards using IP addresses and a router forwards using MAC addresses
- C) A switch only connects wireless devices, while a router only connects wired devices
- D) Both operate identically, but routers have more Ethernet ports

## answer
A

## explanation
Switches operate at Layer 2 (Data Link) and forward frames based on MAC addresses within a single broadcast domain. They maintain a MAC address table learned from incoming traffic. Routers operate at Layer 3 (Network) and forward packets based on IP addresses between different networks, maintaining routing tables. Routers separate broadcast domains, while switches do not.

## hook
Why do routers block broadcast traffic while switches forward it?
