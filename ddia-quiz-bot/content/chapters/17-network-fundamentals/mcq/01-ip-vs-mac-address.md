---
id: network-ip-vs-mac
day: 1
tags: [networking, ip-address, mac-address, osi-model]
related_stories:
  - network-basics
---

# IP Address vs MAC Address

## question
Why do networks need both IP addresses and MAC addresses instead of just one addressing scheme?

## options
- A) MAC addresses are used for routing across multiple networks, while IP addresses work only within a local network
- B) IP addresses enable global routing across networks, while MAC addresses handle local network communication within a broadcast domain
- C) IP addresses are permanent hardware identifiers, while MAC addresses are dynamically assigned by DHCP
- D) They serve identical purposes, and having both is simply a legacy design choice

## answer
B

## explanation
IP addresses (Layer 3) are logical identifiers that enable routing across multiple networks globally. MAC addresses (Layer 2) are hardware identifiers that work within a single broadcast domain (local network). When packets traverse routers, MAC addresses change at each hop while IP addresses remain constant (except with NAT).

## hook
What happens to MAC and IP addresses as a packet travels from your laptop to a remote server?
