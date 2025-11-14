---
id: debug-k8s-service-clusterip-unreachable
day: 5
tags: [networking, kubernetes, troubleshooting, service, iptables, kube-proxy]
related_stories:
  - kubernetes-networking
  - iptables-nat
---

# Kubernetes Service ClusterIP Not Reachable

## question
A Kubernetes Service with ClusterIP 10.96.0.50 is not reachable from pods, but the backend pods are running and healthy. Running `kubectl get endpoints` shows the correct pod IPs. Direct pod-to-pod communication works. What is the most likely issue?

## options
- A) The Service's selector doesn't match the pod labels
- B) The kube-proxy daemon is not running or has failed to create the necessary iptables NAT rules for the ClusterIP
- C) The CoreDNS service is down
- D) The pod CIDR and service CIDR are overlapping

## answer
B

## explanation
ClusterIP Services are virtual IPs that don't exist on any interface. kube-proxy watches Services and creates iptables (or IPVS) rules to DNAT traffic to ClusterIPs to backend pod IPs. If kube-proxy isn't running or has errors, these rules aren't created, so traffic to the ClusterIP is never translated to pod IPs. Since direct pod-to-pod works, the CNI and basic networking are functional. Since endpoints exist (A is incorrect), the selector is matching correctly. CoreDNS (C) only affects name resolution, not IP connectivity. CIDR overlap (D) would cause routing confusion but is usually prevented by Kubernetes.

## hook
How would you inspect the iptables rules to verify if kube-proxy has correctly configured the DNAT rules for this ClusterIP Service?
