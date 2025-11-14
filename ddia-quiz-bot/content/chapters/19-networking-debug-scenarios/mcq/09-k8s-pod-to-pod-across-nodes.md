---
id: debug-k8s-pod-to-pod-across-nodes
day: 9
tags: [networking, kubernetes, troubleshooting, cni, routing, encapsulation]
related_stories:
  - kubernetes-networking
  - cni-plugins
---

# Kubernetes Pod-to-Pod Communication Across Nodes Fails

## question
Pods on the same node can communicate with each other, but pods on different nodes cannot reach each other in a Kubernetes cluster. The CNI plugin (Calico) is installed. Nodes can ping each other. What is a likely cause?

## options
- A) The kube-proxy is not running on one of the nodes
- B) The BGP peering between nodes is broken, or routing tables on nodes don't have routes for remote pod CIDRs
- C) The docker0 bridge is using the wrong subnet
- D) CoreDNS is not properly configured

## answer
B

## explanation
CNI plugins like Calico use BGP (or other routing protocols) to distribute pod CIDR routes across nodes. Each node needs routes to reach pod CIDRs on other nodes. If BGP peering fails or routes aren't propagated, nodes don't know how to route packets to remote pod IPs. Since intra-node communication works, the CNI, veth pairs, and local networking are functional. kube-proxy (A) handles Service IPs, not pod-to-pod routing. docker0 (C) is typically not used in custom CNI setups. CoreDNS (D) only affects name resolution.

## hook
How would you check the BGP peering status in Calico and verify that each node has routes for all pod CIDRs in the cluster?
