---
id: virtual-network-port-publishing
day: 4
tags: [networking, docker, port-mapping, nat]
related_stories:
  - docker-networking
---

# Docker Port Publishing Mechanism

## question
When you run `docker run -p 8080:80 nginx`, what mechanism allows external clients to reach the container?

## options
- A) The container's eth0 interface is directly attached to the host's physical network
- B) Docker creates an iptables DNAT rule to forward host port 8080 to container port 80
- C) A proxy process running in the container forwards traffic from port 8080 to port 80
- D) Docker assigns the container a public IP address accessible from the Internet

## answer
B

## explanation
Docker uses iptables DNAT (Destination NAT) rules in the PREROUTING chain to forward traffic arriving at the host's port 8080 to the container's IP and port 80. This allows external access while keeping containers on a private network.

## hook
What happens to the return traffic from the container back to the external client?
