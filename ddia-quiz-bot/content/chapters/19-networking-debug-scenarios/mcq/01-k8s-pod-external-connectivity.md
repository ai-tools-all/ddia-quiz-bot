---
id: debug-k8s-pod-external-connectivity
day: 1
tags: [networking, kubernetes, troubleshooting, nat, cni, dns]
related_stories:
  - kubernetes-networking
  - container-networking
---

# Kubernetes Pod Cannot Reach External Service

## question
A Kubernetes pod can successfully ping other pods in the cluster but cannot reach external services like google.com. The pod has a correct IP address from the pod CIDR range. What is the most likely root cause?

## options
- A) The pod's DNS resolver is misconfigured in /etc/resolv.conf
- B) The CNI plugin hasn't created the veth pair correctly
- C) The host node is missing iptables MASQUERADE rules for pod CIDR, preventing SNAT for outbound traffic
- D) The default gateway is not set in the pod's routing table

## answer
C

## explanation
If the pod can reach other pods but not external services, internal routing is working (CNI, veth pairs, and internal routing are functional). The issue is with outbound traffic to external networks. Kubernetes nodes use NAT (MASQUERADE) to translate pod IPs to node IPs for external communication. Without proper iptables MASQUERADE rules in the POSTROUTING chain for the pod CIDR, packets leave with private pod IPs that external routers cannot route back. DNS issues (A) would cause name resolution failures but not prevent IP-based connectivity. If veth pairs were broken (B), pod-to-pod communication would also fail.

## hook
How would you verify the iptables MASQUERADE rules and test if packets are being properly SNAT'd before leaving the node?
