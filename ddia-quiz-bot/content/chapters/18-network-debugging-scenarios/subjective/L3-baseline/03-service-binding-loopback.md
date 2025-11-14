---
id: debug-subjective-L3-003
type: subjective
level: L3
category: baseline
topic: network-debugging
subtopic: service-binding-interfaces
estimated_time: 5-7 minutes
---

# question_title - Debugging Service Binding and Interface Selection

## main_question - Core Question
"A developer reports that their web service starts successfully on port 3000, and they can access it via `curl http://localhost:3000` from the same machine, but colleagues cannot access it remotely. How would you diagnose and fix this issue?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Check Binding Address**: Use `netstat -tulpn` or `ss -tulpn` to see which interface the service is bound to
- **Loopback vs All Interfaces**: Understand difference between 127.0.0.1 (loopback only) and 0.0.0.0 (all interfaces)
- **Fix Configuration**: Service must bind to 0.0.0.0 or specific external IP to accept remote connections

### expected_keywords
- Primary keywords: bind, listening address, loopback, 127.0.0.1, 0.0.0.0, interface
- Commands: netstat, ss, lsof, curl
- Concepts: socket binding, network interfaces, localhost vs external IP

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Firewall Rules**: After fixing binding, check if firewall blocks external access
- **IPv4 vs IPv6**: Service might bind to ::1 (IPv6 loopback) vs 0.0.0.0 (IPv4 all)
- **Multiple Interfaces**: When to bind to specific IP vs 0.0.0.0
- **Security Implications**: Why services might intentionally bind to 127.0.0.1

### bonus_keywords
- Implementation: bind() system call, socket options, SO_REUSEADDR
- Tools: tcpdump to verify if packets arrive, iptables -L for firewall
- Security: principle of least privilege, attack surface reduction

## sample_excellent - Example Excellence
"First, I'd check what address the service is actually listening on using `netstat -tulpn | grep 3000`. If it shows `127.0.0.1:3000`, that's the problem - the service is bound to the loopback interface, which only accepts connections from the same machine. Loopback traffic never leaves the host's network stack, so remote clients can't reach it. The fix is to configure the service to bind to `0.0.0.0:3000` (all interfaces) or the specific external IP address. This is a common issue when developers test locally but forget that localhost/127.0.0.1 is not accessible remotely. After changing the binding, I'd also verify the firewall isn't blocking port 3000 with `iptables -L INPUT -n`. Some services bind to 127.0.0.1 intentionally for security (like databases that should only be accessed locally), so it's important to understand if remote access is actually desired."

## sample_acceptable - Minimum Acceptable
"I'd run `netstat -tulpn` to check what address the service is listening on. If it shows 127.0.0.1, the service is only listening on localhost and won't accept remote connections. The developer needs to change the configuration to bind to 0.0.0.0 instead."

## common_mistakes - Watch Out For
- Jumping to firewall without checking binding first
- Not understanding what 127.0.0.1 means
- Confusing localhost (loopback) with the machine's hostname
- Not knowing the difference between 0.0.0.0 and a specific IP
- Thinking this is a networking issue rather than application configuration

## follow_up_excellent - Depth Probe
**Question**: "If the service binds to a specific interface IP (like 192.168.1.100:3000) instead of 0.0.0.0:3000, what happens when clients try to connect via a different interface IP on the same machine?"
- **Looking for**: Understanding that binding is per-interface, not per-machine
- **Red flags**: Thinking one binding covers all IPs on the machine

## follow_up_partial - Guided Probe
**Question**: "What does 0.0.0.0 mean in the context of binding a socket?"
- **Hint embedded**: Leads to "all interfaces" understanding
- **Concept testing**: Special meaning of 0.0.0.0 in socket programming

## follow_up_weak - Foundation Check
**Question**: "Why does localhost work from the same machine but not from remote machines?"
- **Simplification**: Basic loopback concept
- **Building block**: Understanding localhost is local-only

## bar_raiser_question - L3→L4 Challenge
"The developer changes the binding to 0.0.0.0:3000 and remote access works. However, for security, they want ONLY machines on the 192.168.1.0/24 subnet to access the service. What are three different ways to achieve this, and what are the trade-offs of each?"

### bar_raiser_concepts
- Binding to specific interface IP (192.168.1.100)
- Firewall rules (iptables source filtering)
- Application-level filtering (reverse proxy, app middleware)
- Trade-offs: flexibility, performance, complexity, defense in depth

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 2-3 min answer + 3-4 min discussion
- **Common next topics**: Socket programming, firewall configuration, multi-homed hosts
