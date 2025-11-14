---
id: debug-loopback-binding
day: 3
tags: [debugging, loopback, socket, binding, network-interfaces]
related_stories:
  - network-debugging
---

# Debugging Service Binding Issue

## question
A web service starts successfully and `netstat -tulpn` shows it listening on `127.0.0.1:3000`. However, external clients cannot connect, and even `curl http://<host-ip>:3000` from another machine times out. What's the problem?

## options
- A) The firewall is blocking port 3000
- B) The service is bound to the loopback interface only, not all interfaces
- C) The routing table doesn't have a route to the host IP
- D) ARP is not resolving the host's MAC address

## answer
B

## explanation
When a service binds to 127.0.0.1 (loopback), it only accepts connections from the same machine. External clients cannot reach it because packets destined for the loopback interface never leave the host's network stack. To accept external connections, the service should bind to 0.0.0.0 (all interfaces) or a specific external interface IP. Options A, C, and D would show different symptoms (firewall would reset connections, routing issues would affect all services, ARP issues would be local network only).

## hook
What happens when you bind to 0.0.0.0 vs a specific IP address?
