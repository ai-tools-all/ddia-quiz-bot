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
