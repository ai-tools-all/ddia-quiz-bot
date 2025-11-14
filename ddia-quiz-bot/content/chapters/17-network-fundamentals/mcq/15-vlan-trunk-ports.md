---
id: network-vlan-trunk-tagging
day: 15
tags: [networking, vlan, switching, trunking, 802.1q]
related_stories:
  - network-basics
  - switching
---

# VLAN Trunk Ports and 802.1Q Tagging

## question
You connect two switches with a trunk port configured to carry VLANs 10, 20, and 30. A frame from VLAN 20 traverses the trunk. What happens to the frame?

## options
- A) The frame is dropped because trunk ports don't carry data frames, only control traffic
- B) The frame is tagged with VLAN ID 20 using 802.1Q, allowing the receiving switch to forward it to the correct VLAN
- C) The frame's source MAC address is rewritten to include the VLAN number
- D) The frame is broadcast to all VLANs on the receiving switch

## answer
B

## explanation
Trunk ports use 802.1Q tagging to carry multiple VLANs across a single physical link. When a frame from VLAN 20 enters the trunk, the switch inserts a 4-byte 802.1Q tag between the source MAC and EtherType fields. This tag contains the VLAN ID (20) and priority information. The receiving switch reads this tag, strips it off, and forwards the frame only to ports belonging to VLAN 20. This mechanism allows a single cable to carry traffic for multiple broadcast domains. Access ports (non-trunk) don't add tags - they're for end devices. Trunk ports are for switch-to-switch or switch-to-router connections where multiple VLANs must traverse a single link.

## hook
What is the "native VLAN" on a trunk port, and why does it remain untagged?
