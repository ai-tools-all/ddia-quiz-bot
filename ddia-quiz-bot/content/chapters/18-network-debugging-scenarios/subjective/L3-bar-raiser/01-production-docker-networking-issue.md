---
id: debug-subjective-L3-bar-001
type: subjective
level: L3
category: bar-raiser
topic: network-debugging
subtopic: production-networking-troubleshooting
estimated_time: 7-10 minutes
---

# question_title - Production Web Server Becomes Unreachable After Docker Installation

## main_question - Core Question
"A production web server running on port 80 becomes unreachable from external clients immediately after installing Docker. The web service is still running, and you can access it via `curl http://localhost:80` from the server itself. What could cause this, and how would you systematically diagnose and resolve it?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **iptables FORWARD Policy**: Docker may have changed the default FORWARD policy to DROP
- **Port Conflicts**: Docker's docker-proxy might bind to 0.0.0.0:80 if a container uses that port
- **IP Forwarding Changes**: Docker enables IP forwarding which might interact with existing firewall rules
- **Systematic Diagnosis**: Check service binding, firewall rules, and Docker's impact on networking

### expected_keywords
- Primary keywords: iptables, FORWARD policy, docker-proxy, port binding, netstat
- Commands: iptables -L, iptables -t nat -L, netstat -tulpn, ss -tulpn, docker ps
- Concepts: firewall policies, port conflicts, packet forwarding

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Docker's iptables Modifications**: Docker adds rules to FORWARD, POSTROUTING, and creates DOCKER chain
- **docker-proxy Process**: Docker's userspace proxy for port publishing
- **Connection Tracking**: conntrack entries might be affected
- **Interface Binding**: Service might have been listening on a specific interface that Docker affected
- **Packet Flow Tracing**: Using tcpdump to see where packets are dropped

### bonus_keywords
- Implementation: iptables-save, iptables-restore, systemd service conflicts
- Tools: conntrack -L, tcpdump, iptables -v -L -n --line-numbers
- Advanced: Docker's --iptables flag, docker network ls

## sample_excellent - Example Excellence
"This is likely an iptables FORWARD policy issue. Docker modifies iptables extensively when installed, and one common change is setting the FORWARD chain policy to DROP (for security) and then adding specific ACCEPT rules for Docker networks. If the web server was relying on packet forwarding (uncommon for a simple web server, but possible if behind a NAT or load balancer configuration), this could break connectivity.

I'd diagnose systematically: First, `netstat -tulpn | grep :80` to confirm the web server is still bound to 0.0.0.0:80 or the external IP. Second, check if a docker-proxy process grabbed the port with `ps aux | grep docker-proxy` and `docker ps` to see if any container is using port 80. Third, run `iptables -L -v -n` and check the FORWARD chain policy - if it's DROP and there's no rule allowing established connections, that could block return traffic in certain network configurations. Fourth, check `iptables -t nat -L PREROUTING` to see if Docker added any DNAT rules affecting port 80.

The fix depends on the root cause: if it's iptables FORWARD, add a rule to allow established connections or web server traffic; if it's a port conflict, stop the conflicting container; if Docker's iptables integration is causing issues, could restart Docker with --iptables=false (though this has security implications)."

## sample_acceptable - Minimum Acceptable
"Docker heavily modifies iptables when installed, so I'd check `iptables -L` to see if it changed the FORWARD chain policy. I'd also check with `netstat` if something else is now bound to port 80, like a Docker container or docker-proxy. If iptables is the issue, I'd add rules to allow traffic to port 80."

## common_mistakes - Watch Out For
- Not considering iptables changes Docker makes
- Forgetting about docker-proxy port conflicts
- Not checking both filter table and nat table
- Assuming the web server crashed without verifying
- Not understanding when FORWARD chain matters vs INPUT chain
- Jumping to conclusions without systematic diagnosis

## follow_up_excellent - Depth Probe
**Question**: "Under what circumstances would a web server be affected by the FORWARD chain policy? Isn't INPUT chain what matters for incoming connections to local services?"
- **Looking for**: Understanding of when traffic goes through FORWARD vs INPUT (bridged networks, VMs, certain NAT setups)
- **Red flags**: Not knowing the difference between FORWARD and INPUT chains

## follow_up_partial - Guided Probe
**Question**: "You see the FORWARD policy is DROP. What iptables command would you use to temporarily test if this is the cause?"
- **Hint embedded**: Leads to iptables -P FORWARD ACCEPT or adding specific ACCEPT rule
- **Concept testing**: Understanding of policy vs rules, and how to test safely

## follow_up_weak - Foundation Check
**Question**: "What's the difference between a service being stopped vs being unreachable due to firewall rules?"
- **Simplification**: Basic troubleshooting distinction
- **Building block**: Layered debugging approach

## bar_raiser_question - L3→L4 Challenge
"After investigation, you find that Docker set FORWARD policy to DROP, but the web server is on the host itself (not in a container), so INPUT chain should handle it. Yet external clients still can't connect. You discover the server is behind a NAT router that forwards port 80 to this machine. How does this explain why FORWARD policy matters, and what's happening to the packets?"

### bar_raiser_concepts
- **NAT and Packet Flow**: Packets from internet → NAT router (DNAT) → web server involves forwarding on NAT router
- **Multi-hop Packet Path**: Understanding where packets traverse
- **FORWARD Chain Applicability**: FORWARD applies to packets being routed through the system
- **DNAT Impact**: Destination NAT changes packet destination, affecting which chain processes it

### expected_depth
- Understanding that if the server itself is doing NAT/forwarding for other reasons, FORWARD applies
- Recognizing packet path complexity in real production environments
- Ability to trace packet flow through multiple hops and iptables chains

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 3-4 min answer + 4-6 min discussion
- **Common next topics**: iptables advanced rules, Docker networking deep dive, production troubleshooting strategies
- **Difficulty indicators**: Requires understanding of Docker's impact on host networking, iptables chain selection, and systematic debugging under pressure
