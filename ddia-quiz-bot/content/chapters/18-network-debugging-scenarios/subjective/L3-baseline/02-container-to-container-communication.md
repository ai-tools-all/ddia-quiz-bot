---
id: debug-subjective-L3-002
type: subjective
level: L3
category: baseline
topic: network-debugging
subtopic: container-to-container-communication
estimated_time: 5-7 minutes
---

# question_title - Debugging Container-to-Container Communication

## main_question - Core Question
"Two containers on the same Docker bridge network (docker0) cannot ping each other by IP, even though both can reach the internet. What would you check to diagnose this issue?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Same Subnet Check**: Verify both containers are on same subnet (172.17.0.0/16)
- **Bridge Connectivity**: Confirm both veth interfaces are attached to docker0 bridge
- **iptables FORWARD**: Check if FORWARD chain policy blocks inter-container traffic

### expected_keywords
- Primary keywords: bridge, veth pairs, FORWARD chain, same subnet, Layer 2
- Commands: docker network inspect, brctl show, iptables -L FORWARD, ip link
- Concepts: bridge forwarding, network namespace isolation

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **MAC Learning**: Bridge should learn MAC addresses from both containers
- **Bridge STP**: Spanning Tree Protocol could block ports
- **Network Policies**: Docker network-specific iptables rules
- **Packet Capture**: Use tcpdump on bridge to see if packets arrive

### bonus_keywords
- Implementation: bridge FDB (forwarding database), brctl showmacs
- Tools: tcpdump -i docker0, nsenter for namespace debugging
- Layer 2 concepts: MAC address table, broadcast domain

## sample_excellent - Example Excellence
"Since both containers can reach the internet, I know their individual connectivity is working. The issue must be with communication through the bridge. First, I'd run `docker network inspect bridge` to confirm both containers are actually on the same network. Then I'd check `brctl show docker0` to verify both veth interfaces are listed as ports on the bridge. Next, I'd check `iptables -L FORWARD -v` because even though containers are on the same bridge (Layer 2), Linux still routes packets through the FORWARD chain when they cross network namespaces. If the policy is DROP and there are no ACCEPT rules for docker0, that would block communication. I'd also use `brctl showmacs docker0` to verify the bridge has learned both containers' MAC addresses. Finally, I could use `tcpdump -i docker0` to see if packets from one container are reaching the bridge."

## sample_acceptable - Minimum Acceptable
"I'd first verify both containers are on the same Docker network with `docker network inspect`. Then check if both veth interfaces are attached to the docker0 bridge using `brctl show`. The most common issue would be iptables FORWARD chain blocking the traffic, so I'd check `iptables -L FORWARD`."

## common_mistakes - Watch Out For
- Not understanding that containers use network namespaces
- Thinking routing table is the issue (it's Layer 2 on same bridge)
- Confusing INPUT/OUTPUT chains with FORWARD chain
- Not knowing the difference between bridge forwarding and IP routing
- Forgetting that iptables can filter even same-subnet traffic

## follow_up_excellent - Depth Probe
**Question**: "If packets are reaching the bridge but not being forwarded to the destination container, what two Linux subsystems could be responsible, and how would you check each?"
- **Looking for**: iptables filtering (FORWARD chain), bridge filtering (ebtables), or bridge port blocking (STP)
- **Red flags**: Only mentioning one possibility, not understanding the difference

## follow_up_partial - Guided Probe
**Question**: "You mentioned the bridge. How does a bridge decide whether to forward a frame between two ports?"
- **Hint embedded**: Leads to MAC address learning and FDB
- **Concept testing**: Understanding of Layer 2 switching

## follow_up_weak - Foundation Check
**Question**: "When two containers are on the same bridge, do they need a router to communicate, or does the bridge handle it directly?"
- **Simplification**: Basic Layer 2 vs Layer 3 distinction
- **Building block**: Understanding bridge vs router functionality

## bar_raiser_question - L3→L4 Challenge
"The containers can ping each other by IP, but a web service in container A cannot be reached from container B using `curl http://172.17.0.2:8080`. The service is running and bound to 0.0.0.0:8080. What additional layer should you investigate?"

### bar_raiser_concepts
- Layer 4 (transport) vs Layer 3 (network) debugging
- iptables INPUT chain in container namespace
- Application binding (0.0.0.0 vs 127.0.0.1)
- Docker's per-container iptables rules

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 2-3 min answer + 3-4 min discussion
- **Common next topics**: Network namespace isolation, iptables chains, bridge learning
