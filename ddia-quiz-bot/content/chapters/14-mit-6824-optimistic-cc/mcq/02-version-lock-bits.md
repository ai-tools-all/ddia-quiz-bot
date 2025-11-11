---
id: farm-version-lock-bits
day: 1
tags: [occ, version, locking, commit]
---

# Version Numbers and Lock Bits

## question
What is the purpose of the lock bit in an object's header during the FaRM commit protocol?

## options
- A) To permanently lock the object for the transaction's entire execution phase
- B) To prevent concurrent conflicting transactions from both committing to the same object simultaneously
- C) To indicate that the object has been deleted and should not be read
- D) To mark objects that require backup replication before commit

## answer
B

## explanation
The lock bit is set atomically during the commit phase (LOCK messages) to ensure that only one transaction at a time can commit changes to a given object. If another transaction tries to lock an already-locked object, it detects the set lock bit and aborts, preventing conflicting concurrent commits.

## hook
What happens if a transaction crashes while holding a lock bit?
