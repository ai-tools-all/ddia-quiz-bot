---
id: debug-docker-container-intermittent-connectivity
day: 4
tags: [networking, troubleshooting, docker, bridge, veth-pair, iptables]
related_stories:
  - docker-networking
  - virtual-interfaces
---

# Docker Container Intermittently Loses Connectivity

## question
A Docker container intermittently loses network connectivity. Sometimes it works fine, sometimes it cannot reach anything (not even other containers or the host). The container is using the default bridge network. What is a likely cause?

## options
- A) The docker0 bridge interface is flapping up and down
- B) The veth pair is being deleted and recreated, possibly due to container restarts or network namespace issues
- C) The host's default gateway is intermittently unreachable
- D) DNS resolution is timing out periodically

## answer
B

## explanation
Intermittent complete network loss (including inability to reach local containers) points to Layer 2/interface issues rather than Layer 3 routing or application-level problems. veth pairs connect containers to the bridge; if the veth pair is deleted, recreated, or has issues (often during container lifecycle events, network plugin bugs, or namespace problems), connectivity is completely lost until it's restored. The docker0 bridge (A) rarely flaps in stable systems. Host gateway issues (C) wouldn't affect container-to-container communication on the same bridge. DNS (D) wouldn't cause complete connectivity loss to IP addresses.

## hook
How would you monitor the veth interfaces and bridge to detect when the veth pair is being removed or recreated?
