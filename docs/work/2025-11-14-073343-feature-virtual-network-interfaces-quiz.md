# Virtual Network Interfaces Quiz Creation

**Date:** 2025-11-14
**Category:** feature
**Status:** completed

## Objective

Create a comprehensive quiz on virtual network interfaces and Docker networking covering:
- Virtual interface types (loopback, veth, bridge, TAP/TUN, etc.)
- Docker networking modes (bridge, host, overlay, macvlan)
- iptables and packet flow
- Network namespaces
- Container-to-container communication

## Tasks

- [x] Examine existing quiz structure and format
- [x] Create new chapter folder: `17-virtual-network-interfaces`
- [x] Create MCQ folder structure
- [x] Create 10 MCQ questions with increasing difficulty
  - Level 1 (Easy): 3 questions - Basic concepts
  - Level 2 (Medium): 4 questions - Understanding flows and interactions
  - Level 3 (Hard): 3 questions - Complex scenarios and troubleshooting
- [x] Review and validate all questions
- [x] Commit and push changes

## Quiz Structure

```
ddia-quiz-bot/content/chapters/17-virtual-network-interfaces/
  └── mcq/
      ├── 01-loopback-interface.md (Easy)
      ├── 02-veth-pairs.md (Easy)
      ├── 03-docker-bridge-network.md (Easy)
      ├── 04-docker-port-publishing.md (Medium)
      ├── 05-iptables-nat.md (Medium)
      ├── 06-container-dns-resolution.md (Medium)
      ├── 07-packet-flow-outbound.md (Medium)
      ├── 08-network-namespace-isolation.md (Hard)
      ├── 09-docker-host-network.md (Hard)
      └── 10-iptables-packet-flow.md (Hard)
```

## Notes

- Following existing format with YAML front matter
- Each question includes: question, options, answer, explanation, and hook
- Topics derived from the networking fundamentals content provided
- Progressive difficulty: basic concepts → flows/interactions → complex scenarios

## Second Batch (Questions 11-20)

**Additional Topics:**
- TAP/TUN interfaces for VPN
- Connection tracking (conntrack)
- Custom bridge networks
- iptables tables (filter vs nat)
- Container-to-container communication
- Docker network drivers comparison
- MACVLAN/IPVLAN
- Network troubleshooting scenarios
- Bridge interface internals
- Overlay networks basics

**Questions Created:**
- 11-tun-tap-difference.md (Easy)
- 12-bridge-mac-learning.md (Easy)
- 13-custom-bridge-benefits.md (Easy)
- 14-connection-tracking.md (Medium)
- 15-iptables-tables.md (Medium)
- 16-container-communication-same-bridge.md (Medium)
- 17-macvlan-driver.md (Medium)
- 18-overlay-network-purpose.md (Hard)
- 19-troubleshooting-container-network.md (Hard)
- 20-packet-flow-inbound-complete.md (Hard)

**Status:** Completed - 20 total MCQ questions created
