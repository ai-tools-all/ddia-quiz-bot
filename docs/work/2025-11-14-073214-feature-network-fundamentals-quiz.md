# Network Fundamentals Quiz Creation

**Date:** 2025-11-14 07:32:14
**Category:** feature
**Task:** Create MCQ quiz for network fundamentals guide

## Overview
Creating 10 MCQ questions based on the network fundamentals guide covering:
- IP/MAC addresses
- ARP
- Gateway & Routing
- Switches/Bridges/Routers
- NAT (SNAT, DNAT)

## Quiz Structure
- **Level 1 (Easy):** 3 questions - Basic concepts
- **Level 2 (Medium):** 4 questions - Application and analysis
- **Level 3 (Hard):** 3 questions - Complex scenarios and edge cases

## Tasks
- [x] Analyze existing MCQ format
- [x] Create chapter folder: `17-network-fundamentals`
- [x] Create 10 MCQ files following the established format
- [x] Review for accuracy and progressive difficulty

## Notes
- Following format from existing chapters (10-mit-6824-primary-backup, etc.)
- Each MCQ has: id, day, tags, question, options (A/B/C/D), answer, explanation, hook
- Questions should build on each other in difficulty

## Created Questions

### Level 1 (Easy - Days 1-3):
1. **IP vs MAC Address** - Why both addressing schemes are needed
2. **ARP Purpose** - Understanding address resolution protocol
3. **Default Gateway** - Role of gateway in routing

### Level 2 (Medium - Days 4-7):
4. **Subnet Routing Decision** - Applying subnet masks to determine routing
5. **Switch vs Router** - Layer 2 vs Layer 3 forwarding
6. **NAT Translation Table** - Understanding SNAT operation
7. **ARP Cache** - Why caching reduces broadcast traffic

### Level 3 (Hard - Days 8-10):
8. **Multi-Hop Routing** - MAC/IP behavior across router hops with TTL
9. **DNAT Port Forwarding** - Destination NAT for hosting services
10. **NAT Traversal P2P** - Challenges in peer-to-peer connectivity

## Quiz Summary
- Total Questions: 10 MCQs
- Difficulty Levels: 3 (progressive difficulty)
- Topics Covered: IP/MAC addressing, ARP, routing, switches/routers, NAT (SNAT/DNAT)
- Location: `ddia-quiz-bot/content/chapters/17-network-fundamentals/mcq/`

## Status
✅ Completed - All 10 MCQ questions created and verified
