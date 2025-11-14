---
id: virtual-network-container-dns
day: 6
tags: [networking, docker, dns, service-discovery]
related_stories:
  - docker-networking
---

# Docker DNS Resolution

## question
Why can containers on a custom bridge network resolve each other by name, but containers on the default bridge network cannot?

## options
- A) Custom bridge networks have Docker's embedded DNS server (127.0.0.11) enabled, while the default bridge does not
- B) The default bridge network uses a different IP address range that doesn't support DNS
- C) Custom bridge networks require manual /etc/hosts file configuration
- D) Container name resolution only works with the overlay network driver

## answer
A

## explanation
User-defined (custom) bridge networks include automatic DNS resolution through Docker's embedded DNS server at 127.0.0.11, allowing containers to reach each other by name. The default bridge network lacks this feature and requires using legacy --link or manual IP management.

## hook
How does Docker's embedded DNS server know which IP to return for a container name?
