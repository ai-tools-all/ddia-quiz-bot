---
id: network-arp-purpose
day: 2
tags: [networking, arp, address-resolution, layer2]
related_stories:
  - network-basics
---

# ARP (Address Resolution Protocol)

## question
What is the primary purpose of ARP (Address Resolution Protocol)?

## options
- A) To translate domain names (like google.com) into IP addresses
- B) To discover the MAC address associated with a known IP address on the local network
- C) To assign IP addresses dynamically to devices joining the network
- D) To encrypt network traffic between two devices

## answer
B

## explanation
ARP is used to discover the MAC address corresponding to a known IP address on the local network. When a device wants to send data to an IP address in its local subnet, it broadcasts an ARP request asking "Who has IP X.X.X.X?" and the device with that IP replies with its MAC address. DNS handles domain-to-IP translation, DHCP handles IP assignment, and encryption is handled by protocols like TLS/IPsec.

## hook
How does a device avoid sending ARP requests for every single packet it sends?
