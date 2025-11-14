# Network Fundamentals Advanced Quiz Generation

**Date**: 2025-11-14 07:57:08
**Category**: Feature
**Status**: In Progress

## Objective
Create 15 advanced network fundamentals MCQ questions:
- 5 questions testing depth of understanding of overall networking knowledge
- 10 practical scenario-based questions for common scenarios

## Existing Content Analysis
The `17-network-fundamentals/mcq/` directory already has 10 questions (01-10) covering:
1. IP vs MAC addresses - layer 2/3 addressing
2. ARP purpose - address resolution
3. Default gateway - routing basics
4. Subnet routing decision - routing logic
5. Switch vs router - forwarding mechanisms
6. NAT translation table - NAT basics
7. ARP cache broadcast - ARP details
8. Multi-hop routing - packet traversal
9. DNAT port forwarding - port mapping
10. NAT traversal P2P - peer-to-peer connectivity

## New Questions Plan (11-25)

### Depth-Testing Questions (5 questions)
These will test deeper understanding of networking concepts:
1. TCP flow control and congestion control interaction
2. DNS resolution process (recursive vs iterative)
3. TCP three-way handshake with SYN cookies (security)
4. BGP routing and AS path selection
5. VLAN tagging and trunk ports

### Practical Scenario Questions (10 questions)
Common scenarios developers/engineers encounter:
1. Troubleshooting connection timeout vs connection refused
2. Load balancer health check failures
3. TCP retransmissions and packet loss diagnosis
4. DNS caching issues and TTL
5. SSL/TLS handshake failures
6. MTU/MSS issues causing packet fragmentation
7. Debugging slow network performance
8. Firewall blocking vs routing issues
9. Connection pooling and TCP TIME_WAIT exhaustion
10. HTTP keep-alive and connection reuse

## File Format
Following existing format:
- Front matter: id, day, tags, related_stories
- Question section
- Options (A, B, C, D)
- Answer
- Explanation
- Hook (follow-up question)

## Files to Create
- 11-tcp-flow-control.md
- 12-dns-resolution-process.md
- 13-tcp-syn-cookies.md
- 14-bgp-as-path-selection.md
- 15-vlan-trunk-ports.md
- 16-connection-timeout-vs-refused.md
- 17-load-balancer-health-checks.md
- 18-tcp-retransmission-diagnosis.md
- 19-dns-caching-ttl.md
- 20-tls-handshake-failures.md
- 21-mtu-mss-fragmentation.md
- 22-network-performance-diagnosis.md
- 23-firewall-vs-routing-issues.md
- 24-tcp-time-wait-exhaustion.md
- 25-http-keepalive-reuse.md

## Tasks
- [x] Examine existing structure
- [ ] Create 15 new MCQ files
- [ ] Commit and push changes
