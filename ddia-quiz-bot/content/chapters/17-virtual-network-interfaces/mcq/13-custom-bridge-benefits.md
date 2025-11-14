---
id: virtual-network-custom-bridge-benefits
day: 13
tags: [networking, docker, bridge, dns, best-practices]
related_stories:
  - docker-networking
---

# Custom Bridge Network Benefits

## question
What feature do custom bridge networks provide that the default docker0 bridge does not?

## options
- A) Higher network throughput and lower latency
- B) Automatic DNS resolution between containers by name
- C) Support for IPv6 addressing
- D) Direct connection to the host's physical network

## answer
B

## explanation
Custom user-defined bridge networks include Docker's embedded DNS server (127.0.0.11), enabling containers to resolve each other by name. The default bridge network lacks this feature and requires using legacy --link flags or manual IP address management.

## hook
Why is DNS-based service discovery more flexible than using IP addresses directly?
