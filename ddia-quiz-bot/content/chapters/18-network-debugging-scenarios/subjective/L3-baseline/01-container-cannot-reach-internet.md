---
id: debug-subjective-L3-001
type: subjective
level: L3
category: baseline
topic: network-debugging
subtopic: container-internet-connectivity
estimated_time: 5-7 minutes
---

# question_title - Debugging Container Internet Connectivity

## main_question - Core Question
"A Docker container on the default bridge network cannot reach the internet (ping 8.8.8.8 fails), but the host can reach the internet fine. Walk me through the systematic steps you would take to debug this issue, explaining what each step tells you."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Check Default Gateway**: Verify container has correct gateway (172.17.0.1) using `ip route` or `route -n`
- **Verify NAT Rules**: Check iptables POSTROUTING for MASQUERADE rule on 172.17.0.0/16 subnet
- **Test Layer by Layer**: Test connectivity at each hop (container→bridge→host→internet)

### expected_keywords
- Primary keywords: default gateway, NAT, MASQUERADE, iptables, routing, POSTROUTING
- Commands: ip route, iptables -t nat -L, ping, traceroute, tcpdump
- Network layers: Layer 2 (bridge), Layer 3 (routing), Layer 4 (transport)

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **IP Forwarding**: Check if net.ipv4.ip_forward=1 on host
- **FORWARD Chain**: Verify iptables FORWARD chain allows traffic from docker0
- **Packet Tracing**: Use tcpdump to see where packets are dropped
- **DNS vs Connectivity**: Distinguish between DNS resolution and actual connectivity issues

### bonus_keywords
- Implementation: sysctl, /proc/sys/net/ipv4/ip_forward, veth pairs, docker0 bridge
- Tools: tcpdump, iptables -v (verbose), conntrack
- Systematic approach: test with IP first, then DNS; check outbound then inbound

## sample_excellent - Example Excellence
"I would approach this systematically, layer by layer. First, I'd exec into the container and check `ip route` to verify it has a default gateway (should be 172.17.0.1, the docker0 bridge IP). Then I'd ping the gateway to confirm Layer 2/3 connectivity to the host. Next, from the host, I'd check `iptables -t nat -L POSTROUTING` to verify there's a MASQUERADE rule for the 172.17.0.0/16 subnet - this is critical because container IPs are private and need NAT to reach the internet. I'd also verify `cat /proc/sys/net/ipv4/ip_forward` returns 1, enabling packet forwarding. Finally, I'd check `iptables -L FORWARD` to ensure container traffic isn't being blocked. If packets are leaving but not returning, I might use `tcpdump -i docker0` to trace the packet flow."

## sample_acceptable - Minimum Acceptable
"First, check if the container has a default gateway set with `ip route`. Then verify the NAT rules exist with `iptables -t nat -L` - Docker needs MASQUERADE to translate the container's private IP. Also check that IP forwarding is enabled on the host. If those look good, check if iptables is blocking the traffic."

## common_mistakes - Watch Out For
- Jumping straight to DNS without checking basic connectivity
- Not understanding the role of NAT/MASQUERADE for private IPs
- Confusing FORWARD chain with NAT tables
- Not taking a systematic layer-by-layer approach
- Forgetting to check if IP forwarding is enabled

## follow_up_excellent - Depth Probe
**Question**: "If you see packets leaving the container with tcpdump on docker0, but they never reach the external network, what specific iptables chain would be responsible, and how would you verify it?"
- **Looking for**: Understanding of FORWARD chain for packet forwarding, POSTROUTING for NAT
- **Red flags**: Confusion between INPUT/OUTPUT and FORWARD chains

## follow_up_partial - Guided Probe
**Question**: "You mentioned checking the gateway. What command would show you the gateway, and what IP should it be for a container on the default bridge?"
- **Hint embedded**: Guides toward `ip route` and 172.17.0.1
- **Concept testing**: Understanding of default Docker bridge configuration

## follow_up_weak - Foundation Check
**Question**: "Why do you think the container needs a gateway to reach the internet? What would happen if it tried to send packets directly?"
- **Simplification**: Basic routing concept
- **Building block**: Understanding that containers are in separate network namespace

## bar_raiser_question - L3→L4 Challenge
"The container can reach the internet by IP (ping 8.8.8.8 works) but cannot resolve domain names (ping google.com fails). The host's /etc/resolv.conf has nameserver 127.0.0.53. How does this explain the problem, and how would you fix it?"

### bar_raiser_concepts
- systemd-resolved and 127.0.0.53
- Docker's /etc/resolv.conf copying behavior
- Loopback addresses not reachable from containers
- DNS vs IP connectivity distinction

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 2-3 min answer + 3-4 min discussion
- **Common next topics**: iptables packet flow, network namespaces, Docker networking modes
