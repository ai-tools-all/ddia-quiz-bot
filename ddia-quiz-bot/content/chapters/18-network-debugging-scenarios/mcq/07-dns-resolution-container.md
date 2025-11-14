---
id: debug-dns-container
day: 7
tags: [debugging, dns, docker, resolv.conf]
related_stories:
  - network-debugging
---

# Debugging DNS in Containers

## question
A container can ping 8.8.8.8 successfully but `ping google.com` fails with "Temporary failure in name resolution". What should you check first?

## options
- A) The container's /etc/resolv.conf to verify DNS server configuration
- B) The iptables NAT rules for DNS traffic
- C) The container's default gateway setting
- D) Whether port 53 is published on the container

## answer
A

## explanation
Since the container can reach 8.8.8.8 by IP, network connectivity and routing are working. The DNS resolution failure indicates a problem with DNS configuration. Docker copies the host's /etc/resolv.conf to the container, but if it's misconfigured or empty, DNS won't work. Check if nameserver entries exist. NAT rules aren't specific to DNS (option B), the gateway works (ping succeeds), and containers don't need to publish port 53 to use DNS.

## hook
What happens when Docker cannot find a valid nameserver in the host's /etc/resolv.conf?
