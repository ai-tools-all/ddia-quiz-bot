---
id: ch04-rolling-upgrade-compatibility
day: 16
tags: [deployment, rolling-upgrade, compatibility, encoding]
related_stories: []
---

# Rolling Upgrades and Compatibility

## question
During a rolling upgrade where old and new versions of a service run simultaneously, which compatibility guarantee is most critical?

## options
- A) Only backward compatibility (new code reads old data)
- B) Only forward compatibility (old code reads new data)
- C) Both backward and forward compatibility
- D) Neither, since it's a temporary state

## answer
C

## explanation
During rolling upgrades, you need BOTH backward and forward compatibility. New nodes must read data written by old nodes (backward compatibility), and old nodes must read data written by new nodes (forward compatibility). Even though it's temporary, the deployment window can be significant, and rollbacks require forward compatibility.

## hook
What's the blast radius when your schema change breaks during a rolling deployment?
