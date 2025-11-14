---
id: virtual-network-namespace-isolation
day: 8
tags: [networking, namespaces, isolation, containers]
related_stories:
  - docker-networking
---

# Network Namespace Isolation

## question
Three containers are running on the same host. Container A listens on port 80, Container B listens on port 80, and Container C tries to connect to Container A on port 80. What happens?

## options
- A) Container B fails to start because port 80 is already bound by Container A
- B) Container C's connection fails because multiple containers cannot use the same port
- C) Both containers successfully listen on port 80 in their own namespaces; Container C needs Container A's IP address to connect
- D) The kernel automatically port-forwards all connections to the first container that bound port 80

## answer
C

## explanation
Network namespaces provide complete isolation of the network stack, including port spaces. Each container has its own port namespace, so multiple containers can bind to the same port number without conflict. To connect, Container C must know Container A's IP address (or use DNS on a custom bridge network).

## hook
What else besides ports are isolated in separate network namespaces?
