---
id: virtual-network-loopback
day: 1
tags: [networking, virtual-interfaces, loopback]
related_stories:
  - docker-networking
---

# Loopback Interface Purpose

## question
What is the primary purpose of the loopback interface (lo) present on every system?

## options
- A) To provide a backup network interface when the physical NIC fails
- B) To allow a system to communicate with itself without using the physical network
- C) To bridge multiple network namespaces for container communication
- D) To enable wireless connectivity as an alternative to Ethernet

## answer
B

## explanation
The loopback interface (127.0.0.1) exists entirely in software and allows local processes to communicate with each other without packets leaving the system. This is useful for IPC, testing network code, and services that only need local access.

## hook
Why would a database listen on 127.0.0.1 instead of 0.0.0.0?
