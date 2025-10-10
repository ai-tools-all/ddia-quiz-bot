---
id: ch09-write-data-flow
day: 1
tags: [gfs, data-flow, write-operations]
related_stories: []
---

# GFS Write Data Flow

## question
How does data flow during a write operation in GFS?

## options
- A) Client sends data to master, master forwards to chunk servers
- B) Client sends data to all replicas simultaneously
- C) Client sends data to nearest replica, which pipelines to others
- D) Client sends data only to primary, primary replicates to secondaries

## answer
C

## explanation
GFS separates data flow from control flow. The client pushes data to all replicas in a pipelined fashion, starting with the nearest one. Each chunk server forwards data to the next closest replica that hasn't received it, utilizing network topology efficiently.

## hook
How do you efficiently send 64MB of data to three different servers in different locations?
