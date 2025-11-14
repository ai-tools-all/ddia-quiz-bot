---
id: virtual-network-tun-tap-difference
day: 11
tags: [networking, virtual-interfaces, vpn, tun, tap]
related_stories:
  - docker-networking
---

# TUN vs TAP Interfaces

## question
What is the key difference between TUN and TAP virtual network interfaces?

## options
- A) TUN is for Linux systems while TAP is for Windows systems
- B) TUN operates at Layer 3 (IP packets) while TAP operates at Layer 2 (Ethernet frames)
- C) TUN is encrypted while TAP is unencrypted
- D) TUN is faster but TAP is more secure

## answer
B

## explanation
TUN interfaces operate at Layer 3 and handle IP packets only, making them ideal for most VPN applications. TAP interfaces operate at Layer 2 and handle full Ethernet frames, allowing them to bridge entire network segments and carry non-IP protocols.

## hook
Why would most VPN software prefer TUN over TAP interfaces?
