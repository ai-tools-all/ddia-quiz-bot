---
id: occ-conflict-detection
day: 4
tags: [occ, conflicts, validation]
---

# Conflict Detection Conditions

## question
Which condition causes an OCC transaction to abort at commit?

## options
- A) Any object it read has the same version and no lock bit set
- B) Any object it read or intends to modify has a changed version or a lock bit set compared to the initial read
- C) The coordinator observes low network latency
- D) The backup replica is temporarily unreachable

## answer
B

## explanation
OCC validates that all read objects are unchanged and that write targets can be locked. A version mismatch or observed lock bit indicates a conflicting concurrent transaction, requiring abort.

## hook
What info must a client remember from initial reads to validate later?
