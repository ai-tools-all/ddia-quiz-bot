---
id: debug-port-mapping-failure
day: 2
tags: [debugging, docker, port-mapping, iptables, dnat]
related_stories:
  - network-debugging
---

# Debugging Port Mapping Failure

## question
You start a container with `-p 8080:80` but cannot access it from outside the host. Running `iptables -t nat -L DOCKER` shows no DNAT rule for port 8080. Which command would help diagnose why Docker didn't create the port mapping rule?

## options
- A) `docker inspect <container>` to check if the port mapping is recorded
- B) `netstat -tulpn | grep 8080` to see if something else is using port 8080
- C) `ip route show` to check the routing table
- D) `arp -a` to verify ARP entries for the container

## answer
A

## explanation
`docker inspect` will show the port mappings in the container's configuration. If it's missing from there, the container wasn't started with the port mapping at all. If it's present in docker inspect but not in iptables, then there's an issue with Docker's iptables rule creation. Option B could reveal port conflicts but wouldn't explain missing iptables rules. Options C and D are unrelated to port mapping configuration.

## hook
What's the difference between a port being mapped in the container config vs having the actual DNAT rule in iptables?
