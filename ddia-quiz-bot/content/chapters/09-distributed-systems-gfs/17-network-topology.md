---
id: ch09-network-topology
day: 1
tags: [gfs, network-topology, replica-placement]
related_stories: []
---

# Network Topology Awareness

## question
How does GFS use network topology information in replica placement?

## options
- A) It doesn't consider network topology
- B) Places all replicas in the same rack for fast access
- C) Spreads replicas across racks for availability and bandwidth
- D) Only considers geographic location

## answer
C

## explanation
GFS is network-topology aware and typically places replicas on chunk servers in different racks. This ensures availability even if an entire rack fails (power or network switch failure) and can increase aggregate bandwidth by using multiple rack switches.

## hook
Where should you store your backup copies if your entire server rack can lose power?
