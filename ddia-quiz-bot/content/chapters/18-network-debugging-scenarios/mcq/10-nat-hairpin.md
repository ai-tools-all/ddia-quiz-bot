---
id: debug-nat-hairpin
day: 10
tags: [debugging, nat, hairpin, docker, dnat, loopback]
related_stories:
  - network-debugging
---

# Debugging NAT Hairpin Problem

## question
A container running a web service is published on host port 8080 (`-p 8080:80`). From inside the container, trying to access `http://<host-public-ip>:8080` fails, but external clients can access it fine. What's the issue?

## options
- A) The container needs a hairpin NAT rule to access itself via the host's public IP
- B) The container's firewall is blocking outbound connections to port 8080
- C) The host's public IP is not reachable from the docker0 bridge
- D) Port 8080 is not published correctly in the iptables DNAT rule

## answer
A

## explanation
This is called NAT hairpin or NAT loopback. When a container tries to reach its own service via the host's public IP, the packet goes out to the host, gets DNAT'd to the container's IP, but needs special routing to loop back. Without hairpin NAT rules, the return path fails. External clients work because they come from outside the host. Option B would affect all outbound traffic, option C is incorrect (the bridge can reach host IPs), and option D is wrong since external access works.

## hook
Why don't containers have this problem when accessing each other via the docker bridge IPs instead of the host's public IP?
