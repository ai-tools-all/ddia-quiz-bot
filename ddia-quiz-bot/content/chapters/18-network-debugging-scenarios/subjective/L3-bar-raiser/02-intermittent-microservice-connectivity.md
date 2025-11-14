---
id: debug-subjective-L3-bar-002
type: subjective
level: L3
category: bar-raiser
topic: network-debugging
subtopic: intermittent-connectivity-debugging
estimated_time: 7-10 minutes
---

# question_title - Debugging Intermittent Microservice Connectivity

## main_question - Core Question
"Two microservices running in Docker containers on the same host experience intermittent connectivity failures. Service A (172.17.0.2) can sometimes reach Service B (172.17.0.3), but about 10% of requests fail with connection timeouts. Both services work fine when tested individually. How would you approach debugging this intermittent issue?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Intermittent vs Persistent**: Recognize that intermittent issues require different debugging (logging, metrics, packet capture over time)
- **Layer-by-Layer Testing**: Test at different layers (ping for ICMP, telnet/nc for TCP, actual application protocol)
- **Resource Exhaustion**: Check for connection limits, port exhaustion, or resource contention
- **Timing and Patterns**: Look for patterns in when failures occur (time-based, load-based, specific operations)

### expected_keywords
- Primary keywords: intermittent, connection timeout, ephemeral ports, conntrack, resource limits
- Commands: tcpdump, netstat -ant, ss -s, conntrack -S, docker stats
- Concepts: connection pooling, port exhaustion, NAT connection tracking, race conditions

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Connection Tracking Table**: Linux conntrack table can fill up, causing random drops
- **ARP Cache Expiry**: Intermittent ARP resolution delays
- **Bridge MAC Learning**: Transient MAC table issues on bridge
- **Application-Level Issues**: Connection pool exhaustion, timeout configurations
- **Kernel Tuning**: net.ipv4.ip_local_port_range, nf_conntrack_max

### bonus_keywords
- Implementation: sysctl parameters, connection tracking, TCP states
- Tools: conntrack, sar, vmstat, container resource monitoring
- Performance: connection pooling, keep-alive, timeout tuning
- Debugging: correlation between failures and system events

## sample_excellent - Example Excellence
"Intermittent issues are challenging because they require capturing the failure as it happens. I'd start by setting up continuous monitoring with `tcpdump -i docker0 -w capture.pcap host 172.17.0.2 and host 172.17.0.3` to capture traffic over time, including failures. While that runs, I'd check several likely causes:

First, connection tracking table size - run `conntrack -S` or `sysctl net.netfilter.nf_conntrack_count` vs `net.netfilter.nf_conntrack_max`. If the table is full, new connections get dropped randomly. Docker creates many connections, and the default limit can be too low.

Second, ephemeral port exhaustion - check `netstat -ant | grep TIME_WAIT | wc -l` and review `sysctl net.ipv4.ip_local_port_range`. If Service A makes many short-lived connections to Service B, it might run out of local ports.

Third, resource limits - use `docker stats` to see if either container is hitting CPU/memory limits, causing intermittent slowdowns that trigger timeouts.

Fourth, application-level issues - connection pool configuration, timeout settings, or race conditions in the code.

I'd correlate timing of failures with system metrics. For the 10% failure rate, I'd analyze the tcpdump to see if it's TCP SYN packets being dropped (conntrack/firewall issue), TCP SYN-ACK not returning (service B overload), or connections establishing but then timing out (application issue)."

## sample_acceptable - Minimum Acceptable
"For intermittent issues, I'd use tcpdump to capture traffic and see what happens during failures. I'd check if the connection tracking table is full with `conntrack -S`, check for port exhaustion with `netstat`, and monitor resource usage with `docker stats`. The pattern of failures would help identify if it's a resource issue or configuration problem."

## common_mistakes - Watch Out For
- Trying to debug intermittent issues with one-time tests
- Not setting up monitoring/logging before testing
- Assuming it's a network issue without checking application layer
- Not correlating failures with resource usage or timing patterns
- Forgetting about connection tracking limits in Docker environments
- Not understanding the difference between connection failures at different layers

## follow_up_excellent - Depth Probe
**Question**: "Your tcpdump shows that during failures, Service A sends SYN packets that receive SYN-ACK responses, but then sends RST immediately. What does this tell you about where the problem is?"
- **Looking for**: Understanding TCP handshake, RST means application rejected connection (likely resource limit)
- **Red flags**: Blaming network when it's application-level issue

## follow_up_partial - Guided Probe
**Question**: "How would you determine if the issue is connection tracking table exhaustion vs port exhaustion?"
- **Hint embedded**: Leads to conntrack -S and netstat analysis
- **Concept testing**: Understanding different types of resource exhaustion

## follow_up_weak - Foundation Check
**Question**: "Why might an issue appear intermittent instead of happening all the time?"
- **Simplification**: Resource exhaustion, race conditions, timing dependencies
- **Building block**: Understanding intermittent vs persistent failures

## bar_raiser_question - L3→L4 Challenge
"After investigation, you find that connection failures spike exactly every 60 seconds and last for about 2 seconds. During these spikes, new connections fail but existing connections continue working. What could cause this periodic pattern, and what subsystem would you investigate?"

### bar_raiser_concepts
- **Periodic System Tasks**: Cron jobs, garbage collection, health checks, log rotation
- **ARP Cache Expiry**: Default ARP timeout is often 60 seconds
- **Connection Tracking GC**: Periodic garbage collection of connection tracking table
- **Application GC**: JVM garbage collection or similar runtime pauses
- **Correlation Analysis**: Connecting periodic failures to system-level events

### expected_depth
- Recognizing 60-second pattern suggests timeout or periodic task
- Understanding ARP cache behavior and its impact on connection establishment
- Ability to hypothesize about multiple possible causes and design tests to distinguish them
- Knowing how to monitor system-level events (dmesg, /var/log, tcpdump with timestamps)

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 3-4 min answer + 4-6 min discussion
- **Common next topics**: Performance debugging, connection tracking tuning, Docker production best practices
- **Difficulty indicators**: Requires understanding of intermittent failures, multiple subsystems, correlation analysis, and production debugging techniques
