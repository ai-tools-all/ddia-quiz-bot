---
id: cops-garbage-collection
day: 1
tags: [cops, garbage-collection, version-cleanup, storage]
---

# Version Garbage Collection

## question
When can a COPS replica safely garbage collect an old version of a key?

## options
- A) Immediately after a newer version becomes visible locally
- B) After all replicas acknowledge receiving the newer version
- C) After ensuring no pending writes at any datacenter have this version in their dependency list
- D) After a fixed time window (e.g., 24 hours) expires

## answer
C

## explanation
A replica can only garbage collect version X:v after ensuring no future writes will depend on it. Since dependencies are transitive and can come from any datacenter, the system must track: (1) all replicas have received version X:v+1 or higher, AND (2) no in-flight writes have X:v in their dependency list. This requires coordination—typically using version vector protocols or explicit garbage collection messages. Premature GC breaks causal consistency: if a write W depending on X:v arrives after X:v is deleted, W cannot become visible (permanently stalled). This is a key operational challenge in COPS.

## hook
What trade-offs exist between aggressive garbage collection (saving storage) and conservative GC (preventing stalls)?
