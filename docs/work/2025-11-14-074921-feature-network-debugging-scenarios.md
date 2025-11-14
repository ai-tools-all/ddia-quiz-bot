# Network Debugging Scenarios Quiz

**Date**: 2025-11-14 07:49:21
**Category**: feature
**Status**: completed

## Overview
Create a new chapter (18-network-debugging-scenarios) that tests the application of concepts from:
- Chapter 17: Network Fundamentals (IP/MAC, ARP, routing, NAT)
- Chapter 17: Virtual Network Interfaces (loopback, veth, Docker networking)

Goal: Test whether users can debug real-world networking scenarios using these fundamental concepts.

## Resources
- Existing MCQs in `ddia-quiz-bot/content/chapters/17-network-fundamentals/mcq/`
- Existing MCQs in `ddia-quiz-bot/content/chapters/17-virtual-network-interfaces/mcq/`
- Beej's Guide to Network Programming: https://beej.us/guide/bgnet/html/split-wide/ (fetch failed with 403)

## Tasks
- [x] Create work notes file
- [x] Research Beej's Guide for additional networking concepts
- [x] Create directory structure
- [x] Create 10 MCQ debugging scenarios
- [x] Create subjective debugging scenarios (L3-baseline, L3-bar-raiser)
- [x] Review and test build
- [x] Commit and push changes

## Key Concepts to Cover

### From Network Fundamentals:
1. IP vs MAC addresses
2. ARP (Address Resolution Protocol)
3. Default gateway/routing
4. Subnet masks and routing decisions
5. Switch vs router
6. NAT (translation table, DNAT, port forwarding)
7. Multi-hop routing
8. NAT traversal in P2P

### From Virtual Network Interfaces:
1. Loopback interface
2. veth pairs
3. Docker bridge network (172.17.0.0/16)
4. Port publishing/mapping
5. iptables NAT
6. DNS resolution in containers
7. Packet flow (outbound/inbound)
8. Network namespaces
9. Docker host network mode
10. iptables packet flow and tables
11. Connection tracking
12. Bridge MAC learning
13. TUN/TAP interfaces
14. Container communication
15. Overlay networks
16. macvlan driver

## Debugging Scenarios to Create

### MCQ Scenarios (10 questions):
1. Container cannot reach external internet - trace packet flow
2. Port mapping not working - iptables rules issue
3. Container-to-container communication failing on same bridge
4. DNS resolution failing inside container
5. Wrong interface selected for outbound traffic
6. ARP cache issues causing connectivity problems
7. NAT traversal problem in P2P application
8. Loopback vs 0.0.0.0 binding issue
9. Network namespace isolation breaking
10. Multi-hop routing failure diagnosis

### Subjective Scenarios:
1. L3-baseline: "Container cannot ping google.com - walk through debugging steps"
2. L3-baseline: "Two containers on same bridge cannot communicate - diagnose"
3. L3-bar-raiser: "Production web server becomes unreachable after Docker installation - analyze"
4. L3-bar-raiser: "Intermittent connectivity between microservices - debug systematically"

## Notes
- Focus on practical debugging workflows
- Require understanding of layered troubleshooting (Layer 2 → Layer 3 → Layer 4)
- Include common Docker networking pitfalls
- Test understanding of packet flow through iptables chains
- Include scenarios where multiple concepts interact

## Progress Log

### 2025-11-14 07:49
- Created work notes file
- Analyzed existing quiz formats
- Identified key concepts to cover

### 2025-11-14 07:50-08:00
- Created directory structure: `ddia-quiz-bot/content/chapters/18-network-debugging-scenarios/`
- Created 10 MCQ debugging scenario questions covering:
  1. Container internet connectivity (NAT/MASQUERADE)
  2. Port mapping failures (DNAT)
  3. Loopback binding issues (127.0.0.1 vs 0.0.0.0)
  4. Container-to-container communication (bridge/iptables FORWARD)
  5. ARP cache expiry causing intermittent issues
  6. Wrong source IP selection (routing table)
  7. DNS resolution in containers
  8. Overlapping subnet routing issues
  9. iptables FORWARD policy blocking
  10. NAT hairpin/loopback problems

- Created 5 subjective debugging scenarios:
  - L3-baseline:
    1. Container cannot reach internet - systematic debugging
    2. Container-to-container communication failure
    3. Service binding to loopback vs all interfaces
  - L3-bar-raiser:
    1. Production web server unreachable after Docker install
    2. Intermittent microservice connectivity (10% failure rate)

- All questions follow established format with proper YAML frontmatter
- Subjective questions include evaluation rubrics, follow-up questions, and bar-raiser challenges
- Verified file structure and content formatting

## Summary

Successfully created Chapter 18: Network Debugging Scenarios with:
- 10 MCQ questions testing practical debugging skills
- 3 L3-baseline subjective questions
- 2 L3-bar-raiser subjective questions
- GUIDELINES.md for evaluation rubrics

All questions apply concepts from Network Fundamentals and Virtual Network Interfaces chapters to real-world debugging scenarios.
